package appthrift

import (
	"container/list"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"git.apache.org/thrift.git/lib/go/thrift"
)

type Connection struct {
	Transport   thrift.TTransport
	createTime  time.Time
	calledTimes int64
}

type ConnectionPool struct {
	sync.Mutex

	pool_size int
	live_time time.Duration
	time_out  time.Duration
	hosts     string

	activeList   *list.List
	inactiveList *list.List

	clientTimes int64 //the connection will be closed after clientTimes,0 is unlimited
}

//Create Pool
func NewConnectionPool(
	pool_size int,
	live_time time.Duration,
	time_out time.Duration,
	clientTimes int64,
	hosts string) *ConnectionPool {

	p := &ConnectionPool{
		pool_size:    pool_size,
		live_time:    live_time,
		time_out:     time_out,
		clientTimes:  clientTimes,
		hosts:        hosts,
		activeList:   list.New(),
		inactiveList: list.New(),
	}
	return p
}

var connectionReleaseNotifyCh = make(chan interface{}, 128)

//get a connection from pool
func (p *ConnectionPool) GetConnection() (*Connection, error) {
	removeList := list.New()

	if p.activeList.Len() >= p.pool_size {
		select {
		case <-connectionReleaseNotifyCh:
			//
		case <-time.After(time.Millisecond * 200):

			p.Lock()
			now := time.Now()
			for e := p.activeList.Front(); e != nil; e = e.Next() {
				connection := e.Value.(*Connection)
				diff := now.Sub(connection.createTime)
				if diff >= p.live_time {
					removeList.PushBack(connection)
					p.activeList.Remove(e)
				}
			}
			p.Unlock()
		}
	}

	var client *Connection = nil
	p.Lock()
	for e := p.inactiveList.Front(); e != nil; e = e.Next() {
		client = e.Value.(*Connection)
		if !p.isConnectionOpen(client) {
			p.inactiveList.Remove(e)
			client = nil
			continue
		}

		diff := time.Now().Sub(client.createTime)
		if diff >= p.live_time {
			removeList.PushBack(client)
			p.inactiveList.Remove(e)
			client = nil
			continue
		}

		p.inactiveList.Remove(e)
		p.activeList.PushBack(client)
		break
	}
	p.Unlock()

	for e := removeList.Front(); e != nil; e = e.Next() {
		client := e.Value.(*Connection)
		p.closeConnection(client)
	}

	if client != nil {
		return client, nil
	}

	transport := p.createTransport()
	if transport == nil {
		return nil, fmt.Errorf("Can not open transport")
	}

	connection := &Connection{Transport: transport, createTime: time.Now(), calledTimes: 0}
	p.Lock()
	p.activeList.PushBack(connection)
	p.Unlock()
	return connection, nil
}

func (p *ConnectionPool) ReturnConnection(client *Connection) {
	var connection *Connection
	var findInActiveList bool = false

	p.Lock()
	for e := p.activeList.Front(); e != nil; e = e.Next() {
		connection = e.Value.(*Connection)
		if client == connection {
			findInActiveList = true
			p.activeList.Remove(e)
			break
		}
	}
	p.Unlock()

	if findInActiveList {
		if p.clientTimes > 0 && connection.calledTimes >= p.clientTimes {
			p.closeConnection(client)
			return
		}
		connection.calledTimes++
		connection.createTime = time.Now()
		p.Lock()
		if p.inactiveList.Len() < p.pool_size {
			p.inactiveList.PushBack(connection)
			client = nil
		}
		p.Unlock()

		if client != nil {
			fmt.Println("Reach Max Close")
			p.closeConnection(client)
		}
	}
}

func (p *ConnectionPool) ReportErrorConnection(client *Connection) {
	p.Lock()
	for e := p.activeList.Front(); e != nil; e = e.Next() {
		connection := e.Value.(*Connection)
		if client == connection {
			p.activeList.Remove(e)
			break
		}
	}
	p.Unlock()

	p.closeConnection(client)
}

func (p *ConnectionPool) createTransport() thrift.TTransport {
	hostlist := strings.Split(p.hosts, ";")
	if len(hostlist) == 0 {
		return nil
	}

	var transportFactory thrift.TTransportFactory
	transportFactory = thrift.NewTBufferedTransportFactory(8192)
	transportFactory = thrift.NewTFramedTransportFactory(transportFactory)
	//transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())

	hostlen := len(hostlist)
	hostidx := rand.Intn(hostlen)
	for i := 0; i < hostlen; i++ {
		hostidx++

		socket, err := thrift.NewTSocketTimeout(hostlist[hostidx%hostlen], time.Second*2)
		if err != nil {
			continue
		}

		if err := socket.Open(); err != nil {
			continue
		}

		useTransport, err := transportFactory.GetTransport(socket)
		if err != nil {
			continue
		}
		return useTransport
	}

	return nil
}

func (p *ConnectionPool) isConnectionOpen(client *Connection) bool {
	return client.Transport.IsOpen()
}

func (p *ConnectionPool) closeConnection(client *Connection) error {
	return client.Transport.Close()
}
