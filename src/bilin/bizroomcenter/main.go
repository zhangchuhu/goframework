package main

import (
	"bilin/bizroomcenter/config"
	"bilin/bizroomcenter/handler"
	"bilin/protocol"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/tars/servant"
)

func main() {
	appzaplog.Debug("Enter main")

	if err := config.InitAndSubConfig("appconfig.json"); err != nil {
		appzaplog.Error("InitAndSubConfig failed", zap.Error(err))
		return
	}

	srvObj := handler.NewBizRoomCenterPbObj()
	dispObj := bilin.NewBizRoomCenterServantDispatcher()
	if err := servant.AddServant(dispObj, srvObj, "BizRoomCenterPbObj"); err != nil {
		appzaplog.Error("AddPbServant failed", zap.Error(err))
		return
	}

	servant.Run()
}
