// @author kordenlu
// @创建时间 2018/03/02 17:54
// 功能描述:

package servant

import (
	"sync"
	"code.yy.com/yytars/goframework/tars/servant/sd"
)

type ICommunicator interface {
	GetServantProxy(objname string) *ServantProxy
}

type TarCommunicator struct {
	s          *ServantProxyFactory
	Client     *clientConfig
	properties sync.Map
	sd         sd.SDHelper
}

func (c *TarCommunicator) GetServantProxy(objname string) *ServantProxy {
	return c.s.getServantProxy(objname)
}