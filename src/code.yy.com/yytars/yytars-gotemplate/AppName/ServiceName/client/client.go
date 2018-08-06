package main

import (
	"context"
	"%{AppName}/protocol"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"fmt"
	"code.yy.com/yytars/goframework/tars"
)

func main() {
	comm := tars.NewCommunicator()
	objName := fmt.Sprintf("%{AppName}.%{ServiceName}.%{ServantName}Obj")
	client := %{AppName}.New%{ServantName}Client(objName,comm)
	resp, err := client.SayGreeting(context.TODO(), &%{AppName}.Request{Name: "Leo"})
	if err != nil {
		appzaplog.Error("SayGreeting err",zap.Error(err))
		return
	}
	appzaplog.Debug("resp msg",zap.Any("resp",resp))
}
