package pushproxy

import (
	"bilin/common/thriftpool"
	"bilin/thrift/gen-go/tunnel"
	"context"
	"fmt"

	"bilin/bcserver/config"
	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"git.apache.org/thrift.git/lib/go/thrift"
)

const (
	TunnelService = "tunnel"
)

var (
	tunnelPool thriftpool.Pool
)

func init() {
	var err error
	tunnelPool, err = thriftpool.NewChannelPool(0, 1000, createTunnelConn)
	if err != nil {
		log.Panic("can not create thrift connection pool tunnel", zap.Any("err", err))
	}
}

func createTunnelConn() (*thriftpool.Conn, error) {
	var (
		protocolFactory  thrift.TProtocolFactory
		transportFactory thrift.TTransportFactory
		transport        thrift.TTransport
		err              error
		client           *tunnel.TunnelClient
	)
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()
	transportFactory = thrift.NewTBufferedTransportFactory(8192)
	transportFactory = thrift.NewTFramedTransportFactory(transportFactory)
	transport, err = thrift.NewTSocket(config.GetAppConfig().PushProxyAddr)
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
	client = tunnel.NewTunnelClient(thrift.NewTStandardClient(iprot, oprot))
	return &thriftpool.Conn{
		Socket: transport,
		Client: client,
	}, nil
}

func PushToUser(sid int32, uid int32, msg_type int32, msg string) (err error) {
	var ret int32
	err = thriftpool.Invoke(TunnelService, tunnelPool, func(client interface{}) (err error) {
		c := client.(*tunnel.TunnelClient)
		ret, err = c.UnicastToRoomByUidEx(context.TODO(), int64(tunnel.AppidType_NEW_BCSERVER), sid, uid, msg, msg_type)
		return
	})
	if err != nil {
		log.Error("PushToUser", zap.Any("err", err))
		return
	}
	return
}

func PushToRoom(sid int32, msg_type int32, msg string) (err error) {
	var ret int32
	err = thriftpool.Invoke(TunnelService, tunnelPool, func(client interface{}) (err error) {
		c := client.(*tunnel.TunnelClient)
		ret, err = c.BroadcastBySidEx(context.TODO(), int64(tunnel.AppidType_NEW_BCSERVER), sid, msg_type, msg)
		return
	})
	if err != nil {
		log.Error("PushToRoom", zap.Any("err", err))
		return
	}
	return
}
