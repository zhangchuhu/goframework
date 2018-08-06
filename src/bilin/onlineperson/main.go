package main

import (
	"bilin/protocol"
	"bilin/onlineperson/handler"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/tars"
)

func main() {
	appzaplog.Debug("Enter main")
	srvObj := handler.NewOnlinePersonObjObj()
	dispObj := bilin.NewOnlinePersonObjDispatcher()
	if err := tars.AddServant(dispObj, srvObj, "OnlinePersonObjObj");err != nil{
		appzaplog.Error("AddPbServant failed",zap.Error(err))
		return
	}
	tars.Run()
}
