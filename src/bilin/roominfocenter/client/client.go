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
	comm.SetProperty("locator", "tars.tarsregistry.QueryObj@tcp -h 183.36.111.89 -p 17890")
	objName := fmt.Sprintf("bilin.roominfocenter.RoomInfoCenterServantObj")
	client := bilin.NewRoomInfoServantClient(objName, comm)
	resp, err := client.LivingRoomsInfo(context.TODO(), &bilin.LivingRoomsInfoReq{})
	if err != nil {
		appzaplog.Error("LivingRoomsInfo err", zap.Error(err))
		return
	}
	appzaplog.Debug("resp msg", zap.Any("resp", resp))
}
