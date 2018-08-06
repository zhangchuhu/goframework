package handler

import (
	"bilin/clientcenter"
	"bilin/official/dao"
	"bilin/protocol/userinfocenter"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"context"
	"github.com/gin-gonic/gin"
	"strconv"
)

type ContractRecord struct {
	BILINNumber          uint64 `json:"bilin_number"` // 比邻号
	NickName             string `json:"nick_name"`
	ContractStartTime    int64  `json:"contract_start_time"`
	ContractEndTime      int64  `json:"contract_end_time"`
	GuildSharePercentage uint64 `json:"guild_share_percentage"`
	HostSharePercentage  uint64 `json:"host_share_percentage"`
}

type ContractRecords struct {
	Records []ContractRecord `json:"records"`
}

func GetGuildContractRecords(c *gin.Context) *HttpError {
	ret := successHttp
	var (
		records []ContractRecord
	)
	for {
		cookieuid := c.GetInt64("uid")
		guildid, err := UInt64Param("id", c)
		if err != nil {
			appzaplog.Error("GetGuildContractRecords UInt64Param err", zap.Error(err))
			ret = parseURLHttpErr
			break
		}
		if !owQueryGuild(cookieuid, guildid) {
			appzaplog.Warn("GetGuildContractRecords not ow", zap.Uint64("requid", guildid), zap.Int64("cookieuid", cookieuid))
			ret = authHttpErr
			break
		}

		blid := c.Query("hostBiLinNumber")
		var blidU64 uint64
		if blid != "" {
			blidU64, err = strconv.ParseUint(blid, 10, 64)
			if err != nil {
				appzaplog.Error("GetGuildContractRecords ParseUint err", zap.Error(err), zap.String("blid", blid))
				ret = daoGetHttpErr
				break
			}
		}

		records, err = takeRecords(guildid, blidU64)
		if err != nil {
			appzaplog.Error("GetGuildContractRecords takeRecords err", zap.Error(err))
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
		jsonret.Data = &ContractRecords{
			Records: records,
		}
	}
	c.JSON(200, jsonret)
	return ret
}

func takeRecords(guildid uint64, bid uint64) ([]ContractRecord, error) {
	var (
		records []ContractRecord
	)
	contractinfo, err := dao.GetContractsByGuildID(guildid)
	if err != nil {
		appzaplog.Error("[+]takeRecords GetContractsByGuildID err", zap.Error(err), zap.Uint64("guildid", guildid))
		return nil, err
	}

	var (
		uids []uint64
	)
	for _, v := range contractinfo {
		uids = append(uids, v.HostUid)
	}

	// 批量获取用户bilin信息
	blidmap, err := clientcenter.UserInfoClient().BatchUserBiLinId(context.TODO(),
		&userinfocenter.BatchUserBiLinIdReq{
			Uid: uids,
		})
	if err != nil {
		appzaplog.Error("[+]takeRecords BatchUserBiLinId err", zap.Error(err), zap.Uint64("guildid", guildid))
		return records, err
	}

	// 批量获取用户昵称等信息
	userinfo, err := clientcenter.TakeUserInfo(uids)
	if err != nil {
		appzaplog.Error("[+]takeRecords TakeUserInfo err", zap.Error(err), zap.Uint64("guildid", guildid))
		return records, err
	}

	for _, v := range contractinfo {
		record := ContractRecord{
			ContractStartTime:    v.ContractStartTime.Unix() * 1000,
			ContractEndTime:      v.ContractEndTime.Unix() * 1000,
			HostSharePercentage:  v.HostSharePercentage,
			GuildSharePercentage: v.GuildSharePercentage,
		}
		if nick, ok := userinfo[v.HostUid]; ok {
			record.NickName = nick.NickName
		}
		if blnum, ok := blidmap.Uid2Bilinid[v.HostUid]; ok {
			record.BILINNumber = blnum
		}
		if bid > 0 && bid != record.BILINNumber {
			continue
		}
		records = append(records, record)
	}
	return records, nil
}
