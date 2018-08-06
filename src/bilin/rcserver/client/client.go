package main

import (
	"context"
	"code.yy.com/yytars/goframework/tars/servant"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"fmt"
	"bilin/protocol"
)

func main() {
	comm := servant.NewPbCommunicator()
	objName := fmt.Sprintf("bilin.rcserver.RcServantObj")
	client := bilin.NewRcServantClient(objName,comm)
	resp, err := client.StartRandomCall(context.TODO(), &bilin.StartRandomCallReq{ClickCount: 10})
	if err != nil {
		appzaplog.Error("SayGreeting err",zap.Error(err))
		return
	}
	appzaplog.Debug("resp msg",zap.Any("resp",resp))
}
