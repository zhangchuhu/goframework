package appthrift

import (
	"sync"
)

type Balancer interface {
	pick() string //ip:port
}

type rrBalancer struct {
	lock    sync.Mutex
	next    int
	ipports []string
}

func NewRRBalancer(ipports ...string) Balancer {
	r := &rrBalancer{}
	r.ipports = append(r.ipports, ipports...)
	return r
}

func (p *rrBalancer) pick() string {
	p.lock.Lock()
	ipport := p.ipports[p.next]
	p.next = (p.next + 1) % len(p.ipports)
	p.lock.Unlock()
	return ipport
}
