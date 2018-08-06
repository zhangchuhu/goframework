package main

import (
	"bilin/protocol"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/tars"
	"context"
	"fmt"
)

func main() {
	comm := tars.NewCommunicator()
	comm.SetProperty("locator", "tars.tarsregistry.QueryObj@tcp -h 58.215.138.213 -p 17890")
	objName := fmt.Sprintf("bilin.confinfocenter.ConfInfoServantObj")
	client := bilin.NewConfInfoServantClient(objName, comm)
	resp, err := client.LivingCategorys(context.TODO(), &bilin.LivingCategorysReq{})
	if err != nil {
		appzaplog.Error("SayGreeting err", zap.Error(err))
		return
	}
	appzaplog.Debug("resp msg", zap.Any("resp", resp))
}
