package main

import (
	"bilin/bcserver/config"
	"bilin/operationManagement/handler"
	"bilin/protocol"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/tars/servant"
	"net/http"
)

func main() {
	appzaplog.Debug("Enter main")

	if err := config.InitAndSubConfig("appconfig.json"); err != nil {
		appzaplog.Error("InitAndSubConfig failed", zap.Error(err))
		return
	}

	mux := servant.NewTarsHttpMux()
	httpObjSrv := handler.NewOptManagementHttpObj()
	mux.HandleFunc("/hello", httpObjSrv.Hello)
	mux.HandleFunc("/api/vip/headgear", httpObjSrv.HeadgearOperation)
	mux.HandleFunc("/api/totaldata", httpObjSrv.GetAllHeadgears)
	mux.HandleFunc("/api/token", httpObjSrv.GenBS2Token)

	//主播开播流水
	mux.HandleFunc("/api/livingrecord", httpObjSrv.GetAllLivingRecord)

	//static 文件入口
	mux.Handle("/", http.FileServer(assetFS()))

	if err := servant.AddHttpServant(mux, "OptManagementHttpObj"); err != nil {
		appzaplog.Error("AddHttpServant failed", zap.Error(err))
		return
	}

	// 提供给客户端，调用直播间信令
	srvObj := handler.NewOptManagementPbObj()
	dispObj := bilin.NewOperationManagementServantDispatcher()
	if err := servant.AddServant(dispObj, srvObj, "OptManagementPbObj"); err != nil {
		appzaplog.Error("AddPbServant failed", zap.Error(err))
		return
	}

	servant.Run()
}
