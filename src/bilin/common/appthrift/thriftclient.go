// @author kordenlu
// @创建时间 2017/12/01 11:54
// 功能描述:

package appthrift

import (
	"errors"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/juju/ratelimit"
	"github.com/sony/gobreaker"
	"time"
)

const (
	default_rw_timeout = 5 * time.Second

	KConnectTimeout = 1 * time.Second
)

var (
	NoAvailbaleHost         = errors.New("no available host")
	RateLimitingTriggered   = errors.New("rate limiting triggered")
	CircuitBreakerTriggered = errors.New("circuit breaker triggered")
)

type ThriftClient struct {
	//IPPORT            string
	balancer          Balancer
	readwriterTimeOut time.Duration
	ratelimiter       *ratelimit.Bucket
	circuitbreaker    *gobreaker.CircuitBreaker
}

func NewThriftClient(banlancer Balancer, opts ...ThriftClientOption) (*ThriftClient, error) {
	c := &ThriftClient{
		balancer:          banlancer,
		readwriterTimeOut: default_rw_timeout,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

// 根据balancer的配置，取一个
func (this *ThriftClient) InvokeWithTimeoutAndRetry(callback func(transport thrift.TTransport) error) error {
	if this.balancer == nil {
		return NoAvailbaleHost
	}
	ipport := this.balancer.pick()
	return this.invoke(ipport, callback)
}

func (this *ThriftClient) constructTransport(ipport string) (thrift.TTransport, error) {
	var (
		transportFactory = thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
		tsocket          *thrift.TSocket
		transport        thrift.TTransport
		err              error
	)

	if tsocket, err = thrift.NewTSocketTimeout(ipport, KConnectTimeout); err == nil {
		tsocket.SetTimeout(this.readwriterTimeOut)
		transport, err = transportFactory.GetTransport(tsocket)
	}

	return transport, err
}

func (this *ThriftClient) invoke(ipport string, callback func(transport thrift.TTransport) error) error {
	var (
		transport thrift.TTransport
		err       error
	)

	if this.ratelimiter != nil && this.ratelimiter.TakeAvailable(1) == 0 {
		return RateLimitingTriggered
	}

	circuitbreakfunc := func() (interface{}, error) {
		if transport, err = this.constructTransport(ipport); err == nil {
			if err = transport.Open(); err == nil {
				defer transport.Close()
				return nil, callback(transport)
			}
		}
		return nil, err
	}

	if this.circuitbreaker != nil {
		_, err = this.circuitbreaker.Execute(circuitbreakfunc)
	} else {
		_, err = circuitbreakfunc()
	}

	return err
}
