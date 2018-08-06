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
	objName := fmt.Sprintf("bilin.chattagtars.ChatTagTarsObj")
	client := bilin.NewChatTagTarsClient(objName, comm)
	resp, err := client.RChatTag(context.TODO(), &bilin.RChatTagReq{})
	if err != nil {
		appzaplog.Error("SayGreeting err", zap.Error(err))
		return
	}
	appzaplog.Debug("resp msg", zap.Any("resp", resp))
}
