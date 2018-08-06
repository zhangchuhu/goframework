// @author kordenlu
// @创建时间 2018/03/29 18:03
// 功能描述:

package service

import (
	"bilin/bcserver/config"
	"bilin/common/thriftpool"
	"bilin/thrift/gen-go/hotline"
	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"math/rand"
	"time"
)

const (
	// HotLineService 是对端定义的MultiplexedProtocol服务名
	HotLineService     = "hotLine"
	HotLineDataService = "hotLineDataService"

	//直播间内活动定义
	UserEnterRoomTask   = "3"
	HostStartLivingTask = "4"
)

type JoinHotLineService struct {
}

func CreateHotLineServiceConn() (*thriftpool.Conn, error) {
	var (
		protocolFactory  thrift.TProtocolFactory
		transportFactory thrift.TTransportFactory
		transport        thrift.TTransport
		err              error
		client           *hotline.HotLineServiceClient
	)

	log.Info("CreateHotLineServiceConn begin", zap.Any("JavaThriftAddr length", len(config.GetAppConfig().JavaThriftAddr)))

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
	iprot := thrift.NewTMultiplexedProtocol(protocolFactory.GetProtocol(transport), HotLineService)
	oprot := thrift.NewTMultiplexedProtocol(protocolFactory.GetProtocol(transport), HotLineService)
	client = hotline.NewHotLineServiceClient(thrift.NewTStandardClient(iprot, oprot))
	return &thriftpool.Conn{
		Socket: transport,
		Client: client,
	}, nil
}

func CreateHotLineDataServiceConn() (*thriftpool.Conn, error) {
	var (
		protocolFactory  thrift.TProtocolFactory
		transportFactory thrift.TTransportFactory
		transport        thrift.TTransport
		err              error
		client           *hotline.DataServiceClient
	)
	protocolFactory = thrift.NewTCompactProtocolFactory()
	transportFactory = thrift.NewTBufferedTransportFactory(8192)
	transportFactory = thrift.NewTFramedTransportFactory(transportFactory)
	transport, err = thrift.NewTSocketTimeout(config.GetAppConfig().ActTaskThriftAddr, 2*time.Second)
	if err != nil {
		return nil, fmt.Errorf("error new thrift transport: %v", err)
	}
	//transport.(*thrift.TSocket).SetTimeout(100 * time.Millisecond) // read write timeout

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
	client = hotline.NewDataServiceClient(thrift.NewTStandardClient(iprot, oprot))
	return &thriftpool.Conn{
		Socket: transport,
		Client: client,
	}, nil
}

//直播间任务
func NewNotifyTask(uid uint64, roomid uint64, taskId string) (task *hotline.TaskReq) {
	return &hotline.TaskReq{
		UID:     int64(uid),
		TaskKey: "cashRetainTask",
		TaskID:  taskId,
		RoomID:  int64(roomid),
	}
}
