package servant

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"math"
	"sync"
	"time"
	end "code.yy.com/yytars/goframework/jce/servant/taf"
	"code.yy.com/yytars/goframework/jce/taf"
	"code.yy.com/yytars/goframework/jce_parser/gojce"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/tars/servant/protocol"
	"code.yy.com/yytars/goframework/tars/tarsserver"
	"code.yy.com/yytars/goframework/kissgo/gobreaker"
	"net"
	"strconv"
)

type TarsClientProtocolImpl struct {
	ada *AdapterProxy
}

func (c *TarsClientProtocolImpl) ParsePackage(buff []byte) (int, int) {
	return TafRequest(buff)
}

func (c *TarsClientProtocolImpl) Recv(pkg []byte) {
	defer func() {
		//TODO readCh在load之后一定几率被超时关闭了,这个时候需要recover恢复
		//或许有更好的办法吧
		if err := recover(); err != nil {
			appzaplog.Error("recv pkg painc", zap.Any("err", err))
		}
	}()
	packet := taf.ResponsePacket{}
	is := gojce.NewInputStream(pkg)
	err := packet.ReadFrom(is)
	if err != nil {
		appzaplog.Error("decode packet error", zap.Error(err))
		return
	}
	chIF, ok := c.ada.resp.Load(packet.IRequestId)
	if ok {
		ch := chIF.(chan *taf.ResponsePacket)
		//appzaplog.Debug("IN:", zap.Any("packet", packet))
		ch <- &packet
	} else {
		appzaplog.Warn("timeout resp,drop it",
			zap.Any("packet", packet),
			zap.Int32("RequestId", packet.IRequestId))
	}
}

// PB client prootocol
type TarsClientProtocolPbImpl struct {
	ada *AdapterProxy
}

func (c *TarsClientProtocolPbImpl) ParsePackage(buff []byte) (int, int) {
	return TafRequest(buff)
}

func (c *TarsClientProtocolPbImpl) Recv(pkg []byte) {
	defer func() {
		//TODO readCh在load之后一定几率被超时关闭了,这个时候需要recover恢复
		//或许有更好的办法吧
		if err := recover(); err != nil {
			appzaplog.Error("recv pkg painc", zap.Any("err", err))
		}
	}()
	packet := pbtaf.ResponsePacket{}
	if err := proto.Unmarshal(pkg, &packet); err != nil {
		appzaplog.Error("decode packet error", zap.Error(err))
		return
	}

	chIF, ok := c.ada.resp.Load(packet.IRequestId)
	if ok {
		ch := chIF.(chan *pbtaf.ResponsePacket)
		//appzaplog.Debug("IN:", zap.Any("packet", packet))
		ch <- &packet
	} else {
		appzaplog.Warn("timeout resp,drop it",
			zap.Any("packet", packet),
			zap.Int32("RequestId", packet.IRequestId))
	}
}

type AdapterProxy struct {
	resp       sync.Map
	point      *end.EndpointF
	tarsClient *tarsserver.TarsClient
	comm       *Communicator
	circutbreaker *gobreaker.CircuitBreaker
}

func NewAdapterProxy(point *end.EndpointF, comm *Communicator) *AdapterProxy {
	var st gobreaker.Settings
	strport := strconv.FormatInt(int64(point.Port),10)
	st.Interval = 60*time.Second
	st.Name = net.JoinHostPort(point.Host,strport)
	st.ReadyToTrip = func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return counts.ConsecutiveFailures>100 || (counts.Requests >= 3 && failureRatio >= 0.6)
	}

	c := &AdapterProxy{
		comm:  comm,
		point: point,
		circutbreaker:gobreaker.NewCircuitBreaker(st),
	}

	proto := "tcp"
	if point.Istcp == 0 {
		proto = "udp"
	}
	numconnect, _ := c.comm.GetPropertyInt(clientKeyNetConnectionNum)
	conf := &tarsserver.TarsClientConf{
		Proto:        proto,
		NumConnect:   numconnect,
		QueueLen:     10000,
		IdleTimeout:  time.Second * 600,
		ReadTimeout:  time.Millisecond * 100,
		WriteTimeout: time.Millisecond * 100,
	}

	if _, ok := comm.GetProperty(PBSERVANT); ok {
		appzaplog.Info("start full pb support")
		c.tarsClient = tarsserver.NewTarsClient(
			fmt.Sprintf("%s:%d", point.Host, point.Port),
			&TarsClientProtocolPbImpl{
				ada: c,
			},
			conf)
	} else {
		c.tarsClient = tarsserver.NewTarsClient(
			fmt.Sprintf("%s:%d", point.Host, point.Port),
			&TarsClientProtocolImpl{
				ada: c,
			},
			conf)
	}

	return c
}

func (c *AdapterProxy)Available()bool  {
	return c.circutbreaker.Available()
}

func (c *AdapterProxy) Send(req *taf.RequestPacket) error {
	appzaplog.Debug("send req", zap.Int32("RequestId", req.IRequestId))
	sbuf := bytes.NewBuffer(make([]byte, 0, 256))
	sbuf.Write(make([]byte, 4))
	os := gojce.NewOutputStream()
	req.WriteTo(os)
	bs := os.ToBytes()
	sbuf.Write(bs)
	len := sbuf.Len()
	binary.BigEndian.PutUint32(sbuf.Bytes(), uint32(len))

	return c.tarsClient.Send(sbuf.Bytes())
}

func (c *AdapterProxy) PbSend(req *pbtaf.RequestPacket) error {
	appzaplog.Debug("send req", zap.Int32("RequestId", req.IRequestId))
	rsp, err := proto.Marshal(req)
	if err != nil {
		appzaplog.Error("Marshal rspPackage failed", zap.Error(err))
		return err
	}
	var (
		length uint
	)
	const (
		sizeLen = 4
	)

	if length = uint(len(rsp)) + sizeLen; length > math.MaxUint32 {
		return errors.New(fmt.Sprintf("grpc: message too large (%d bytes)", length))
	}

	var buf = make([]byte, length)
	binary.BigEndian.PutUint32(buf[0:], uint32(length))
	copy(buf[sizeLen:], rsp)

	return c.tarsClient.Send(buf)
}

func (c *AdapterProxy) GetPoint() *end.EndpointF {
	return c.point
}
