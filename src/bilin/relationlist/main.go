package main

import (
	"bilin/protocol"
	"bilin/relationlist/config"
	"bilin/relationlist/handler"
	"bilin/relationlist/service"
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

	service.MysqlInit()
	service.RedisInit()

	go service.ConsumerRabbitMQ()

	mux := servant.NewTarsHttpMux()
	httpObjSrv := handler.NewRelationListHttpObj()
	mux.HandleFunc("/", httpObjSrv.GetRelationListByJsonP)
	mux.HandleFunc("/test", httpObjSrv.GetRelationList)
	if err := servant.AddHttpServant(mux, "RelationListHttpObj"); err != nil {
		appzaplog.Error("AddHttpServant failed", zap.Error(err))
		return
	}

	officialMux := servant.NewTarsHttpMux()
	officialHttpObjSrv := handler.NewOfficialRelationListHttpObj()
	officialMux.HandleFunc("/", officialHttpObjSrv.OfficialRoomChangeOwner)
	if err := servant.AddHttpServant(officialMux, "OfficialRelationListHttpObj"); err != nil {
		appzaplog.Error("AddHttpServant failed", zap.Error(err))
		return
	}

	//内网绑定
	if err := servant.AddHttpServant(mux, "RelationListHttpObjInternal"); err != nil {
		appzaplog.Error("RelationListHttpObjInternal failed", zap.Error(err))
		return
	}

	srvObj := handler.NewRelationListPbObj()
	dispObj := bilin.NewRelationListServantDispatcher()
	if err := servant.AddServant(dispObj, srvObj, "RelationListPbObj"); err != nil {
		appzaplog.Error("AddPbServant failed", zap.Error(err))
		return
	}

	servant.Run()
}
