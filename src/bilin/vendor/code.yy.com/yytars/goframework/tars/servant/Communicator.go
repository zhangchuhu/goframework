// @author kordenlu
// @创建时间 2018/02/05 14:59
// 功能描述: 通信器,用于创建和维护客户端proxy

package servant

import (
	"sync"
	"code.yy.com/yytars/goframework/jce/servant/taf"
	s "code.yy.com/yytars/goframework/tars/servant/model"
	"github.com/juju/ratelimit"
	"code.yy.com/yytars/goframework/tars/servant/sd"
)

type ProxyPrx interface {
	SetServant(s.Servant)
}

type Communicator struct {
	s          *ServantProxyFactory
	Client     *clientConfig
	properties sync.Map
	sd         sd.SDHelper
}

func pbHeaderEnable(c *Communicator)  {
	c.SetProperty(PBSERVANT, "yes")
}

var(
	pbonce           sync.Once
	gPbComm *Communicator
)

// global communicator
func NewPbCommunicator() *Communicator {
	pbonce.Do(func() {
		startFrameWorkComm()
		c := new(Communicator)
		if GetClientConfig() != nil {
			c.Client = GetClientConfig()
			c.SetProperty(clientKeyNetConnectionNum,c.Client.netconnectionnum)
		} else {
			c.Client = &clientConfig{
				Locator:"",
				stat:"",
				property:"",
				modulename:"",
				refreshEndpointInterval:60000,
				reportInterval:10000,
			}
			c.SetProperty(clientKeyNetConnectionNum,2)
		}
		if GetServerConfig() != nil {
			c.SetProperty("notify", GetServerConfig().notify)
			c.SetProperty("node", GetServerConfig().Node)
			c.SetProperty("server", GetServerConfig().Server)
		}
		c.s = NewServantProxyFactory(c)
		//pbHeaderEnable(c)
		gPbComm = c
	})

	return gPbComm
}

func (c *Communicator) GetServantProxy(objname string) *ServantProxy {
	return c.s.getServantProxy(objname)
}

func (c *Communicator) SetProperty(key string, value interface{}) {
	if key == clientKeyLocator {
		if locator, ok := value.(string); ok {
			startFrameWorkComm().setQueryPrx(locator)
			startFrameWorkComm().properties.Store(key,value)
		}
	}
	c.properties.Store(key, value)
}

func (c *Communicator) GetProperty(key string) (string, bool) {
	v, ok := c.properties.Load(key)
	if v == nil {
		return "", ok
	}
	return v.(string), ok
}

func (c *Communicator) GetPropertyInt(key string) (int, bool) {
	v, ok := c.properties.Load(key)
	if v == nil {
		return 0, ok
	}
	return v.(int), ok
}

func (c *Communicator) setQueryPrx(obj string) {
	qf := new(taf.QueryF)
	qf.SetServant(c.GetServantProxy(obj))
	c.sd = sd.NewQueryFHelper(ratelimit.NewBucketWithRate(10,10), qf)
}

var (
	once           sync.Once
	gFrameworkComm *Communicator
)

func startFrameWorkComm() *Communicator {
	once.Do(func() {
		c := new(Communicator)
		//c.init()
		if GetClientConfig() != nil {
			c.Client = GetClientConfig()
			c.SetProperty(clientKeyNetConnectionNum, GetClientConfig().netconnectionnum)
		} else {
			c.Client = &clientConfig{
				Locator:"",
				stat:"",
				property:"",
				modulename:"",
				refreshEndpointInterval:60000,
				reportInterval:10000,
			}
			c.SetProperty(clientKeyNetConnectionNum, 2)
		}

		if GetServerConfig() != nil {
			c.SetProperty(serKeyNotify, GetServerConfig().notify)
			c.SetProperty(serKeyNode, GetServerConfig().Node)
			c.SetProperty(serKeyServer, GetServerConfig().Server)
		}

		c.s = NewServantProxyFactory(c)
		if GetClientConfig() != nil {
			c.setQueryPrx(GetClientConfig().Locator)
			c.properties.Store(clientKeyLocator, GetClientConfig().Locator)
		}
		c.SetProperty(clientKeyNetConnectionNum, 1)
		gFrameworkComm = c
	})
	return gFrameworkComm
}
