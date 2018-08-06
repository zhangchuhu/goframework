package main

import (
	"bilin/chattagtars/cache"
	"bilin/chattagtars/config"
	"bilin/chattagtars/dao"
	"bilin/chattagtars/handler"
	"bilin/protocol"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/tars"
)

func main() {
	appzaplog.Debug("Enter main")
	if err := config.InitAndSubConfig("appconfig.json"); err != nil {
		appzaplog.Error("InitAndSubConfig failed", zap.Error(err))
		return
	}

	if err := dao.InitMysqlDao(); err != nil {
		appzaplog.Error("InitMysqlDao failed", zap.Error(err))
		return
	}

	if err := cache.InitCache(); err != nil {
		appzaplog.Error("InitCache failed", zap.Error(err))
		return
	}

	srvObj := handler.NewChatTagTarsObj()
	dispObj := bilin.NewChatTagTarsDispatcher()
	if err := tars.AddServant(dispObj, srvObj, "ChatTagTarsObj"); err != nil {
		appzaplog.Error("AddPbServant failed", zap.Error(err))
		return
	}
	tars.Run()
}
