package service

import (
	"bilin/common/appthrift"
	"bilin/guildtars/config"
	"bilin/thrift/gen-go/turnover"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"context"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"time"
)

var contractThriftAddr = "58.215.52.27:6907"

var contractThriftPool *appthrift.ConnectionPool

func InitTurnOverService(conf *config.AppConfig) error {
	if conf == nil || conf.ContractThriftAddr == "" {
		return fmt.Errorf("config not init")
	}
	contractThriftPool = appthrift.NewConnectionPool(20, time.Minute*1, time.Second*3, 1000, conf.ContractThriftAddr)
	return nil
}

// Parameters:
//  - UID 主播uid
//  - Sid 频道号
//  - Owuid ow uid
//  - Weight 工会抽成比例
func AddContractInfoExternal(ctx context.Context, uid int64, sid int64, owuid int64, weight int32) (r int32, err error) {
	connection, err := contractThriftPool.GetConnection()
	if err != nil {
		appzaplog.Error("GetConnection err", zap.Error(err))
		return -1, err
	}
	protocolFactory := thrift.NewTCompactProtocolFactory()
	client := turnover.NewTContractServiceClientFactory(connection.Transport, protocolFactory)
	info, err := client.AddContractInfoExternal(ctx, uid, turnover.TAppId_Bilin, sid, owuid, weight, 0)
	if err != nil {
		appzaplog.Error("AddContractInfoExternal err", zap.Error(err))
		contractThriftPool.ReportErrorConnection(connection)
		return -1, err
	}
	contractThriftPool.ReturnConnection(connection)
	return info, nil
}

// Parameters:
//  - UID
//  - Appid
func QueryContractByAnchor(ctx context.Context, uid int64) (r *turnover.TContract, err error) {
	connection, err := contractThriftPool.GetConnection()
	if err != nil {
		appzaplog.Error("GetConnection err", zap.Error(err))
		return nil, err
	}
	protocolFactory := thrift.NewTCompactProtocolFactory()
	client := turnover.NewTContractServiceClientFactory(connection.Transport, protocolFactory)
	info, err := client.QueryContractByAnchor(ctx, uid, turnover.TAppId_Bilin)
	if err != nil {
		appzaplog.Error("AddContractInfoExternal err", zap.Error(err))
		contractThriftPool.ReportErrorConnection(connection)
		return nil, err
	}
	contractThriftPool.ReturnConnection(connection)
	return info, nil
}
