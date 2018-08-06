package main

import "code.yy.com/yytars/goframework/tars/servant"

import (
	"bilin/bcserver/config"
	"bilin/bcserver/handler"
	"bilin/protocol"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"math/rand"
	"time"
)

func main() {
	appzaplog.Debug("Enter main")

	if err := config.InitAndSubConfig("appconfig.json"); err != nil {
		appzaplog.Error("InitAndSubConfig failed", zap.Error(err))
		return
	}

	//利用当前时间的UNIX时间戳初始化rand包
	rand.Seed(time.Now().UnixNano())

	// 提供给客户端，调用直播间信令
	srvObj := handler.NewBCServantObj()
	dispObj := bilin.NewBCServantDispatcher()
	if err := servant.AddServant(dispObj, srvObj, "BCServantObj"); err != nil {
		appzaplog.Error("AddPbServant failed", zap.Error(err))
		return
	}

	//// 提供给服务端业务(bilin_java)，控制直播间官频上下麦，查询直播间人数等
	//mux := servant.NewTarsHttpMux()
	//queryHttpObj := handler.NewQueryHttpObj()
	//mux.HandleFunc("/v1/internal/QueryRoom", queryHttpObj.QueryRoom)
	//mux.HandleFunc("/v1/internal/OfficialOnMike", queryHttpObj.OfficialOnMike)
	//mux.HandleFunc("/v1/internal/OfficialOffMike", queryHttpObj.OfficialOffMike)
	//if err := servant.AddHttpServant(mux, "QueryHttpObj"); err != nil {
	//	appzaplog.Error("AddHttpServant failed", zap.Error(err))
	//	return
	//}

	servant.Run()
}
