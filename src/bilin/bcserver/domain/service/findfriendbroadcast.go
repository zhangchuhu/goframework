package service

import (
	"bilin/bcserver/config"
	"bilin/common/thriftpool"
	"bilin/thrift/gen-go/findfriendsbroadcast"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"math/rand"
)

const (
	FindFriendsBroadcastService = "findFriendsBroadcast"
)

func CreateFindFriendsBroadcastServiceConn() (*thriftpool.Conn, error) {
	var (
		protocolFactory  thrift.TProtocolFactory
		transportFactory thrift.TTransportFactory
		transport        thrift.TTransport
		err              error
		client           *findfriendsbroadcast.FindFriendsBroadcastServiceClient
	)

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
	iprot := thrift.NewTMultiplexedProtocol(protocolFactory.GetProtocol(transport), FindFriendsBroadcastService)
	oprot := thrift.NewTMultiplexedProtocol(protocolFactory.GetProtocol(transport), FindFriendsBroadcastService)
	client = findfriendsbroadcast.NewFindFriendsBroadcastServiceClient(thrift.NewTStandardClient(iprot, oprot))
	return &thriftpool.Conn{
		Socket: transport,
		Client: client,
	}, nil
}
