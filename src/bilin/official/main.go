package main

import (
	"bilin/official/config"
	"bilin/official/dao"
	"bilin/official/handler"
	"bilin/official/service"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/tars/servant"
	"github.com/gin-gonic/gin"
	"github.com/tomwei7/gin-jsonp"
)

func main() {
	appzaplog.Debug("Enter main")
	if err := config.InitAndSubConfig("appconfig.json"); err != nil {
		appzaplog.Error("InitAndSubConfig failed", zap.Error(err))
		return
	}

	if err := service.InitTurnOverService(config.GetAppConfig()); err != nil {
		appzaplog.Error("InitTurnOverService err", zap.Error(err))
		return
	}

	if err := dao.InitMysqlDao(); err != nil {
		appzaplog.Error("InitMysqlDao failed", zap.Error(err))
		return
	}

	mux := gin.New()
	mux.Use(gin.Recovery())
	mux.Use(jsonp.JsonP())
	//获取主播用户信息
	mux.GET("/v1/guildoam/host/:id",
		handler.MetricsMiddleWare("HostInfo",
			handler.AuthMiddleWare(handler.GetHost)),
	)

	var (
		hostguildhandler = handler.AuthMiddleWare(handler.GetHostGuild)
	)
	mux.GET("/v1/guildoam/host/:id/guild",
		handler.MetricsMiddleWare("HostGuild", hostguildhandler),
	)

	var (
		hostcontracthandler = handler.AuthMiddleWare(handler.GetHostContract)
	)
	mux.GET("/v1/guildoam/host/:id/contract",
		handler.MetricsMiddleWare("HostContract", hostcontracthandler),
	)

	// 获取指定host的收入明细
	var (
		incomingrechandler = handler.AuthMiddleWare(handler.GetHostIncomingRecords)
	)
	mux.GET("/v1/guildoam/host/:id/incomingrecord",
		handler.MetricsMiddleWare("HostIncomingRecord", incomingrechandler),
	)

	// 获取指定host的收礼明细
	var (
		incomingdetailhandler = handler.AuthMiddleWare(handler.GetHostIncomingDetail)
	)
	mux.GET("/v1/guildoam/host/:id/incomingdetail",
		handler.MetricsMiddleWare("HostIncomingDetail", incomingdetailhandler),
	)

	// 获取指定host的开播信息记录
	var (
		livingrechandler = handler.AuthMiddleWare(handler.GetHostLivingRecords)
	)
	mux.GET("/v1/guildoam/host/:id/livingrecord",
		handler.MetricsMiddleWare("HostLivingRecord", livingrechandler),
	)

	var (
		// 获取指定OW的工会信息
		guildhandler = handler.AuthMiddleWare(handler.GetGuildByOwID)
		// 修改指定ow的工会信息
		guildupdatehandler = handler.AuthMiddleWare(handler.UpdateGuildByOwID)
	)
	mux.OPTIONS("/v1/guildoam/ow/:id/guild", handler.OptionGuildByOwID).
		GET("/v1/guildoam/ow/:id/guild", handler.MetricsMiddleWare("GetOwGuild", guildhandler)).
		PUT("/v1/guildoam/ow/:id/guild", handler.MetricsMiddleWare("CreateOwGuild", guildupdatehandler))

	// 获取指定工会的本月收入明细
	var (
		guildincomingrechandler = handler.AuthMiddleWare(handler.GetGuildIncomingRecords)
	)
	mux.GET("/v1/guildoam/guild/:id/incomingrecord",
		handler.MetricsMiddleWare("GuildIncomingRecord", guildincomingrechandler),
	)

	// 获取指定工会的本月收礼明细
	var (
		guildincomingdetailhandler = handler.AuthMiddleWare(handler.GetGuildIncomingDetail)
	)
	mux.GET("/v1/guildoam/guild/:id/incomingdetail",
		handler.MetricsMiddleWare("GuildIncomingDetail", guildincomingdetailhandler),
	)

	// 获取指定工会的所有签约记录
	var (
		contractrecordshandler = handler.AuthMiddleWare(handler.GetGuildContractRecords)
	)
	mux.GET("/v1/guildoam/guild/:id/contractrecord",
		handler.MetricsMiddleWare("GuildContractRecord", contractrecordshandler),
	)

	// 获取指定工会的所有开播信息记录
	var (
		guildlivingrechandler = handler.AuthMiddleWare(handler.GetGuildLivingRecords)
	)
	mux.GET("/v1/guildoam/guild/:id/livingrecord",
		handler.MetricsMiddleWare("GuildLivingRecord", guildlivingrechandler),
	)

	var (
		guildroomhandler = handler.AuthMiddleWare(handler.GetGuildRoom)
	)
	mux.GET("/v1/guildoam/guild/:id/room",
		handler.MetricsMiddleWare("GuildRoom", guildroomhandler),
	)

	if err := servant.AddHttpServant(mux, "OfficialObj"); err != nil {
		appzaplog.Error("AddPbServant failed", zap.Error(err))
		return
	}
	// 内网绑定
	if err := servant.AddHttpServant(mux, "OfficialObjInternal"); err != nil {
		appzaplog.Error("AddPbServant failed", zap.Error(err))
		return
	}
	servant.Run()
}
