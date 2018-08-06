package main

import (
	"context"
	"bilin/protocol"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"fmt"
	"code.yy.com/yytars/goframework/tars"
)

func main() {
	comm := tars.NewCommunicator()
	objName := fmt.Sprintf("bilin.flowerserver.FlowerServantObj")
	client := bilin.NewFlowerServantClient(objName,comm)
	resp, err := client.SayGreeting(context.TODO(), &bilin.Request{Name: "Leo"})
	if err != nil {
		appzaplog.Error("SayGreeting err",zap.Error(err))
		return
	}
	appzaplog.Debug("resp msg",zap.Any("resp",resp))
}
