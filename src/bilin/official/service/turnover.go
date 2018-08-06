package service

import (
	"bilin/common/appthrift"
	"bilin/thrift/gen-go/turnover"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"context"
	"git.apache.org/thrift.git/lib/go/thrift"
	"strconv"
	"time"
)

//var turnOverThriftAddr = "58.215.52.27:6903"

var turnOverThrift *appthrift.ConnectionPool

//func init() {
//	turnOverThrift = appthrift.NewConnectionPool(3, time.Minute*1, time.Second*3, 1000, turnOverThriftAddr)
//}

func BilinCumulativeProfit(ctx context.Context, uid uint64) (uint64, error) {
	connection, err := turnOverThrift.GetConnection()
	if err != nil {
		appzaplog.Error("BilinCumulativeProfit GetConnection err", zap.Error(err))
		return 0, err
	}
	protocolFactory := thrift.NewTCompactProtocolFactory()
	client := turnover.NewTUserAccountServiceClientFactory(connection.Transport, protocolFactory)

	count, err := client.BilinCumulativeProfit(ctx, int64(uid))
	if err != nil {
		appzaplog.Error("BilinCumulativeProfit err", zap.Error(err))
		turnOverThrift.ReportErrorConnection(connection)
		return 0, err
	}
	turnOverThrift.ReturnConnection(connection)
	return uint64(count), nil
}

func QueryCurMonthRevenueRecord(ctx context.Context, uid int64, sid int64, revenueUserType int32, anchorUid uint64) (r []*turnover.TRevenueRecord, err error) {
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	info, totalpagesize, err := QueryRevenueRecord(ctx, uid, sid, firstOfMonth.Format(TimeLayoutOther), lastOfMonth.Format(TimeLayoutOther),
		revenueUserType, 1, 10, anchorUid)
	if err != nil {
		appzaplog.Error("QueryRevenueRecord err", zap.Error(err))
		return info, err
	}

	info, _, err = QueryRevenueRecord(ctx, uid, sid, firstOfMonth.Format(TimeLayoutOther), lastOfMonth.Format(TimeLayoutOther),
		revenueUserType, 1, int32(totalpagesize), anchorUid)
	if err != nil {
		appzaplog.Error("QueryRevenueRecord err", zap.Error(err))
		return info, err
	}
	return info, err
}

func QueryRevenueRecord(ctx context.Context, uid int64, sid int64, timeGreaterThan string,
	timeLessThan string, revenueUserType int32, page int32, pagesize int32, anchorUid uint64) (r []*turnover.TRevenueRecord, totalPageSize int64, err error) {
	appzaplog.Debug("[+]QueryRevenueRecord", zap.String("starttime", timeGreaterThan), zap.String("endtime", timeLessThan),
		zap.Uint64("hostuid", anchorUid), zap.Int64("uid", uid))
	startTimeT, err := time.ParseInLocation(TimeLayoutOther, timeGreaterThan, time.Local)
	if err != nil {
		appzaplog.Error("QueryRevenueRecord ParseInLocation err", zap.Error(err))
		return nil, 0, err
	}

	endTimeT, err := time.ParseInLocation(TimeLayoutOther, timeLessThan, time.Local)
	if err != nil {
		appzaplog.Error("QueryRevenueRecord ParseInLocation err", zap.Error(err))
		return nil, 0, err
	}

	info, err := QueryRevenueRecordPaging(ctx, uid, pagesize, page, turnover.TUserType(revenueUserType), startTimeT.Unix()*1000, endTimeT.Unix()*1000,
		int64(anchorUid), sid)
	if err != nil {
		appzaplog.Error("QueryRevenueRecord QueryRevenueRecordPaging err", zap.Error(err))
		return nil, 0, err
	}
	totalPageSize = int64(info.TotalElement)
	for _, v := range info.Content {
		tr := &turnover.TRevenueRecord{}

		if contributeUid, ok := v["contributeUid"]; ok {
			tr.ContributeUid, err = strconv.ParseInt(contributeUid, 10, 64)
			if err != nil {
				appzaplog.Error("ParseInt err", zap.Error(err), zap.String("contributeUid", contributeUid))
				continue
			}
		}

		if uid_, ok := v["uid"]; ok {
			tr.UID, err = strconv.ParseInt(uid_, 10, 64)
			if err != nil {
				appzaplog.Error("ParseInt err", zap.Error(err), zap.String("uid_", uid_))
				continue
			}
		}

		if sid_, ok := v["sid"]; ok {
			tr.Sid, err = strconv.ParseInt(sid_, 10, 64)
			if err != nil {
				appzaplog.Error("ParseInt err", zap.Error(err), zap.String("sid_", sid_))
				continue
			}
		}

		if income, ok := v["income"]; ok {
			tr.Income, err = strconv.ParseFloat(income, 64)
			if err != nil {
				appzaplog.Error("ParseFloat err", zap.Error(err), zap.String("income", income))
				continue
			}
		}
		if incomeRate, ok := v["incomeRate"]; ok {
			tr.IncomeRate, err = strconv.ParseFloat(incomeRate, 64)
			if err != nil {
				appzaplog.Error("ParseFloat err", zap.Error(err), zap.String("incomeRate", incomeRate))
				continue
			}
		}
		if realIncome, ok := v["realIncome"]; ok {
			tr.RealIncome, err = strconv.ParseFloat(realIncome, 64)
			if err != nil {
				appzaplog.Error("ParseFloat err", zap.Error(err), zap.String("realIncome", realIncome))
				continue
			}
		}

		if revenueDate, ok := v["revenueDate"]; ok {
			tr.RevenueDate, err = strconv.ParseInt(revenueDate, 10, 64)
			if err != nil {
				appzaplog.Error("ParseFloat err", zap.Error(err), zap.String("revenueDate", revenueDate))
				continue
			}
		}
		r = append(r, tr)
	}
	appzaplog.Debug("[+]QueryRevenueRecord", zap.String("starttime", timeGreaterThan), zap.String("endtime", timeLessThan),
		zap.Uint64("hostuid", anchorUid), zap.Int64("uid", uid),
		zap.Int64("totalPageSize", totalPageSize), zap.Any("resp", r))
	return
}

func UserAccountByUidAndType(ctx context.Context, uid int64, appid turnover.TAppId, currencyType turnover.TCurrencyType) (*turnover.TUserAccount, error) {
	connection, err := turnOverThrift.GetConnection()
	if err != nil {
		appzaplog.Error("UserAccountByUidAndType GetConnection err", zap.Error(err))
		return nil, err
	}
	protocolFactory := thrift.NewTCompactProtocolFactory()
	client := turnover.NewTUserAccountServiceClientFactory(connection.Transport, protocolFactory)
	account, err := client.GetUserAccountByUidAndType(ctx, uid, appid, currencyType)
	if err != nil {
		turnOverThrift.ReportErrorConnection(connection)
		return nil, err
	}
	turnOverThrift.ReturnConnection(connection)
	return account, nil
}

/*
 * 查看ow账户的可提现金额
 * currencyType: Bilin_Profit 10000收益币=1元
 */
func ChannelAccountByUidAndType(ctx context.Context, uid int64, sid int64, appid turnover.TAppId, currencyType turnover.TCurrencyType) (*turnover.TChannelAccount, error) {
	connection, err := turnOverThrift.GetConnection()
	if err != nil {
		appzaplog.Error("ChannelAccountByUidAndType GetConnection err", zap.Error(err))
		return nil, err
	}
	protocolFactory := thrift.NewTCompactProtocolFactory()
	client := turnover.NewTCurrencyServiceClientFactory(connection.Transport, protocolFactory)
	info, err := client.GetChannelAccountByUidAndType(ctx, uid, sid, appid, currencyType)
	if err != nil {
		turnOverThrift.ReportErrorConnection(connection)
		return nil, err
	}
	turnOverThrift.ReturnConnection(connection)
	return info, err
}

func QueryRevenueRecordPagingProxy(ctx context.Context, uid int64, appid turnover.TAppId, pagesize int32, page int32, revenueUserType turnover.TUserType, startDate int64, endDate int64, anchorUid int64, sid int64, srcType turnover.TRevenueSrcType) (*turnover.TQueryPageInfo, error) {
	connection, err := turnOverThrift.GetConnection()
	if err != nil {
		appzaplog.Error("QueryRevenueRecordPagingProxy GetConnection err", zap.Error(err))
		return nil, err
	}
	protocolFactory := thrift.NewTCompactProtocolFactory()
	client := turnover.NewTCurrencyServiceClientFactory(connection.Transport, protocolFactory)
	info, err := client.QueryRevenueRecordPaging(ctx, uid, appid, pagesize, page, revenueUserType, startDate, endDate, anchorUid, sid, srcType)
	if err != nil {
		turnOverThrift.ReportErrorConnection(connection)
		return nil, err
	}
	turnOverThrift.ReturnConnection(connection)
	return info, err
}

func QueryRevenueRecordPaging(ctx context.Context, uid int64, pagesize int32, pagenum int32,
	revenueUserType turnover.TUserType,
	startTime int64, endTime int64,
	anchorUid int64, sid int64) (*turnover.TQueryPageInfo, error) {
	return QueryRevenueRecordPagingProxy(
		ctx, uid, turnover.TAppId_Bilin,
		pagesize, pagenum,
		revenueUserType,
		startTime, endTime,
		anchorUid, sid,
		turnover.TRevenueSrcType_Props)
}

const TimeLayoutOthers string = "2006-01-02 15:04:05"

//func TimeStrToUnixOther(timeStr string) (time.Time, error) {
//	return time.ParseInLocation(TimeLayoutOthers, timeStr, time.Local)
//}

const TimeLayoutOther string = "20060102150405"

func unixToStringTime(t time.Time) string {
	return t.Format(TimeLayoutOther)
}
