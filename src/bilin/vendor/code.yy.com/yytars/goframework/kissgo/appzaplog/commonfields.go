package appzaplog

import (
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap/zapcore"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
)

func UID(uid uint64) zapcore.Field {
	return zap.Uint64("uid",uid)
}

func ROOMID(roomid uint64) zapcore.Field {
	return zap.Uint64("roomid",roomid)
}

