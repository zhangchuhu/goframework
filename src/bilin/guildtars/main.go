package main

import (
	"bilin/guildtars/config"
	"bilin/guildtars/dao"
	"bilin/guildtars/handler"
	"bilin/guildtars/service"
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

	if err := service.InitTurnOverService(config.GetAppConfig()); err != nil {
		appzaplog.Error("InitTurnOverService err", zap.Error(err))
		return
	}

	if err := dao.InitMysqlDao(); err != nil {
		appzaplog.Error("InitMysqlDao failed", zap.Error(err))
		return
	}

	srvObj := handler.NewGuildTarsObj()
	dispObj := bilin.NewGuildTarsDispatcher()
	if err := tars.AddServant(dispObj, srvObj, "GuildTarsObj"); err != nil {
		appzaplog.Error("AddPbServant failed", zap.Error(err))
		return
	}
	tars.Run()
}
