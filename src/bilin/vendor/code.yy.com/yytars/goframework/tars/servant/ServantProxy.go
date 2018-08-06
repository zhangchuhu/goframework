package servant

import (
	"context"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"strings"
	"sync"
	"sync/atomic"
	"code.yy.com/yytars/goframework/jce/taf"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/tars/servant/protocol"
	"errors"
)

type ServantProxy struct {
	sid  int32
	name string // appname.servername.objname
	comm *Communicator
	obj  *ObjectProxy
}

func NewServantProxy(comm *Communicator, objName string) *ServantProxy {
	s := &ServantProxy{
		comm: comm,
	}
	pos := strings.Index(objName, "@")
	if pos > 0 {
		s.name = objName[0:pos]
	} else {
		s.name = objName
	}
	s.obj = NewObjectProxy(comm, objName)
	return s
}

func (s *ServantProxy) Taf_invoke(ctx context.Context, ctype byte, sFuncName string,
	buf []byte) (*taf.ResponsePacket, error) {
	//TODO 重置sid，防止溢出
	atomic.CompareAndSwapInt32(&s.sid, 1<<31-1, 1)
	ctxmap,_ := FromOutgoingContext(ctx)

	//appzaplog.Debug("ctxmap info",zap.Any("ctxmap",ctxmap))
	req := taf.RequestPacket{
		IVersion:     1,
		CPacketType:  ctype,
		IRequestId:   atomic.AddInt32(&s.sid, 1),
		SServantName: s.name,
		SFuncName:    sFuncName,
		SBuffer:      buf,
		ITimeout:     3000,
		Context:      ctxmap,
		//Status:       statusmap,
	}

	msg := &Message{Req: &req, Ser: s, Obj: s.obj}
	if key,ok := ctxmap[CONTEXTCONSISTHASHKEY];ok{
		//appzaplog.Debug("consisthashkey set",zap.String("consisthashkey",key))
		msg.setConsistHashCode(key)
	}
	msg.Init()
	//var (
	//	succ    int32 = 1
	//	timeout int32
	//	exec    int32
	//)
	//defer func() {
	//	msg.End()
	//	ReportStat(msg, succ, timeout, exec)
	//}()
	appzaplog.Debug("Taf_invoke", zap.String("sFuncName", sFuncName),
		zap.String("obj", s.name), zap.Int32("IRequestId", req.IRequestId))
	if err := s.obj.Invoke(ctx,msg); err != nil {
		appzaplog.Error("Invoke error", zap.String("ServantName", s.name),
			zap.String("FuncName", sFuncName), zap.Int32("IRequestId", req.IRequestId),zap.Error(err))
		//TODO report exec
		//timeout = 1
		//succ = 0
		return nil,err
	}

	if msg.Resp != nil{
		if errstr,ok := msg.Resp.Status[STATUSERRSTR];ok{
			return msg.Resp,errors.New(errstr)
		}
	}
	return msg.Resp, nil
}

func (s *ServantProxy) Pb_invoke(ctx context.Context, ctype byte, sFuncName string,
	buf []byte, status map[string]string, context map[string]string) (*pbtaf.ResponsePacket, error) {
	//TODO 重置sid，防止溢出
	atomic.CompareAndSwapInt32(&s.sid, 1<<31-1, 1)
	req := pbtaf.RequestPacket{
		IVersion:     1,
		CPacketType:  pbtaf.RequestPacket_PacketType(ctype),
		IRequestId:   atomic.AddInt32(&s.sid, 1),
		SServantName: s.name,
		SFuncName:    sFuncName,
		SBuffer:      buf,
		ITimeout:     3000,
		Context:      context,
		Status:       status,
	}

	msg := &PbMessage{Req: &req, Ser: s, Obj: s.obj}
	msg.Init()
	var (
		succ    int32 = 1
		timeout int32
		exec    int32
	)
	defer func() {
		msg.End()
		ReportStat(msg, succ, timeout, exec)
	}()
	appzaplog.Debug("Taf_invoke", zap.String("sFuncName", sFuncName),
		zap.String("obj", s.name), zap.Int32("IRequestId", req.IRequestId))
	if err := s.obj.PbInvoke(msg); err != nil {
		appzaplog.Error("Invoke error", zap.String("ServantName", s.name),
			zap.String("FuncName", sFuncName),
			zap.Int32("IRequestId",req.IRequestId),
			zap.Error(err))
		//TODO report exec
		timeout = 1
		succ = 0
		return nil,err
	}

	return msg.Resp, nil
}

type ServantProxyFactory struct {
	lk   sync.RWMutex
	objs map[string]*ServantProxy
	comm *Communicator
}

func NewServantProxyFactory(comm *Communicator) *ServantProxyFactory {
	return &ServantProxyFactory{
		comm: comm,
		objs: make(map[string]*ServantProxy),
	}
}

func (o *ServantProxyFactory) getServantProxy(objName string) *ServantProxy {
	proxy := o.getProxy(objName)
	if proxy != nil{
		return proxy
	}
	return o.createProxy(objName)
}

func (o *ServantProxyFactory) getProxy(objName string) *ServantProxy {
	o.lk.RLock()
	defer o.lk.RUnlock()
	if obj, ok := o.objs[objName]; ok {
		return obj
	}
	return nil
}

func (o *ServantProxyFactory) createProxy(objName string) *ServantProxy {
	o.lk.Lock()
	defer o.lk.Unlock()
	if obj, ok := o.objs[objName]; ok {
		return obj
	}
	obj := NewServantProxy(o.comm, objName)
	o.objs[objName] = obj
	return obj
}
