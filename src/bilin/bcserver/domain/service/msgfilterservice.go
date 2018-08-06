package service

import (
	"bilin/bcserver/config"
	"bilin/common/thriftpool"
	"bilin/thrift/gen-go/bilin_msg_filter"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
)

const (
	// HotLineService 是对端定义的MultiplexedProtocol服务名
	MsgFilterService = "msgFilter"
)

func CreateMsgFilterServiceConn() (*thriftpool.Conn, error) {
	var (
		protocolFactory  thrift.TProtocolFactory
		transportFactory thrift.TTransportFactory
		transport        thrift.TTransport
		err              error
		client           *bilin_msg_filter.MsgFilterClient
	)
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()
	//protocolFactory = thrift.NewTCompactProtocolFactory()
	transportFactory = thrift.NewTBufferedTransportFactory(8192)
	transportFactory = thrift.NewTFramedTransportFactory(transportFactory)
	transport, err = thrift.NewTSocket(config.GetAppConfig().MsgFilterThriftAddr)
	if err != nil {
		return nil, fmt.Errorf("error new thrift transport: %v", err)
	}
	transport, err = transportFactory.GetTransport(transport)
	if err != nil {
		return nil, fmt.Errorf("error wrap thrift transport: %v", err)
	}
	err = transport.Open()
	if err != nil {
		return nil, fmt.Errorf("error open thrift transport: %v", err)
	}
	iprot := protocolFactory.GetProtocol(transport)
	oprot := protocolFactory.GetProtocol(transport)
	client = bilin_msg_filter.NewMsgFilterClient(thrift.NewTStandardClient(iprot, oprot))
	return &thriftpool.Conn{
		Socket: transport,
		Client: client,
	}, nil
}
