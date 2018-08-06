package handler

import (
	"bilin/common/thriftpool"
	"bilin/thrift/gen-go/callrecord"
	"bilin/thrift/gen-go/hotline"
	"fmt"
	"math/rand"

	"bilin/thrift/gen-go/meeting"

	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"git.apache.org/thrift.git/lib/go/thrift"
)

const (
	// HotLineService 是对端定义的MultiplexedProtocol服务名
	HotLineService    = "hotLine"
	CallRecordService = "callRecord"
	MeetingService    = "meeting"
	SpamLevelService  = "spamlevel"
	// 非Multiplexed
	HotLineDataService = "hotLineDataService"
)

func createMultiplexedThriftConn(serviceName string) (thrift.TTransport, thrift.TProtocol, thrift.TProtocol, error) {
	var (
		protocolFactory  thrift.TProtocolFactory
		transportFactory thrift.TTransportFactory
		transport        thrift.TTransport
		err              error
	)

	log.Info("createMultiplexedThriftConn "+serviceName, zap.Any("JavaThriftAddr", Conf.JavaThriftAddr))
	pos := rand.Intn(len(Conf.JavaThriftAddr))

	protocolFactory = thrift.NewTCompactProtocolFactory()
	transportFactory = thrift.NewTBufferedTransportFactory(8192)
	transportFactory = thrift.NewTFramedTransportFactory(transportFactory)
	transport, err = thrift.NewTSocket(Conf.JavaThriftAddr[pos])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error new thrift transport: %v", err)
	}
	transport, err = transportFactory.GetTransport(transport)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error wrap thrift transport: %v", err)
	}
	err = transport.Open()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error open thrift transport: %v", err)
	}
	iprot := thrift.NewTMultiplexedProtocol(protocolFactory.GetProtocol(transport), serviceName)
	oprot := thrift.NewTMultiplexedProtocol(protocolFactory.GetProtocol(transport), serviceName)
	return transport, iprot, oprot, nil
}

func CreateHotLineServiceConn() (*thriftpool.Conn, error) {
	transport, iprot, oprot, err := createMultiplexedThriftConn(HotLineService)
	if err != nil {
		return nil, err
	}
	return &thriftpool.Conn{
		Socket: transport,
		Client: hotline.NewHotLineServiceClient(thrift.NewTStandardClient(iprot, oprot)),
	}, nil
}

func CreateCallRecordServiceConn() (*thriftpool.Conn, error) {
	transport, iprot, oprot, err := createMultiplexedThriftConn(CallRecordService)
	if err != nil {
		return nil, err
	}
	return &thriftpool.Conn{
		Socket: transport,
		Client: callrecord.NewCallRecordServiceClient(thrift.NewTStandardClient(iprot, oprot)),
	}, nil
}

func CreateMeetingServiceConn() (*thriftpool.Conn, error) {
	transport, iprot, oprot, err := createMultiplexedThriftConn(MeetingService)
	if err != nil {
		return nil, err
	}
	return &thriftpool.Conn{
		Socket: transport,
		Client: meeting.NewMeetingServiceClient(thrift.NewTStandardClient(iprot, oprot)),
	}, nil
}

func CreateSpamLevelServiceConn() (*thriftpool.Conn, error) {
	transport, iprot, oprot, err := createMultiplexedThriftConn(SpamLevelService)
	if err != nil {
		return nil, err
	}
	return &thriftpool.Conn{
		Socket: transport,
		Client: meeting.NewMeetingServiceClient(thrift.NewTStandardClient(iprot, oprot)),
	}, nil
}

func createThriftConn(serviceName string) (thrift.TTransport, thrift.TProtocol, thrift.TProtocol, error) {
	var (
		protocolFactory  thrift.TProtocolFactory
		transportFactory thrift.TTransportFactory
		transport        thrift.TTransport
		err              error
	)

	log.Info("createThriftConn "+serviceName, zap.Any("ActTaskThriftAddr", Conf.ActTaskThriftAddr))
	pos := rand.Intn(len(Conf.ActTaskThriftAddr))

	protocolFactory = thrift.NewTCompactProtocolFactory()
	transportFactory = thrift.NewTBufferedTransportFactory(8192)
	transportFactory = thrift.NewTFramedTransportFactory(transportFactory)
	transport, err = thrift.NewTSocket(Conf.ActTaskThriftAddr[pos])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error new thrift transport: %v", err)
	}
	transport, err = transportFactory.GetTransport(transport)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error wrap thrift transport: %v", err)
	}
	err = transport.Open()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error open thrift transport: %v", err)
	}
	iprot := protocolFactory.GetProtocol(transport)
	oprot := protocolFactory.GetProtocol(transport)
	return transport, iprot, oprot, nil
}

func CreateHotLineDataServiceConn() (*thriftpool.Conn, error) {
	transport, iprot, oprot, err := createThriftConn(HotLineDataService)
	if err != nil {
		return nil, err
	}
	return &thriftpool.Conn{
		Socket: transport,
		Client: hotline.NewDataServiceClient(thrift.NewTStandardClient(iprot, oprot)),
	}, nil
}
