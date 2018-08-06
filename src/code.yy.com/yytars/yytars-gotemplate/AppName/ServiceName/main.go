package main

import (
	"%{AppName}/protocol"
	"%{AppName}/%{ServiceName}/handler"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/tars"
)

func main() {
	appzaplog.Debug("Enter main")
	srvObj := handler.New%{ServantName}Obj()
	dispObj := %{AppName}.New%{ServantName}Dispatcher()
	if err := tars.AddServant(dispObj, srvObj, "%{ServantName}Obj");err != nil{
		appzaplog.Error("AddPbServant failed",zap.Error(err))
		return
	}
	tars.Run()
}
