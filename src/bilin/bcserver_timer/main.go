package main

import "code.yy.com/yytars/goframework/tars/servant"

import (
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"bilin/bcserver/config"
)

func main() {
	appzaplog.Debug("Enter main")

	if err := config.InitAndSubConfig("appconfig.json"); err != nil {
		appzaplog.Error("InitAndSubConfig failed", zap.Error(err))
		return
	}

	NewBCServantTimerObj()
	go servant.Run()

	//wait forever
	select {}
}
