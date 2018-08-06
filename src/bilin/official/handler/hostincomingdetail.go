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

type HostIncomingDetail struct {
	DateTime              string  `json:"date_time"` // yyyy-mm-dd hh:mm
	ContributeBilinNumber uint64  `json:"contribute_bilin_number"`
	ContributeNickName    string  `json:"contribute_nick_name"`
	PropName              string  `json:"prop_name"`
	TotalPropValue        float64 `json:"total_prop_value"`
	PropNum               int32   `json:"prop_num"`
}

type HostIncomingDetailS struct {
	TotalPagesize int32                `json:"total_pagesize"`
	Records       []HostIncomingDetail `json:"records"`
}

func GetHostIncomingDetail(c *gin.Context) *HttpError {
	ret := successHttp
	var hincomingdetail *HostIncomingDetailS
	appzaplog.Debug("[+]GetHostIncomingDetail")
	for {
		cookieuid := c.GetInt64("uid")
		idInt, err := UInt64Param("id", c)
		if err != nil {
			appzaplog.Error("GetHostIncomingDetail UInt64Param err", zap.Error(err))
			ret = parseURLHttpErr
			break
		}
		if !accessHostRecords(cookieuid, idInt) {
			appzaplog.Error("GetHostIncomingDetail not host", zap.Uint64("requid", idInt), zap.Int64("cookieuid", cookieuid))
			ret = authHttpErr
			break
		}

		hincomingdetail, err = hostIncomingDetail(uint64(idInt), c)
		if err != nil {
			appzaplog.Error("GetHostIncomingDetail err", zap.Error(err), zap.Uint64("requid", idInt), zap.Int64("cookieuid", cookieuid))
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
		jsonret.Data = hincomingdetail
	}
	c.JSON(200, jsonret)
	appzaplog.Debug("[-]GetHostIncomingDetail", zap.Any("resp", jsonret))
	return ret
}

func hostIncomingDetail(hostid uint64, c *gin.Context) (*HostIncomingDetailS, error) {
	appzaplog.Debug("[+]hostIncomingDetail", zap.Uint64("uid", hostid))
	irecords := &HostIncomingDetailS{}
	startTime := c.Query("startTime")

	endTime := c.Query("endTime")

	page := c.Query("pageNum")
	pageInt, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		appzaplog.Error("hostIncomingDetail ParseInt pageNum err", zap.Error(err), zap.Uint64("uid", hostid))
		return irecords, fmt.Errorf("parese pageNum err")
	}
	pagesize := c.Query("pageSize")
	pageSizeInt, err := strconv.ParseInt(pagesize, 10, 64)
	if err != nil {
		appzaplog.Error("hostIncomingDetail ParseInt pageSize err", zap.Error(err), zap.Uint64("uid", hostid))
		return irecords, fmt.Errorf("parese pagesize err")
	}

	revenue, err := service.QueryAnchorWeekPropsRecieve(context.TODO(), int64(hostid), startTime,
		endTime, int32(pageInt), int32(pageSizeInt))
	if err != nil {
		appzaplog.Error("hostIncomingDetail QueryRevenueRecord err", zap.Error(err), zap.Uint64("hostid", hostid))
		return irecords, err
	}
	var hostids []uint64
	for _, v := range revenue.Content {
		if v.UID > 0 {
			hostids = append(hostids, uint64(v.UID))
		}
	}
	uinfo, err := clientcenter.TakeUserInfo(hostids)
	if err != nil {
		appzaplog.Error("hostIncomingDetail TakeUserInfo err", zap.Error(err))
		return irecords, err
	}

	bilinid, err := clientcenter.UserInfoClient().BatchUserBiLinId(context.TODO(), &userinfocenter.BatchUserBiLinIdReq{
		Uid: hostids,
	})
	if err != nil {
		appzaplog.Error("hostIncomingDetail BatchUserBiLinId err", zap.Error(err))
		return irecords, err
	}

	irecords.TotalPagesize = revenue.TotalElement
	for _, v := range revenue.Content {
		record := HostIncomingDetail{
			DateTime:       time.Unix(v.UsedTime/1000, 0).Format("2006-01-02 15:04"),
			PropName:       v.PropName,
			TotalPropValue: v.Amount,
			PropNum:        v.PropCnt,
		}
		if blid, ok := bilinid.Uid2Bilinid[uint64(v.UID)]; ok {
			record.ContributeBilinNumber = blid
		}
		if nick, ok := uinfo[uint64(v.UID)]; ok {
			record.ContributeNickName = nick.NickName
		}
		irecords.Records = append(irecords.Records, record)
	}
	appzaplog.Debug("[-]hostIncomingDetail", zap.Uint64("uid", hostid), zap.Any("resp", irecords))
	return irecords, nil
}
