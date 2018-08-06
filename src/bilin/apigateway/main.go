package main

import (
	"bilin/apigateway/cache"
	"bilin/apigateway/config"
	"bilin/apigateway/handler"
	"bilin/apigateway/services/rank"
	"bilin/apigateway/services/sensitive"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/tars/servant"
	"github.com/juju/ratelimit"
)

func main() {
	if err := config.InitAndSubConfig("appconfig.json"); err != nil {
		appzaplog.Error("InitAndSubConfig failed", zap.Error(err))
		return
	}
	rate := int64(10000)
	if conf := config.GetAppConfig(); conf != nil {
		rate = conf.RateLimiteRate
	}

	if err := rank.InitRank(); err != nil {
		appzaplog.Error("InitRank failed", zap.Error(err))
		return
	}

	// 依赖rank，注意在rank后初始化
	if err := cache.InitCache(); err != nil {
		appzaplog.Error("InitCache failed", zap.Error(err))
		return
	}

	// 载入敏感词
	if err := sensitive.LoadSensiTive(); err != nil {
		appzaplog.Error("LoadSensiTive failed", zap.Error(err))
		return
	}

	mux := servant.NewTarsHttpMux()

	//随机呼叫头像
	var (
		rcheaderimghandler = handler.AuthMiddleWareV2(handler.HandleRCHeaderImg)
		rcheaderimglimiter = ratelimit.NewBucketWithRate(float64(rate), rate)
	)
	rcheaderimghandler = handler.RateLimiteMiddleWareV2(rcheaderimglimiter, rcheaderimghandler)
	mux.HandleFunc("/v1/home/rcHeaderImg", handler.MetricsMiddleWare("rcHeaderImg", rcheaderimghandler))

	//直播间轮播页面
	var (
		carouselhandler          = handler.AuthMiddleWareV2(handler.HandleQueryCarousel)
		querycarouselratelimiter = ratelimit.NewBucketWithRate(float64(rate), rate)
	)
	carouselhandler = handler.RateLimiteMiddleWareV2(querycarouselratelimiter, carouselhandler)
	mux.HandleFunc("/v1/live/carousel", handler.MetricsMiddleWare("carousel", carouselhandler))

	// 直播间榜单入口
	var (
		propranklisthandler = handler.AuthMiddleWareV2(handler.HandleTodayRankList)
		propranklistlimiter = ratelimit.NewBucketWithRate(float64(rate), rate)
	)
	propranklisthandler = handler.RateLimiteMiddleWareV2(propranklistlimiter, propranklisthandler)
	mux.HandleFunc("/v1/live/rankList", handler.MetricsMiddleWare("rankList", propranklisthandler))

	// 直播页品类分类列表
	var (
		categorylisthandler = handler.AuthMiddleWareV2(handler.HandleCategoryList)
		categorylistlimiter = ratelimit.NewBucketWithRate(float64(rate), rate)
	)
	categorylisthandler = handler.RateLimiteMiddleWareV2(categorylistlimiter, categorylisthandler)
	mux.HandleFunc("/v1/live/categoryList", handler.MetricsMiddleWare("categoryList", categorylisthandler))

	// 开播页品类分类列表
	var (
		startlivecategorylisthandler = handler.AuthMiddleWareV2(handler.HandleStartLiveCategoryList)
		startlivecategorylistlimiter = ratelimit.NewBucketWithRate(float64(rate), rate)
	)
	startlivecategorylisthandler = handler.RateLimiteMiddleWareV2(startlivecategorylistlimiter, startlivecategorylisthandler)
	mux.HandleFunc("/v1/live/startLiveCategoryList", handler.MetricsMiddleWare("startLiveCategoryList", startlivecategorylisthandler))

	// 推荐的直播间列表
	var (
		gethomehotlinelistv2handler = handler.AuthMiddleWareV2(handler.HandleGetLivingHotLineList)
		gethomehotlinelistv2limiter = ratelimit.NewBucketWithRate(float64(rate), rate)
	)
	gethomehotlinelistv2handler = handler.RateLimiteMiddleWareV2(gethomehotlinelistv2limiter, gethomehotlinelistv2handler)
	mux.HandleFunc("/v1/live/recHotLineList", handler.MetricsMiddleWare("recHotLineList", gethomehotlinelistv2handler))

	// 推荐的直播间列表V2
	var (
		hotlinelistv2handler        = handler.AuthMiddleWareV2(handler.HandleGetLivingHotLineListV2)
		hotlinelistv2handlerlimiter = ratelimit.NewBucketWithRate(float64(rate), rate)
	)
	hotlinelistv2handler = handler.RateLimiteMiddleWareV2(hotlinelistv2handlerlimiter, hotlinelistv2handler)
	mux.HandleFunc("/v1/live/recHotLineListV2", handler.MetricsMiddleWare("recHotLineListV2", hotlinelistv2handler))

	servant.AddHttpServant(mux, "ApiGateWayObj")
	//内网绑定
	servant.AddHttpServant(mux, "ApiGateWayObjInternal")
	servant.Run()
}
