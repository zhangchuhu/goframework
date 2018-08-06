package dao

import (
	//"bilin/thrift/gen-go/common"
	"bilin/thrift/gen-go/openstatus"
	//d"bilin/userinfocenter/config"
	"bilin/common/appthrift"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"context"
	"errors"
	"git.apache.org/thrift.git/lib/go/thrift"
	"time"
)

const (
	CALL_TITMOUT = time.Second * 1
)

var connectionPool *appthrift.ConnectionPool

func InitThriftConnentPool(hosts string) {
	connectionPool = appthrift.NewConnectionPool(20, time.Minute*5, time.Second*1, 0, hosts)
}

func GetUserOpenStaus(uid uint64, version, clientType, ip string) (int32, error) {
	if connectionPool == nil {
		return 0, errors.New("thrift not init")
	}

	connection, err := connectionPool.GetConnection()
	if err != nil {
		appzaplog.Error("GetConnection error")
		return 0, errors.New("GetConnection error")
	}
	if connection == nil || connection.Transport == nil {
		appzaplog.Error("GetConnection or Transport nil")
		return 0, errors.New("GetConnection or Transport nil")
	}

	protocolFactory := thrift.NewTCompactProtocolFactory()
	iprot := thrift.NewTMultiplexedProtocol(protocolFactory.GetProtocol(connection.Transport), "OpenStatus")
	oprot := thrift.NewTMultiplexedProtocol(protocolFactory.GetProtocol(connection.Transport), "OpenStatus")
	client := openstatus.NewOpenStatusServiceClient(thrift.NewTStandardClient(iprot, oprot))
	if client == nil {
		connectionPool.ReportErrorConnection(connection)
		appzaplog.Error("NewOpenStatusServiceClientFactory error")
		return 0, errors.New("NewOpenStatusServiceClientFactory error")
	}

	ctx, cancel := context.WithTimeout(context.Background(), CALL_TITMOUT)
	defer cancel()
	staus, err := client.GetOpenStatusNew(ctx, int64(uid), version, clientType, ip)
	if err != nil {
		connectionPool.ReportErrorConnection(connection)
		appzaplog.Error("ReportErrorConnection error", zap.Error(err))
		return 0, err
	}
	connectionPool.ReturnConnection(connection)
	appzaplog.Debug("GetUserOpenStaus success", zap.Uint64("uid", uid), zap.Int32("status", staus))
	return staus, nil
}
