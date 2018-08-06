package tarsserver

import (
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"net"
	"time"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"context"
)

const (
	PACKAGE_LESS = iota
	PACKAGE_FULL
	PACKAGE_ERROR
)

type ITarsServer interface {
	Serve() error
	Shutdown()
	IsZombie(timeout time.Duration) bool
}

type TarsProtoCol interface {
	Invoke(ctx context.Context,pkg []byte) ([]byte, error)
	ParsePackage(buff []byte) (int, int)
	InvokeTimeout(pkg []byte) ([]byte, error)
}

type TarsServerConf struct {
	Proto   string
	Address string

	MaxAccept int
	MaxInvoke int

	AcceptTimeout time.Duration
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
	HandleTimeout time.Duration
	IdleTimeout   time.Duration
}
type TarsServer struct {
	protocolImp TarsProtoCol
	listener    net.Listener
	conf        *TarsServerConf

	//worker pool

	acceptCounter chan bool
	invokeCounter chan bool

	isClosed   bool
	lastInvoke time.Time
}

func NewTarsServer(svr TarsProtoCol, conf *TarsServerConf) *TarsServer {
	ts := &TarsServer{protocolImp: svr, conf: conf}
	ts.isClosed = false
	ts.acceptCounter = make(chan bool, conf.MaxAccept)
	ts.invokeCounter = make(chan bool, conf.MaxInvoke)
	ts.lastInvoke = time.Now()
	return ts
}

func (ts *TarsServer) Serve() error {
	cfg := ts.conf
	if cfg.Proto == "tcp" {
		lis, err := net.Listen(cfg.Proto, cfg.Address)
		if err != nil {
			return err
		}
		defer lis.Close()
		appzaplog.Debug("Listening on", zap.Any("addr", cfg.Address))
		ts.listener = lis
		for {
			if ts.isClosed {
				return nil
			}
			lis.(*net.TCPListener).SetDeadline(time.Now().Add(cfg.AcceptTimeout)) // set accept timeout
			ts.acceptCounter <- true
			conn, err := lis.Accept()
			if err != nil {
				<-ts.acceptCounter
				netErr, ok := err.(net.Error)
				if ok && netErr.Timeout() && netErr.Temporary() {
					//no new conn, not error
				} else {
					appzaplog.Error("Error accepting", zap.Error(err))
				}
				continue
			}
			handler := NewConnectHandler(ts, conn)
			appzaplog.Debug("TCP recv", zap.Any("addr", conn.RemoteAddr()))
			go handler.recv()
			go handler.loopsend()
		}
	} else if cfg.Proto == "udp" {
		udpAddr, err := net.ResolveUDPAddr("udp4", cfg.Address)
		if err != nil {
			return err
		}
		conn, err := net.ListenUDP("udp4", udpAddr)
		if err != nil {
			return err
		}
		handler := NewConnectHandler(ts, conn)
		appzaplog.Debug("UDP listen", zap.Any("addr", conn.LocalAddr()))
		ts.acceptCounter <- true
		handler.recv()
	}
	return nil
}

func (ts *TarsServer) Shutdown() {
	ts.isClosed = true
}

func (ts *TarsServer) GetConfig() *TarsServerConf {
	return ts.conf
}

// invokeCounter已满,超过timeout时间未处理
func (ts *TarsServer) IsZombie(timeout time.Duration) bool {
	return len(ts.invokeCounter) == cap(ts.invokeCounter) && ts.lastInvoke.Add(timeout).Before(time.Now())
}
