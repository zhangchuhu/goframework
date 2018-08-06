package main

import (
	"bilin/ccserver/config"
	"bilin/ccserver/handler"
	"bilin/common/onlinepush"
	"bilin/common/onlinequery"
	"bilin/protocol"

	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/tars/servant"
)

const (
	//URL = "http://bilingopush.yy.com/1" // 正式环境
	URL = "http://test-goim.yy.com:7172/1" // 测试环境
)

func configLoaded(appconfig *config.AppConfig) {
	onlinepush.URL = appconfig.OnlinePushURL
	appzaplog.Info("Set onlinepush.URL", zap.Any("url", onlinepush.URL))
	onlinequery.URL = appconfig.OnlineQueryURL
	appzaplog.Info("Set onlinequery.URL", zap.Any("url", onlinequery.URL))
}

func main() {
	appzaplog.Debug("Enter main")
	if err := config.InitAndSubConfig("appconfig.json", configLoaded); err != nil {
		appzaplog.Error("InitAndSubConfig failed", zap.Error(err))
		return
	}

	srvObj := handler.NewCCServantObj()
	dispObj := bilin.NewCCServantDispatcher()
	if err := servant.AddServant(dispObj, srvObj, "CCServantObj"); err != nil {
		appzaplog.Error("AddPbServant failed", zap.Error(err))
		return
	}

	servant.Run()
}
