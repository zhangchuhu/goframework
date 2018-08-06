// @author kordenlu
// @创建时间 2017/12/01 10:57
// 功能描述:

package appthrift

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	"net"
	"strconv"
)

type AppThriftServer struct {
	listerPort       int64
	listerHost       string
	processor        thrift.TProcessor
	transportFactory thrift.TTransportFactory
	protocolFactory  thrift.TProtocolFactory
	server           *thrift.TSimpleServer
}

const defaulthost = "0.0.0.0"

func NewAppThriftServer(port int64, processor thrift.TProcessor, opts ...ThriftServerOption) *AppThriftServer {
	server := &AppThriftServer{
		listerPort:       port,
		listerHost:       defaulthost,
		processor:        processor,
		transportFactory: thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory()),
		protocolFactory:  thrift.NewTBinaryProtocolFactory(true, true),
	}

	for _, opt := range opts {
		opt(server)
	}

	return server
}

func (this *AppThriftServer) Start() error {
	serverTransport, err := thrift.NewTServerSocket(net.JoinHostPort(this.listerHost, strconv.FormatInt(this.listerPort, 10)))
	if err != nil {
		return err
	}
	this.server = thrift.NewTSimpleServer4(this.processor, serverTransport, this.transportFactory, this.protocolFactory)
	go this.server.Serve()
	return nil
}
