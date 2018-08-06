package handler

import (
	"bilin/clientcenter"
	"bilin/official/service"
	"bilin/protocol/userinfocenter"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

type GuildIncomingDetail struct {
	DateTime              string  `json:"date_time"` // yyyy-mm-dd hh:mm
	HostBilinNumber       uint64  `json:"host_bilin_number"`
	HostNickName          string  `json:"host_nick_name"`
	ContributeBilinNumber uint64  `json:"contribute_bilin_number"`
	ContributeNickName    string  `json:"contribute_nick_name"`
	PropName              string  `json:"prop_name"`
	TotalPropValue        float64 `json:"total_prop_value"`
	PropNum               int32   `json:"prop_num"`
}
type GuildIncomingDetailS struct {
	TotalPagesize int32                 `json:"total_pagesize"`
	Records       []GuildIncomingDetail `json:"records"`
}

func GetGuildIncomingDetail(c *gin.Context) *HttpError {
	ret := successHttp
	var gincomingdetail *GuildIncomingDetailS
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

		gincomingdetail, err = guildIncomingDetail(uint64(guildIdInt), c)
		if err != nil {
			appzaplog.Error("guildIncomingDetail not host", zap.Uint64("guildid", guildIdInt), zap.Int64("cookieuid", cookieuid))
			ret = daoGetHttpErr
			break
		}
		break
	}
	jsonret := &HttpRetComm{
		Code: ret.code,
		Desc: ret.desc,
	}
	if ret.code == 0 {
		jsonret.Data = gincomingdetail
	}
	c.JSON(200, jsonret)
	return ret
}

func guildIncomingDetail(guildid uint64, c *gin.Context) (*GuildIncomingDetailS, error) {
	irecords := &GuildIncomingDetailS{}
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

	revenue, err := service.QueryChannelWeekPropsRecieve(context.TODO(), int64(guildid), startTime,
		endTime, int32(pageInt), int32(pageSizeInt), int64(hostuid))
	if err != nil {
		appzaplog.Error("QueryRevenueRecord err", zap.Error(err), zap.Uint64("hostid", hostuid))
		return irecords, err
	}

	var hostids []uint64
	for _, v := range revenue.Content {
		if v.UID > 0 {
			hostids = append(hostids, uint64(v.UID))
		}
		if v.AnchorUid > 0 {
			hostids = append(hostids, uint64(v.AnchorUid))
		}
	}
	uinfo, err := clientcenter.TakeUserInfo(hostids)
	if err != nil {
		appzaplog.Error("TakeUserInfo err", zap.Error(err))
		return irecords, err
	}

	bilinid, err := clientcenter.UserInfoClient().BatchUserBiLinId(context.TODO(), &userinfocenter.BatchUserBiLinIdReq{
		Uid: hostids,
	})
	if err != nil {
		appzaplog.Error("BatchUserBiLinId err", zap.Error(err))
		return irecords, err
	}

	irecords.TotalPagesize = revenue.TotalElement
	for _, v := range revenue.Content {
		record := GuildIncomingDetail{
			DateTime:       time.Unix(v.UsedTime/1000, 0).Format("2006-01-02 15:04"),
			PropName:       v.PropName,
			TotalPropValue: v.Amount,
			PropNum:        v.PropCnt,
		}
		if blid, ok := bilinid.Uid2Bilinid[uint64(v.UID)]; ok {
			record.ContributeBilinNumber = blid
		}
		if blid, ok := bilinid.Uid2Bilinid[uint64(v.AnchorUid)]; ok {
			record.HostBilinNumber = blid
		}
		if nick, ok := uinfo[uint64(v.UID)]; ok {
			record.ContributeNickName = nick.NickName
		}
		if nick, ok := uinfo[uint64(v.AnchorUid)]; ok {
			record.HostNickName = nick.NickName
		}
		irecords.Records = append(irecords.Records, record)
	}

	return irecords, nil
}
