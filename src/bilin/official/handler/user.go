package handler

import (
	"bilin/clientcenter"
	"bilin/official/service"
	"bilin/protocol/userinfocenter"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"context"
	"github.com/gin-gonic/gin"
)

type HostInfo struct {
	UID           uint64 `json:"uid"`
	BILINNumber   uint64 `json:"bilin_number"` // 比邻号
	NickName      string `json:"nick_name"`
	Avatar        string `json:"avatar"`
	TotalHeartNum uint64 `json:"total_heart_num"` // 累计心值
	TotalCharmNum uint64 `json:"total_charm_num"` // 魅力值
	FansNum       uint64 `json:"fans_num"`        // 粉丝数
	AttentionNum  uint64 `json:"attention_num"`   // 关注数目
}

func GetHost(c *gin.Context) *HttpError {
	ret := successHttp
	var (
		hostinfo *HostInfo
	)
	for {
		cookieuid := c.GetInt64("uid")
		hostid, err := UInt64Param("id", c)
		if err != nil {
			appzaplog.Error("GetHost UInt64Param err", zap.Error(err))
			ret = parseURLHttpErr
			break
		}

		if !cookieUserEqReqUser(cookieuid, hostid) {
			appzaplog.Warn("GetHostGuild not host", zap.Uint64("requid", hostid), zap.Int64("cookieuid", cookieuid))
			ret = authHttpErr
			break
		}

		hostinfo, err = hostInfo(hostid)
		if err != nil {
			appzaplog.Error("GetHost hostInfo err", zap.Error(err))
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
		jsonret.Data = hostinfo
	}
	c.JSON(200, jsonret)
	return ret
}

func hostInfo(uid uint64) (*HostInfo, error) {
	var hinfo *HostInfo

	// 批量用户信息
	uinfo, err := clientcenter.TakeUserInfo([]uint64{uid})
	if err != nil {
		appzaplog.Error("GetHost hostInfo err", zap.Error(err), zap.Uint64("uid", uid))
		return hinfo, err
	}

	// 批量获取用户bilin信息
	blidmap, err := clientcenter.UserInfoClient().BatchUserBiLinId(context.TODO(),
		&userinfocenter.BatchUserBiLinIdReq{
			Uid: []uint64{uid},
		})
	if err != nil {
		appzaplog.Error("[+]hostInfo BatchUserBiLinId err", zap.Error(err), zap.Uint64("uid", uid))
		return hinfo, err
	}

	// 用户关注，魅力，粉丝
	attention, err := clientcenter.UserInfoClient().AttentionInfo(context.TODO(),
		&userinfocenter.AttentionInfoReq{
			Uid: uid,
		})
	if err != nil {
		appzaplog.Error("[+]hostInfo AttentionInfo err", zap.Error(err), zap.Uint64("uid", uid))
		return hinfo, err
	}

	heartCount, err := service.BilinCumulativeProfit(context.TODO(), uid)
	if err != nil {
		appzaplog.Error("[+]hostInfo BilinCumulativeProfit err", zap.Error(err), zap.Uint64("uid", uid))
		return hinfo, err
	}
	hinfo = &HostInfo{
		UID:           uid,
		AttentionNum:  attention.Attentionnum,
		FansNum:       attention.Fansnum,
		TotalCharmNum: attention.Glamour,
		TotalHeartNum: heartCount,
	}
	if info, ok := uinfo[uid]; ok {
		hinfo.NickName = info.NickName
		hinfo.Avatar = info.Avatar
	}

	if blnum, ok := blidmap.Uid2Bilinid[uid]; ok {
		hinfo.BILINNumber = blnum
	}
	return hinfo, nil
}
