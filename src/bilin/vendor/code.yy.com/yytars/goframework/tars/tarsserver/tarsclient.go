package tarsserver

import (
	"time"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"errors"
	"sync"
)

type TarsClientProtocol interface {
	Recv(pkg []byte)
	ParsePackage(buff []byte) (int, int)
}

type AdapterClient interface {
	Send(req []byte) error
	Close()
}

type TarsClientConf struct {
	Proto        string
	ClientProto  TarsClientProtocol
	NumConnect   int
	QueueLen     int
	IdleTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type TarsClient struct {
	address string
	//TODO remove it
	connWorkers []*connection

	cp        TarsClientProtocol
	conf      *TarsClientConf
	sendQueue chan []byte
}

func NewTarsClient(address string, cp TarsClientProtocol, conf *TarsClientConf) *TarsClient {
	if conf.QueueLen <= 0 {
		conf.QueueLen = 100
	}
	if conf.NumConnect <= 0 {
		conf.NumConnect = 1
	}
	sendQueue := make(chan []byte, conf.QueueLen)
	tc := &TarsClient{conf: conf, address: address, cp: cp, sendQueue: sendQueue}
	tc.connWorkers = make([]*connection, conf.NumConnect)
	for i := 0; i < conf.NumConnect; i++ {
		c := &connection{tc: tc, isClosed: true, connLock: &sync.Mutex{}}
		tc.connWorkers[i] = c
	}
	return tc
}

func (tc *TarsClient) Send(req []byte) error {
	for _, w := range tc.connWorkers {
		if err := w.reConnect(); err != nil {
			return err
		}
	}
	select {
	case tc.sendQueue <- req:
	default:
		appzaplog.Warn("sendQueue is full")
		return errors.New("sendQueue is full")
	}

	return nil
}

func (tc *TarsClient) Close() {
	for _, w := range tc.connWorkers {
		w.isClosed = true
		w.conn.Close()
	}
}