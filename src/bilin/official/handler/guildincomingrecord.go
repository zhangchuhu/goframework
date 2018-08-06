package handler

import (
	"bilin/clientcenter"
	"bilin/official/service"
	"bilin/protocol/userinfocenter"
	"bilin/thrift/gen-go/turnover"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

type GuildIncomingRecord struct {
	HostBilinID     int64   `json:"host_bilin_id"`
	HostNickName    string  `json:"host_nick_name"`
	HeartNum        float64 `json:"heart_num"`        // 收礼心值
	GuildPercentage int64   `json:"guild_percentage"` // 工会抽成比例
	HOSTIncoming    float64 `json:"host_incoming"`    // 主播收入
	GuildIncoming   float64 `json:"guild_incoming"`   // 工会抽成收入
	RevenueDate     string  `json:"revenue_date"`     //消费的日期
}

type GuildIncomingRecordS struct {
	TotalPagesize      int64                 `json:"total_pagesize"`
	CurMonthIncoming   float64               `json:"cur_month_incoming"` // 本月佣金
	CashIncoming       int64                 `json:"cash_incoming"`      // 可提现佣金
	Records            []GuildIncomingRecord `json:"records"`
	TotalHeartNum      float64               `json:"total_heart_num"` // 总收礼心值
	TotalHOSTIncoming  float64               `json:"total_host_incoming"`
	TotalGuildIncoming float64               `json:"total_guild_incoming"` // 抽成总收入
}

func GetGuildIncomingRecords(c *gin.Context) *HttpError {
	ret := successHttp
	var (
		data *GuildIncomingRecordS
	)
	for {
		cookieuid := c.GetInt64("uid")
		guildIdInt, err := UInt64Param("id", c)
		if err != nil {
			appzaplog.Error("GetGuildIncomingRecords UInt64Param err", zap.Error(err))
			ret = parseURLHttpErr
			break
		}
		if !owQueryGuild(cookieuid, guildIdInt) {
			appzaplog.Warn("GetGuildIncomingRecords not ow", zap.Uint64("requid", guildIdInt), zap.Int64("cookieuid", cookieuid))
			ret = authHttpErr
			break
		}
		account, err := service.ChannelAccountByUidAndType(context.TODO(), c.GetInt64("uid"), int64(guildIdInt), turnover.TAppId_Bilin, turnover.TCurrencyType_Bilin_Profit)
		if err != nil {
			appzaplog.Error("GetGuildIncomingRecords UserAccountByUidAndType err", zap.Error(err), zap.Uint64("uid", guildIdInt))
			ret = parseURLHttpErr
			break
		}
		data, err = guildIncomingRecords(int64(guildIdInt), c)
		if err != nil {
			appzaplog.Error("GetGuildIncomingRecords guildIncomingRecords err", zap.Error(err), zap.Uint64("uid", guildIdInt))
			ret = parseURLHttpErr
			break
		}
		if account != nil {
			data.CashIncoming = account.Amount / 10000
		}

		break
	}

	jsonret := &HttpRetComm{
		Code: ret.code,
		Desc: ret.desc,
	}
	if ret.code == 0 {
		jsonret.Data = data
	}
	appzaplog.Debug("[-]GetGuildIncomingRecords success", zap.Any("resp", jsonret))
	c.JSON(200, jsonret)
	return ret
}

func hostUid(c *gin.Context) uint64 {
	var (
		hostuid uint64
	)

	hostBiLinNumber := c.Query("hostBiLinNumber")
	if hostBiLinNumber != "" {
		blnumInt, err := strconv.ParseUint(hostBiLinNumber, 10, 64)
		if err != nil {
			return hostuid
		}
		resp, err := clientcenter.UserInfoClient().BatchUserIdByBiLinId(context.TODO(), &userinfocenter.BatchUserIdByBiLinIdReq{
			Bilinid: []uint64{blnumInt},
		})
		if err != nil {
			appzaplog.Error("BatchUserIdByBiLinId err", zap.Error(err), zap.Uint64("bilinid", blnumInt))
			return hostuid
		}
		if myuid, ok := resp.Bilinid2Uid[blnumInt]; ok {
			hostuid = myuid
		}
	}
	return hostuid
}

func guildIncomingRecords(guildid int64, c *gin.Context) (*GuildIncomingRecordS, error) {
	irecords := &GuildIncomingRecordS{}
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

	hostuid := hostUid(c)

	revenue, totalPageSize, err := service.QueryRevenueRecord(context.TODO(), c.GetInt64("uid"), guildid,
		startTime, endTime, 2, int32(pageInt), int32(pageSizeInt), hostuid,
	)
	if err != nil {
		appzaplog.Error("QueryRevenueRecord err", zap.Error(err), zap.Int64("guildid", guildid))
		return irecords, err
	}
	var hostids []uint64
	for _, v := range revenue {
		if v.ContributeUid > 0 {
			hostids = append(hostids, uint64(v.ContributeUid))
		}
	}

	uinfo, err := clientcenter.TakeUserInfo(hostids)
	if err != nil {
		appzaplog.Error("TakeUserInfo err", zap.Error(err), zap.Int64("guildid", guildid))
		return irecords, err
	}

	bilinid, err := clientcenter.UserInfoClient().BatchUserBiLinId(context.TODO(), &userinfocenter.BatchUserBiLinIdReq{
		Uid: hostids,
	})
	if err != nil {
		appzaplog.Error("BatchUserBiLinId err", zap.Error(err), zap.Int64("guildid", guildid))
		return irecords, err
	}

	irecords.TotalPagesize = totalPageSize
	for _, v := range revenue {
		record := GuildIncomingRecord{
			HeartNum:        v.Income,
			GuildPercentage: int64(v.IncomeRate),
			HOSTIncoming:    (v.Income - v.RealIncome) / 10000,
			GuildIncoming:   v.RealIncome / 10000,
			RevenueDate:     time.Unix(v.RevenueDate/1000, 0).Format("2006-01-02"),
		}
		if blid, ok := bilinid.Uid2Bilinid[uint64(v.ContributeUid)]; ok {
			record.HostBilinID = int64(blid)
		}
		if nick, ok := uinfo[uint64(v.ContributeUid)]; ok {
			record.HostNickName = nick.NickName
		}
		irecords.Records = append(irecords.Records, record)
	}

	// 获取全部的信息
	totalrevenue, _, err := service.QueryRevenueRecord(context.TODO(), c.GetInt64("uid"), guildid,
		startTime, endTime, 2, int32(1), int32(totalPageSize), hostuid,
	)
	if err != nil {
		appzaplog.Error("QueryRevenueRecord err", zap.Error(err), zap.Int64("guildid", guildid))
		return irecords, err
	}
	for _, v := range totalrevenue {
		record := GuildIncomingRecord{
			HeartNum:        v.Income,
			GuildPercentage: int64(v.IncomeRate),
			HOSTIncoming:    (v.Income - v.RealIncome) / 10000,
			GuildIncoming:   v.RealIncome / 10000,
			RevenueDate:     time.Unix(v.RevenueDate/1000, 0).Format("2006-01-02"),
		}
		irecords.TotalHeartNum += record.HeartNum
		irecords.TotalHOSTIncoming += record.HOSTIncoming
		irecords.TotalGuildIncoming += record.GuildIncoming
	}

	info, err := service.QueryCurMonthRevenueRecord(context.TODO(), guildid, guildid,
		2, 0)
	if err != nil {
		return irecords, err
	}
	for _, v := range info {
		irecords.CurMonthIncoming += v.RealIncome / 10000
	}
	return irecords, nil
}
