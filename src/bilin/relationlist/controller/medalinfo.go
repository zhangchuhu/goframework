package controller

import (
	"bilin/relationlist/entity"
	"bilin/relationlist/service"
	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
)

var (
	medalsConfig map[int32]entity.MedalInfo
)

//从数据库中后取勋章相关的配置信息,启动时初始化
func InitMedalConfig() {
	var err error
	medalsConfig, err = service.MysqlGetMedalsConfig()
	if err != nil {
		panic("MysqlGetMedalsConfig error!")
	}
	log.Info("InitMedalConfig end!", zap.Any("medalsConfig", medalsConfig))
}

func GetMedalConfig() map[int32]entity.MedalInfo {
	return medalsConfig
}

func GetUserMedalInfo(owner uint64, guest_uid uint64) (rank uint32, name string, medalUrl string) {
	const prefix = "GetUserMedalInfo "

	//先从缓存中查，查不到再查db,勋章数据存5分钟
	medalsRet, err := service.RedisGetOwnerMedals(owner)
	if err != nil || len(medalsRet) == 0 { // 直接从db中获取
		medalsRet, err = service.MysqlGetOwnerMedalsInfo(owner)
		service.RedisSetOwnerMedals(owner, medalsRet)
	}

	medalID, exist := medalsRet[guest_uid]
	if !exist {
		log.Info(prefix+"user have no medal with owner ", zap.Any("owner", owner), zap.Any("guest_uid", guest_uid))
		return 0, "", ""
	}

	rank = uint32(medalID)
	name = medalsConfig[medalID].MedalName
	medalUrl = medalsConfig[medalID].MedalUrl

	log.Info(prefix+"end", zap.Any("owner", owner), zap.Any("guest_uid", guest_uid), zap.Any("rank", rank), zap.Any("name", name))
	return
}
