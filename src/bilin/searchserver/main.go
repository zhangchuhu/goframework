package main

import (
	"bilin/protocol"
	"bilin/searchserver/config"
	"bilin/searchserver/handler"
	"bilin/searchserver/updater"

	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/tars/servant"
)

func configLoaded(appconfig *config.AppConfig) {
}

func main() {
	log.Info("enter main")
	if err := config.InitAndSubConfig("appconfig.json", configLoaded); err != nil {
		log.Error("call InitAndSubConfig fail", zap.Error(err))
		return
	}

	srvObj := handler.NewSearchServantObj()
	dispObj := bilin.NewSearchServantDispatcher()
	if err := servant.AddServant(dispObj, srvObj, "SearchServantObj"); err != nil {
		log.Error("call AddServant fail", zap.Error(err))
		return
	}

	go updater.KafkaLoop()
	servant.Run()
}
