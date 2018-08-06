package servant

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	base "code.yy.com/yytars/goframework/jce/servant/taf"
	"code.yy.com/yytars/goframework/jce/taf"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/tars/servant/protocol"
	"code.yy.com/yytars/goframework/tars/tarsserver"
)

type ObjectProxy struct {
	manager  *EndpointManager
	comm     *Communicator
	queueLen int32
}

func NewObjectProxy(comm *Communicator, objName string) *ObjectProxy {
	return &ObjectProxy{
		comm:    comm,
		manager: NewEndpointManager(objName, comm),
	}
}

func (obj *ObjectProxy) Invoke(ctx context.Context, msg *Message) error {
	adp := obj.manager.SelectAdapterProxy(msg)
	msg.Adp = adp
	if adp == nil {
		msg.Resp = &taf.ResponsePacket{
			IRet:        base.JCEADAPTERNULL,
			SResultDesc: "no adapter proxy selected",
		}
		return NoAdapterErr
	}
	if atomic.LoadInt32(&obj.queueLen) > 10000 {
		msg.Resp = &taf.ResponsePacket{
			IRet:        base.JCESERVEROVERLOAD,
			SResultDesc: "invoke queue is full",
		}
		return OverloadErr
	}
	atomic.AddInt32(&obj.queueLen, 1)
	readCh := make(chan *taf.ResponsePacket, 1)
	//todo move to TarsClientProtocol
	adp.resp.Store(msg.Req.IRequestId, readCh)
	ctx, cancle := context.WithTimeout(ctx, takeclientsynctimeinvoketimeout())
	defer func() {
		checkPanic()
		atomic.AddInt32(&obj.queueLen, -1)
		adp.resp.Delete(msg.Req.IRequestId)
		close(readCh)
		cancle()
	}()

	_, err := adp.circutbreaker.Execute(func() (interface{}, error) {
		if err := adp.Send(msg.Req); err != nil {
			appzaplog.Error("Send msg failed", zap.Error(err))
			return nil, SendErr
		}

		select {
		case <-ctx.Done():
			appzaplog.Warn("req timeout", zap.Int32("RequestId", msg.Req.IRequestId))

			msg.Resp = &taf.ResponsePacket{
				IRet:        base.JCEINVOKETIMEOUT,
				SResultDesc: "req timeout",
			}
			return nil, ReqTimeoutErr
		case msg.Resp = <-readCh:
			appzaplog.Debug("recv msg succ ", zap.Int32("RequestId", msg.Req.IRequestId))
		}
		return nil, nil
	})

	return err
}

func (obj *ObjectProxy) PbInvoke(msg *PbMessage) error {
	//now := time.Now()
	adp := obj.manager.SelectAdapterProxy(msg)
	msg.Adp = adp
	if adp == nil {
		msg.Resp = &pbtaf.ResponsePacket{
			IRet:        base.JCEADAPTERNULL,
			SResultDesc: "no adapter Proxy selected",
		}
		return NoAdapterErr
	}
	if atomic.LoadInt32(&obj.queueLen) > 10000 {
		msg.Resp = &pbtaf.ResponsePacket{
			IRet:        base.JCESERVEROVERLOAD,
			SResultDesc: "invoke queue is full",
		}
		return OverloadErr
	}
	atomic.AddInt32(&obj.queueLen, 1)
	readCh := make(chan *pbtaf.ResponsePacket, 1)
	adp.resp.Store(msg.Req.IRequestId, readCh)
	defer func() {
		checkPanic()
		atomic.AddInt32(&obj.queueLen, -1)
		adp.resp.Delete(msg.Req.IRequestId)
		close(readCh)
	}()
	//TODO adp active check
	_, err := adp.circutbreaker.Execute(func() (interface{}, error) {
		if err := adp.PbSend(msg.Req); err != nil {
			//adp.failAdd()
			appzaplog.Error("Send msg failed", zap.Error(err))
			if err == tarsserver.NetDialTimeoutErr && obj.manager != nil {
				// need to refresh endpoint cache with ratelimit
				//obj.manager.findAndSetObj(startFrameWorkComm().sd)
			}
			return nil, SendErr
		}
		//httpmetrics.DefReport(msg.Req.SFuncName, 0, now)
		select {
		//TODO USE TIMEOUT
		case <-time.After(3 * time.Second):
			appzaplog.Warn("req timeout", zap.Int32("RequestId", msg.Req.IRequestId))

			msg.Resp = &pbtaf.ResponsePacket{
				IRet:        base.JCEINVOKETIMEOUT,
				SResultDesc: "req timeout",
			}
			//adp.failAdd()
			return nil, ReqTimeoutErr
		case msg.Resp = <-readCh:
			appzaplog.Debug("recv msg succ ", zap.Int32("RequestId", msg.Req.IRequestId))
		}
		return nil, nil
	})
	return err
}

type ObjectProxyFactory struct {
	objs map[string]*ObjectProxy
	comm *Communicator
	om   *sync.Mutex
}

//func NewObjectProxyFactory(comm *Communicator) *ObjectProxyFactory {
//	o := &ObjectProxyFactory{
//		om:   new(sync.Mutex),
//		comm: comm,
//		objs: make(map[string]*ObjectProxy),
//	}
//	return o
//}
//
//func (o *ObjectProxyFactory) GetObjectProxy(objName string) *ObjectProxy {
//	if obj, ok := o.objs[objName]; ok {
//		return obj
//	}
//	obj := NewObjectProxy(o.comm, objName)
//
//	o.om.Lock()
//	defer o.om.Unlock()
//	o.objs[objName] = obj
//	return obj
//}
