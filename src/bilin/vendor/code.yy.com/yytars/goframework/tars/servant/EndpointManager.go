package servant

import (
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"code.yy.com/yytars/goframework/jce/servant/taf"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/tars/util/endpoint"
	"code.yy.com/yytars/goframework/tars/servant/sd"
    "github.com/serialx/hashring"
	"net"
	"strconv"
)

type EndpointManager struct {
	comm            *Communicator
	objName         string
	refreshInterval int
	directproxy     bool

	mlock           *sync.Mutex
	consistadapters *hashring.HashRing
	adapters        map[string]*AdapterProxy // [ip:port]
	points          map[endpoint.Endpoint]int
	index           []endpoint.Endpoint
	pos             int32
}

func NewEndpointManager(objName string, comm *Communicator) *EndpointManager {
	e := &EndpointManager{
		comm:            comm,
		mlock:           new(sync.Mutex),
		adapters:        make(map[string]*AdapterProxy),
		consistadapters: hashring.New([]string{}),
		points:          make(map[endpoint.Endpoint]int),
		refreshInterval: comm.Client.refreshEndpointInterval,
	}
	e.setObjName(objName)
	return e
}

func (e *EndpointManager) setObjName(objName string) {
	pos := strings.Index(objName, "@")
	if pos > 0 {
		//[direct]
		e.objName = objName[0:pos]
		endpoints := objName[pos+1:]
		e.directproxy = true
		for _, end := range strings.Split(endpoints, ":") {
			e.points[endpoint.Parse(end)] = 0
		}
		e.index = []endpoint.Endpoint{}
		for ep, _ := range e.points {
			e.index = append(e.index, ep)
			e.consistadapters = e.consistadapters.AddNode(ep.IPPort)
			appzaplog.Debug("consist AddNode",zap.String("ipport",ep.IPPort),zap.Int("size",e.consistadapters.Size()))
		}
	} else {
		//[proxy] TODO singleton
		appzaplog.Debug("proxy mode", zap.String("objName", objName))
		e.objName = objName
		e.findAndSetObj(startFrameWorkComm().sd)
		go func() {
			loop := time.NewTicker(time.Duration(e.refreshInterval) * time.Millisecond)
			for range loop.C {
				//TODO exit
				e.findAndSetObj(startFrameWorkComm().sd)
			}
		}()
	}
}

func (e *EndpointManager) GetNextValidProxy() *AdapterProxy {
	//TODO
	var (
		ep endpoint.Endpoint
	)
	e.mlock.Lock()
	defer e.mlock.Unlock()

	length := len(e.points)
	if length == 0 {
		return nil
	}

	if e.pos > int32(length-1) {
		atomic.SwapInt32(&e.pos, 0)
	}
	n := e.pos
	ep = e.index[n]
	ipport := ep.IPPort
	atomic.AddInt32(&e.pos, 1)
	if adp, ok := e.adapters[ipport]; ok {
		if adp.Available(){
			return adp
		} else if length == 1{
			appzaplog.Warn("single node in circuitbreak open stat",zap.String("ipport",ipport))
			return nil
		}else {
			if n < int32(length-1) {
				return e.adapters[e.index[n+1].IPPort]
			} else {
				return e.adapters[e.index[n-1].IPPort]
			}
		}
	}

	adp,err := e.createProxy(ipport)
	if err != nil {
		appzaplog.Error("create adapter fail", zap.Any("endpoint", ep), zap.Error(err))
		return nil
	}
	return adp
}

func (e *EndpointManager) createProxy(ipport string) (*AdapterProxy,error) {
	appzaplog.Debug("create adapter", zap.Any("ipport", ipport))
	host,port,err := net.SplitHostPort(ipport)
	if err != nil {
		return nil,err
	}
	intPort,err := strconv.ParseInt(port,10,64)
	if err != nil {
		return nil,err
	}
	end := taf.EndpointF{
		Host:host,
		Port:int32(intPort),
		Istcp:1,
	}

	e.adapters[ipport] = NewAdapterProxy(&end, e.comm)
	return e.adapters[ipport],nil
}

func (e *EndpointManager) GetHashProxy(hashcode string) *AdapterProxy {
	e.mlock.Lock()
	length := len(e.points)
	if length == 0 {
		e.mlock.Unlock()
		return nil
	}
	intHashCode,err := strconv.ParseInt(hashcode,10,64)
	if err != nil {
		e.mlock.Unlock()
		appzaplog.Error("create adapter fail", zap.String("hashcode", hashcode),zap.Error(err))
		return nil
	}
	pos := intHashCode % int64(length)
	ep := e.index[pos]
	ipport := ep.IPPort
	if adp, ok := e.adapters[ipport]; ok {
		e.mlock.Unlock()
		if adp.Available(){
			return adp
		}
		return nil
	}
	adp,err := e.createProxy(ipport)
	e.mlock.Unlock()
	if err != nil {
		appzaplog.Error("create adapter fail", zap.Any("endpoint", ep), zap.Error(err))
		return nil
	}
	return adp
}

func (e *EndpointManager)GetConsistHashProxy(hashcode string) *AdapterProxy {
	// hot code, dont user defer for lock/unlock
	var (
		localadapter *AdapterProxy
		ipport string
		ok bool
		err error
	)
	e.mlock.Lock()
	ipport,ok = e.consistadapters.GetNode(hashcode)
	if !ok {
		e.mlock.Unlock()
		appzaplog.Warn("GetConsistHashProxy not found", zap.String("hashcode", hashcode),zap.Any("consist",e.consistadapters.Size()))
		return localadapter
	}

	if localadapter, ok = e.adapters[ipport]; ok {
		e.mlock.Unlock()
		if localadapter.Available(){
			return localadapter
		}
		appzaplog.Warn("selected node in circuitbreak open stat,wait 60s",zap.String("ipport",ipport))
		return nil
	}

	localadapter,err = e.createProxy(ipport)
	e.mlock.Unlock()
	if err != nil{
		appzaplog.Error("create adapter fail", zap.Any("ipport", ipport), zap.Error(err))
		return nil
	}

	return localadapter
}

func (e *EndpointManager) SelectAdapterProxy(msg IMessage) *AdapterProxy {
	switch  {
	case msg.consistHashEnable():
		return e.GetConsistHashProxy(msg.HashCode())
	case msg.hashEnable():
		return e.GetHashProxy(msg.HashCode())
	default:
		return e.GetNextValidProxy()
	}
}

func (e *EndpointManager) findAndSetObj(sdhelper sd.SDHelper) error{
	if sdhelper == nil {
		return NilParamsErr
	}

	activeEp := new([]taf.EndpointF)
	inactiveEp := new([]taf.EndpointF)
	ret, err := sdhelper.FindObjectByIdInSameGroup(e.objName, activeEp, inactiveEp)
	if err != nil {
		appzaplog.Error("find obj end fail", zap.Error(err))
		return err
	}
	appzaplog.Debug("find obj endpoint", zap.String("obj", e.objName), zap.Int32("ret", ret),
		zap.Any("activeEp", *activeEp),
		zap.Any("inactiveEp", *inactiveEp))

	e.mlock.Lock()
	defer e.mlock.Unlock()
	if (len(*inactiveEp)) > 0 {
		for _, ep := range *inactiveEp {
			end := endpoint.Taf2endpoint(ep)
			if _, ok := e.points[end]; ok {
				delete(e.points, end)
				ipport := end.IPPort
				delete(e.adapters, ipport)
				e.consistadapters = e.consistadapters.RemoveNode(ipport)
			}
		}
	}
	if (len(*activeEp)) > 0 {
		tag := int(time.Now().Unix())
		e.index = []endpoint.Endpoint{}
		for _, ep := range *activeEp {
			e.points[endpoint.Taf2endpoint(ep)] = tag
			e.index = append(e.index, endpoint.Taf2endpoint(ep))
			e.consistadapters = e.consistadapters.AddNode(net.JoinHostPort(ep.Host,strconv.FormatInt(int64(ep.Port),10)))
		}
		for ep, _ := range e.points {
			if e.points[ep] != tag {
				appzaplog.Info("remove ep", zap.Any("endpoint", ep))
				delete(e.points, ep)
				delete(e.adapters, ep.IPPort)
				e.consistadapters = e.consistadapters.RemoveNode(ep.IPPort)
			}
		}
	}
	return nil
}
