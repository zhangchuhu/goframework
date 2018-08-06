package handler

import (
	"bilin/official/dao"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
)

func cookieUserEqReqUser(realuid int64, requid uint64) bool {
	return uint64(realuid) == requid
}

func accessHostRecords(realuid int64, requestuid uint64) bool {
	return cookieUserEqReqUser(realuid, requestuid) || owQueryHost(uint64(realuid), requestuid)
}

func owQueryHost(ow uint64, host uint64) bool {
	contract, err := dao.GetContractByHostUid(host)
	if err != nil {
		appzaplog.Error("owQueryHost GetContractByHostUid err", zap.Error(err),
			zap.Uint64("ow", ow),
			zap.Uint64("host", host),
		)
		return false
	}
	if contract == nil {
		appzaplog.Debug("not contract", zap.Uint64("ow", ow), zap.Uint64("host", host))
		return false
	}
	return owQueryGuild(int64(ow), contract.GuildID)
}

func owQueryGuild(ow int64, guildid uint64) bool {
	guild, err := dao.GetByGuildID(guildid)
	if err != nil {
		return false
	}
	return guild.OW == uint64(ow)
}
