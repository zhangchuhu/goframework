package controller

import (
	"bilin/relationlist/service"
	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
)

func UserOnMike(owner uint64, guest_uid uint64, onMikeTime int64) (err error) {
	const prefix = "UserOnMike "

	// 更新redis
	service.RedisUserOnMike(owner, guest_uid, onMikeTime)

	log.Info(prefix, zap.Any("owner", owner), zap.Any("guest_uid", guest_uid), zap.Any("onMikeTime", onMikeTime))
	return
}

func UserOffMike(owner uint64, guest_uid uint64) (err error) {
	const prefix = "UserOffMike "

	// 更新redis
	service.RedisUserOffMike(owner, guest_uid)

	log.Info(prefix, zap.Any("owner", owner), zap.Any("guest_uid", guest_uid))
	return
}
