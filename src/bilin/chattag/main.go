package main

import (
	"bilin/chattag/cache"
	"bilin/chattag/config"
	"bilin/chattag/dao"
	"bilin/chattag/handler"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/tars/servant"
	"github.com/gin-gonic/gin"
)

func main() {
	appzaplog.Debug("Enter main")
	if err := config.InitAndSubConfig("appconfig.json"); err != nil {
		appzaplog.Error("InitAndSubConfig failed", zap.Error(err))
		return
	}

	if err := dao.InitMysqlDao(); err != nil {
		appzaplog.Error("InitMysqlDao failed", zap.Error(err))
		return
	}

	if err := dao.InitRedisDao(config.GetAppConfig()); err != nil {
		appzaplog.Error("InitRedisDao failed", zap.Error(err))
		return
	}
	defer dao.RedisClient.Close()

	if err := cache.InitCache(); err != nil {
		appzaplog.Error("InitCache failed", zap.Error(err))
		return
	}

	mux := gin.New()
	mux.Use(gin.Recovery())

	mux.POST("/v1/rc/chattag/taglist", handler.MetricsMiddleWare("taglist",
		handler.AuthMiddleWare(handler.GetChatTagsList)))
	mux.POST("/v1/rc/chattag/settag", handler.MetricsMiddleWare("settag",
		handler.AuthMiddleWare(handler.SetChatTag)))
	mux.POST("/v1/rc/chattag/tagstatus", handler.MetricsMiddleWare("tagstatus",
		handler.AuthMiddleWare(handler.HandleTagStatus)))

	//玩法
	// 玩法列表
	mux.POST("/v1/rc/plugins", handler.MetricsMiddleWare("plugins",
		handler.AuthMiddleWare(handler.HandleRCPlugin)))
	mux.POST("/v1/rc/puatopic", handler.MetricsMiddleWare("puatopic",
		handler.AuthMiddleWare(handler.HandlePUATopic)))
	mux.POST("/v1/rc/truthtopic", handler.MetricsMiddleWare("truthtopic",
		handler.AuthMiddleWare(handler.HandleTruthGame)))
	mux.POST("/v1/rc/userdetailincall", handler.MetricsMiddleWare("userdetailincall",
		handler.AuthMiddleWare(handler.HandleUserDetailInCall)))

	if err := servant.AddHttpServant(mux, "ChatTagObj"); err != nil {
		appzaplog.Error("AddHttpServant ChatTagObj failed", zap.Error(err))
		return
	}

	if err := servant.AddHttpServant(mux, "ChatTagObjInternal"); err != nil {
		appzaplog.Error("AddHttpServant ChatTagObjInternal failed", zap.Error(err))
		return
	}

	servant.Run()
}
