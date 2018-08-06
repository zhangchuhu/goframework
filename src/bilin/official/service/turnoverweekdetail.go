package service

import (
	"bilin/thrift/gen-go/turnover"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"context"
	"git.apache.org/thrift.git/lib/go/thrift"
	"time"
)

// Parameters:
//  - UID
//  - Appid
//  - StartTime
//  - EndTime
//  - Page
//  - Pagesize
//  - PropIds
//  - PlayTypes
func QueryAnchorWeekPropsRecieve(ctx context.Context, uid int64, startTime string, endTime string, page int32, pagesize int32) (r *turnover.TWeekPropsRecvInfoQueryPage, err error) {
	startTimeT, err := time.ParseInLocation(TimeLayoutOther, startTime, time.Local)
	if err != nil {
		appzaplog.Error("QueryAnchorWeekPropsRecieve time.Parse err", zap.Error(err), zap.String("startTime", startTime))
		return nil, err
	}

	endTimeT, err := time.ParseInLocation(TimeLayoutOther, endTime, time.Local)
	if err != nil {
		appzaplog.Error("QueryAnchorWeekPropsRecieve time.Parse err", zap.Error(err), zap.String("endTime", endTime))
		return nil, err
	}

	connection, err := turnOverThrift.GetConnection()
	if err != nil {
		appzaplog.Error("QueryAnchorWeekPropsRecieve GetConnection err", zap.Error(err))
		return nil, err
	}
	protocolFactory := thrift.NewTCompactProtocolFactory()
	client := turnover.NewTPropsServiceClientFactory(connection.Transport, protocolFactory)
	info, err := client.QueryAnchorWeekPropsRecieve(ctx, uid, turnover.TAppId_Bilin, startTimeT.Unix()*1000, endTimeT.Unix()*1000, page, pagesize, []int32{}, []int32{})
	if err != nil {
		turnOverThrift.ReportErrorConnection(connection)
		return nil, err
	}
	turnOverThrift.ReturnConnection(connection)
	return info, err
}

// Parameters:
//  - Sid
//  - Appid
//  - StartTime
//  - EndTime
//  - Page
//  - Pagesize
//  - UsedUid
func QueryChannelWeekPropsRecieve(ctx context.Context, sid int64, startTime string, endTime string, page int32, pagesize int32, usedUid int64) (r *turnover.TWeekPropsRecvInfoQueryPage, err error) {
	startTimeT, err := time.ParseInLocation(TimeLayoutOther, startTime, time.Local)
	if err != nil {
		appzaplog.Error("time.Parse err", zap.Error(err), zap.String("startTime", startTime))
		return nil, err
	}

	endTimeT, err := time.ParseInLocation(TimeLayoutOther, endTime, time.Local)
	if err != nil {
		appzaplog.Error("time.Parse err", zap.Error(err), zap.String("endTime", endTime))
		return nil, err
	}

	connection, err := turnOverThrift.GetConnection()
	if err != nil {
		appzaplog.Error("QueryChannelWeekPropsRecieve GetConnection err", zap.Error(err))
		return nil, err
	}
	protocolFactory := thrift.NewTCompactProtocolFactory()
	client := turnover.NewTPropsServiceClientFactory(connection.Transport, protocolFactory)
	info, err := client.QueryChannelWeekPropsRecieve(ctx, sid, turnover.TAppId_Bilin, startTimeT.Unix()*1000, endTimeT.Unix()*1000, page, pagesize, usedUid)
	if err != nil {
		turnOverThrift.ReportErrorConnection(connection)
		return nil, err
	}
	turnOverThrift.ReturnConnection(connection)
	return info, err
}
