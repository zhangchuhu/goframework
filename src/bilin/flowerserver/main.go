package main

import (
	"bilin/flowerserver/handler"
	"bilin/protocol"
	"encoding/json"

	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/tars"
)

func main() {
	appzaplog.Debug("Enter main")

	confBytes, err := tars.ReadConf("appconfig.json")
	if err != nil {
		appzaplog.Error("tars.ReadConf failed", zap.Error(err))
		return
	}
	var conf handler.AppConfig
	if err := json.Unmarshal(confBytes, &conf); err != nil {
		appzaplog.Error("json.Unmarshal failed", zap.Error(err))
		return
	}

	srvObj := handler.NewFlowerServantObj(conf)
	dispObj := bilin.NewFlowerServantDispatcher()
	if err := tars.AddServant(dispObj, srvObj, "FlowerServantObj"); err != nil {
		appzaplog.Error("AddPbServant failed", zap.Error(err))
		return
	}
	tars.Run()
}
