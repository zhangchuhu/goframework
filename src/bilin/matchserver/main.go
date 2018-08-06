package main

import (
	"bilin/matchserver/handler"
	"bilin/protocol"
	"encoding/json"

	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/tars"
)

func main() {
	log.Debug("enter main")

	configFileName := "appconfig.json"
	confBytes, err := tars.ReadConf(configFileName)
	if err != nil {
		log.Fatal("tars.ReadConf failed", zap.String("configFileName", configFileName), zap.Error(err))
		return
	}
	if err := json.Unmarshal(confBytes, &handler.Conf); err != nil {
		log.Fatal("json.Unmarshal failed", zap.Error(err))
		return
	}
	handler.Init()

	srvObj := handler.NewMatchServantObj()
	dispObj := bilin.NewMatchServantDispatcher()
	if err := tars.AddServant(dispObj, srvObj, "MatchServantObj"); err != nil {
		log.Fatal("tars.AddServant failed", zap.Error(err))
		return
	}

	go handler.HandleMatchTimer()
	go handler.HandleTalkingHeartTimer()
	go handler.HandleOnlineStatTimer()

	tars.Run()
}
