package handler

import (
	"bilin/official/service"
	"bilin/thrift/gen-go/turnover"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

type IncomingRecord struct {
	DateTime        string  `json:"date_time"`        // yyyy-mm-dd
	HeartNum        float64 `json:"heart_num"`        // 收礼心值
	GuildPercentage int64   `json:"guild_percentage"` // 工会抽成比例
	HOSTIncoming    float64 `json:"host_incoming"`    // 主播收入
	GuildIncoming   float64 `json:"guild_incoming"`   // 工会抽成收入
}

type IncomingRecordS struct {
	TotalPagesize     int64            `json:"total_pagesize"`
	CurMonthIncoming  float64          `json:"cur_month_incoming"` // 主播本月佣金
	CashIncoming      int64            `json:"cash_incoming"`      // 主播可提现佣金
	Records           []IncomingRecord `json:"records"`
	TotalHeartNum     float64          `json:"total_heart_num"`     // 总收礼心值
	TotalHostIncoming float64          `json:"total_host_incoming"` // 主播抽成总收入
}

func GetHostIncomingRecords(c *gin.Context) *HttpError {
	appzaplog.Debug("[+]GetHostIncomingRecords")
	ret := successHttp
	var (
		data *IncomingRecordS
	)
	for {
		cookieuid := c.GetInt64("uid")
		idInt, err := UInt64Param("id", c)
		if err != nil {
			appzaplog.Error("GetHostIncomingRecords UInt64Param err", zap.Error(err))
			ret = parseURLHttpErr
			break
		}
		if !accessHostRecords(cookieuid, idInt) {
			appzaplog.Warn("GetHostIncomingRecords not host", zap.Uint64("requid", idInt), zap.Int64("cookieuid", cookieuid))
			ret = authHttpErr
			break
		}
		account, err := service.UserAccountByUidAndType(context.TODO(), int64(idInt), turnover.TAppId_Bilin, turnover.TCurrencyType_Bilin_Profit)
		if err != nil {
			appzaplog.Error("GetHostIncomingRecords UserAccountByUidAndType err", zap.Error(err), zap.Uint64("uid", idInt))
			ret = parseURLHttpErr
			break
		}

		//todo
		data, err = hostIncomingRecords(idInt, c)
		if err != nil {
			appzaplog.Error("GetHostIncomingRecords hostIncomingRecords err", zap.Error(err), zap.Uint64("uid", idInt))
			ret = parseURLHttpErr
			break
		}
		data.CashIncoming = account.Amount / 10000
		break
	}

	jsonret := &HttpRetComm{
		Code: ret.code,
		Desc: ret.desc,
	}
	if ret.code == 0 {
		jsonret.Data = data
	}
	c.JSON(200, jsonret)
	appzaplog.Debug("[-]GetHostIncomingRecords", zap.Any("resp", jsonret))
	return ret
}

func hostIncomingRecords(hostid uint64, c *gin.Context) (*IncomingRecordS, error) {
	irecords := &IncomingRecordS{}
	startTime := c.Query("startTime")

	endTime := c.Query("endTime")

	page := c.Query("pageNum")
	pageInt, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		return irecords, fmt.Errorf("parese pageNum err")
	}
	pagesize := c.Query("pageSize")
	pageSizeInt, err := strconv.ParseInt(pagesize, 10, 64)
	if err != nil {
		return irecords, fmt.Errorf("parese pagesize err")
	}
	revenue, totalPageSize, err := service.QueryRevenueRecord(context.TODO(), int64(hostid), 0,
		startTime, endTime, 1, int32(pageInt), int32(pageSizeInt), hostid,
	)
	if err != nil {
		appzaplog.Error("QueryRevenueRecord err", zap.Error(err), zap.Uint64("hostid", hostid))
		return irecords, err
	}
	irecords.TotalPagesize = totalPageSize
	for _, v := range revenue {
		record := IncomingRecord{
			DateTime:        time.Unix(v.RevenueDate/1000, 0).Format("2006-01-02"),
			HeartNum:        v.Income,
			GuildPercentage: int64(100 - v.IncomeRate),
			HOSTIncoming:    v.RealIncome / 10000,
			GuildIncoming:   (v.Income - v.RealIncome) / 10000,
		}
		irecords.Records = append(irecords.Records, record)
		//irecords.TotalHeartNum += record.HeartNum
		//irecords.TotalHostIncoming += record.HOSTIncoming
	}

	//为了获取全部的信息，只能再取一次了
	totalrevenue, _, err := service.QueryRevenueRecord(context.TODO(), int64(hostid), 0,
		startTime, endTime, 1, int32(1), int32(totalPageSize), hostid,
	)
	if err != nil {
		appzaplog.Error("QueryRevenueRecord err", zap.Error(err), zap.Uint64("hostid", hostid))
		return irecords, err
	}
	for _, v := range totalrevenue {
		record := IncomingRecord{
			DateTime:        time.Unix(v.RevenueDate/1000, 0).Format("2006-01-02"),
			HeartNum:        v.Income,
			GuildPercentage: int64(100 - v.IncomeRate),
			HOSTIncoming:    v.RealIncome / 10000,
			GuildIncoming:   (v.Income - v.RealIncome) / 10000,
		}
		irecords.TotalHeartNum += record.HeartNum
		irecords.TotalHostIncoming += record.HOSTIncoming
	}

	info, err := service.QueryCurMonthRevenueRecord(context.TODO(), int64(hostid), 0,
		1, hostid)
	if err != nil {
		appzaplog.Error("QueryCurMonthRevenueRecord err", zap.Error(err), zap.Uint64("hostid", hostid))
		return irecords, err
	}
	for _, v := range info {
		irecords.CurMonthIncoming += v.RealIncome / 10000
	}
	return irecords, nil
}
