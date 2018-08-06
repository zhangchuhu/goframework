package service

import (
	"bilin/bcserver/config"
	"bilin/common/thriftpool"
	"bilin/thrift/gen-go/user"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"math/rand"
)

const (
	// HotLineService 是对端定义的MultiplexedProtocol服务名
	UserServervice = "user"
)

func CreateUserServiceConn() (*thriftpool.Conn, error) {
	var (
		protocolFactory  thrift.TProtocolFactory
		transportFactory thrift.TTransportFactory
		transport        thrift.TTransport
		err              error
		client           *user.UserServiceClient
	)
	//利用当前时间的UNIX时间戳初始化rand包
	pos := rand.Intn(len(config.GetAppConfig().JavaThriftAddr))

	//protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()
	protocolFactory = thrift.NewTCompactProtocolFactory()
	transportFactory = thrift.NewTBufferedTransportFactory(8192)
	transportFactory = thrift.NewTFramedTransportFactory(transportFactory)
	transport, err = thrift.NewTSocket(config.GetAppConfig().JavaThriftAddr[pos])
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
	iprot := thrift.NewTMultiplexedProtocol(protocolFactory.GetProtocol(transport), UserServervice)
	oprot := thrift.NewTMultiplexedProtocol(protocolFactory.GetProtocol(transport), UserServervice)
	client = user.NewUserServiceClient(thrift.NewTStandardClient(iprot, oprot))
	return &thriftpool.Conn{
		Socket: transport,
		Client: client,
	}, nil
}
