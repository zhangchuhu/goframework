package main

import "code.yy.com/yytars/goframework/tars/servant"

import (
	"bilin/rcserver/handler"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"bilin/protocol"
)

func main() {
	appzaplog.Debug("Enter main")
	srvObj := handler.NewRcServantObj()
	dispObj := bilin.NewRcServantDispatcher()
	if err := servant.AddServant(dispObj, srvObj, "RcServantObj");err != nil{
		appzaplog.Error("AddPbServant failed",zap.Error(err))
		return
	}
	servant.Run()
}
