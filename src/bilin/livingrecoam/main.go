package main

import (
	"bilin/livingrecoam/config"
	"bilin/livingrecoam/handler"
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

	mux := gin.New()
	mux.Use(gin.Recovery())
	mux.Use(jsonp.JsonP())
	mux.Use(handler.CORSMiddleWare())
	//工会信息推荐
	mux.OPTIONS("/v1/livingrec/room", handler.CORSOptionHandler).
		POST("/v1/livingrec/room", handler.UpdateGuildRec).
		GET("/v1/livingrec/room", handler.GetGuildRec).
		PUT("/v1/livingrec/room", handler.CreateGuildRec)

	mux.OPTIONS("/v1/livingrec/room/:roomid", handler.CORSOptionHandler).
		DELETE("/v1/livingrec/room/:id", handler.DelGuildRec)

	//主播信息推荐
	mux.OPTIONS("/v1/livingrec/host", handler.CORSOptionHandler).
		GET("/v1/livingrec/host", handler.GetHostRec).
		PUT("/v1/livingrec/host", handler.CreateHostRec).
		POST("/v1/livingrec/host", handler.UpdateHostRec)
	mux.OPTIONS("/v1/livingrec/host/:id", handler.CORSOptionHandler).
		DELETE("/v1/livingrec/host/:id", handler.DelHostRec)

	// 置顶区推荐
	mux.OPTIONS("/v1/livingrec/stickypost", handler.CORSOptionHandler).
		GET("/v1/livingrec/stickypost", handler.GetStickyPost).
		PUT("/v1/livingrec/stickypost", handler.AddStickyPost).
		POST("/v1/livingrec/stickypost", handler.UpdateStickyPost)

	mux.OPTIONS("/v1/livingrec/stickypost/:id", handler.CORSOptionHandler).
		DELETE("/v1/livingrec/stickypost/:id", handler.DelStickyPost)

	// 工会房间关系管理
	mux.OPTIONS("/v1/guildoam/room", handler.CORSOptionHandler).
		GET("/v1/guildoam/room", handler.GetGuildRoom).
		PUT("/v1/guildoam/room", handler.CreateGuildRoom).
		POST("/v1/guildoam/room", handler.UpdateGuildRoom)

	mux.OPTIONS("/v1/guildoam/room/:id", handler.CORSOptionHandler).
		DELETE("/v1/guildoam/room/:id", handler.DelGuildRoom)

	// 工会房间关系管理
	mux.OPTIONS("/v1/guildoam/guild", handler.CORSOptionHandler).
		GET("/v1/guildoam/guild", handler.GetGuild).
		PUT("/v1/guildoam/guild", handler.CreateGuild).
		POST("/v1/guildoam/guild", handler.UpdateGuild)

	mux.OPTIONS("/v1/guildoam/guild/:id", handler.CORSOptionHandler).
		DELETE("/v1/guildoam/guild/:id", handler.DelGuild)

	// 工会签约关系管理
	mux.OPTIONS("/v1/guildoam/contract", handler.CORSOptionHandler).
		GET("/v1/guildoam/contract", handler.GetContract).
		PUT("/v1/guildoam/contract", handler.CreateContract).
		POST("/v1/guildoam/contract", handler.UpdateContract)

	mux.OPTIONS("/v1/guildoam/contract/:id", handler.CORSOptionHandler).
		DELETE("/v1/guildoam/contract/:id", handler.DelContract)

	// 聊妹套话
	mux.OPTIONS("/v1/rc/puatopic", handler.CORSOptionHandler).
		GET("/v1/rc/puatopic", handler.GetPUATopic).
		PUT("/v1/rc/puatopic", handler.CreatePUATopic).
		POST("/v1/rc/puatopic", handler.UpdatePUATopic)
	mux.OPTIONS("/v1/rc/puatopic/:id", handler.CORSOptionHandler).
		DELETE("/v1/rc/puatopic/:id", handler.DelPUATopic)

	// 真心话
	mux.OPTIONS("/v1/rc/truthtopic", handler.CORSOptionHandler).
		GET("/v1/rc/truthtopic", handler.GetTruthTopic).
		PUT("/v1/rc/truthtopic", handler.CreateTruthTopic).
		POST("/v1/rc/truthtopic", handler.UpdateTruthTopic)
	mux.OPTIONS("/v1/rc/truthtopic/:id", handler.CORSOptionHandler).
		DELETE("/v1/rc/truthtopic/:id", handler.DelTruthTopic)

	// 聊天标签
	mux.OPTIONS("/v1/rc/chattag", handler.CORSOptionHandler).
		GET("/v1/rc/chattag", handler.GetChatTag).
		PUT("/v1/rc/chattag", handler.CreateChatTag).
		POST("/v1/rc/chattag", handler.UpdateChatTag)
	mux.OPTIONS("/v1/rc/chattag/:id", handler.CORSOptionHandler).
		DELETE("/v1/rc/chattag/:id", handler.DelChatTag)

	//登录
	mux.OPTIONS("/v1/user/login", handler.CORSOptionHandler).
		POST("/v1/user/login", handler.Login)

	// 登出
	mux.OPTIONS("/v1/user/logout", handler.CORSOptionHandler).
		POST("/v1/user/logout", handler.LoginOut)

	mux.OPTIONS("/v1/user/info", handler.CORSOptionHandler).
		GET("/v1/user/info", handler.UserInfo)

	if err := servant.AddHttpServant(mux, "LivingRecOAMObj"); err != nil {
		appzaplog.Error("AddPbServant failed", zap.Error(err))
		return
	}
	// 内网绑定
	//if err := servant.AddHttpServant(mux, "OfficialObjInternal"); err != nil {
	//	appzaplog.Error("AddPbServant failed", zap.Error(err))
	//	return
	//}
	servant.Run()
}
