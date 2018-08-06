// @author kordenlu
// @创建时间 2017/12/01 14:12
// 功能描述:

package appthrift

import (
	"github.com/juju/ratelimit"
	"github.com/sony/gobreaker"
	"time"
)

type ThriftClientOption func(*ThriftClient)

func ReadWriteTimeout(duration time.Duration) ThriftClientOption {
	return func(client *ThriftClient) {
		client.readwriterTimeOut = duration
	}
}

func RateLimiting(rate float64, capacity int64) ThriftClientOption {
	return func(client *ThriftClient) {
		client.ratelimiter = ratelimit.NewBucketWithRate(rate, capacity)
	}
}

func CirCuitBreaker(name string, timeout time.Duration) ThriftClientOption {
	return func(client *ThriftClient) {
		client.circuitbreaker = gobreaker.NewCircuitBreaker(
			gobreaker.Settings{
				Name:    name,
				Timeout: timeout,
				OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
				},
			},
		)
	}
}
