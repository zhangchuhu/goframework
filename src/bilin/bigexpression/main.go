package main

import (
	"bilin/protocol"
	"bilin/bigexpression/handler"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/tars"
	"time"
	"bilin/bigexpression/config"
	"bilin/common/onlinepush"
	"math/rand"
)

func main() {
	appzaplog.Debug("Enter main")
	if err := config.InitAndSubConfig("appconfig.json"); err != nil {
		appzaplog.Error("InitAndSubConfig failed", zap.Error(err))
		return
	}
	onlinepush.SetUrl(config.GetAppConfig().GetPushUrl())
	rand.Seed(time.Now().UnixNano())
	srvObj := handler.NewBigExpressionObjObj()
	dispObj := bilin.NewBigExpressionObjDispatcher()
	if err := tars.AddServant(dispObj, srvObj, "BigExpressionObjObj");err != nil{
		appzaplog.Error("AddPbServant failed",zap.Error(err))
		return
	}
	tars.Run()
}
