package handler

import (
	"bilin/clientcenter"
	"bilin/official/dao"
	"bilin/protocol/userinfocenter"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"context"
	"github.com/gin-gonic/gin"
)

type Contract struct {
	GuildID              uint64 `json:"guild_id"`
	HostUid              uint64 `json:"host_uid"`
	BILINNumber          uint64 `json:"bilin_number"` // 比邻号
	NickName             string `json:"nick_name"`
	ContractStartTime    int64  `json:"contract_start_time"`
	ContractEndTime      int64  `json:"contract_end_time"`
	GuildSharePercentage uint64 `json:"guild_share_percentage"`
	HostSharePercentage  uint64 `json:"host_share_percentage"`
}

func GetHostContract(c *gin.Context) *HttpError {
	ret := successHttp
	var (
		contract *dao.Contract
		userinfo map[uint64]*userinfocenter.UserInfo
		blidmap  *userinfocenter.BatchUserBiLinIdResp
	)
	for {
		cookieuid := c.GetInt64("uid")
		idInt, err := UInt64Param("id", c)
		if err != nil {
			appzaplog.Error("GetHostContract UInt64Param err", zap.Error(err))
			ret = parseURLHttpErr
			break
		}
		if !cookieUserEqReqUser(cookieuid, idInt) {
			appzaplog.Warn("GetHostGuild not host", zap.Uint64("requid", idInt), zap.Int64("cookieuid", cookieuid))
			ret = authHttpErr
			break
		}
		contract, err = dao.GetContractByHostUid(idInt)
		if err != nil {
			appzaplog.Error("GetHostContract GetContractByHostUid err", zap.Error(err))
			ret = daoGetHttpErr
			break
		}

		if contract == nil {
			appzaplog.Debug("no contract", zap.Uint64("uid", idInt))
			break
		}
		// 批量获取用户bilin信息
		blidmap, err = clientcenter.UserInfoClient().BatchUserBiLinId(context.TODO(),
			&userinfocenter.BatchUserBiLinIdReq{
				Uid: []uint64{idInt},
			})
		if err != nil {
			appzaplog.Error("[+]takeRecords BatchUserBiLinId err", zap.Error(err), zap.Int64("cookieuid", cookieuid))
			ret = daoGetHttpErr
			break
		}

		// 批量获取用户昵称等信息
		userinfo, err = clientcenter.TakeUserInfo([]uint64{idInt})
		if err != nil {
			appzaplog.Error("[+]takeRecords TakeUserInfo err", zap.Error(err), zap.Int64("cookieuid", cookieuid))
			ret = daoGetHttpErr
			break
		}

		break
	}

	jsonret := &HttpRetComm{
		Code: ret.code,
		Desc: ret.desc,
	}
	if ret.code == 0 && contract != nil {
		contract := &Contract{
			GuildID:              contract.GuildID,
			HostUid:              contract.HostUid,
			ContractStartTime:    contract.ContractStartTime.Unix() * 1000,
			ContractEndTime:      contract.ContractEndTime.Unix() * 1000,
			GuildSharePercentage: contract.GuildSharePercentage,
			HostSharePercentage:  contract.HostSharePercentage,
		}
		if info, ok := userinfo[contract.HostUid]; ok {
			contract.NickName = info.NickName
		}
		if blid, ok := blidmap.Uid2Bilinid[contract.HostUid]; ok {
			contract.BILINNumber = blid
		}
		jsonret.Data = contract
	}
	c.JSON(200, jsonret)
	return ret
}
