package main

import (
	"bilin/adpromotion/config"
	"bilin/adpromotion/handler"
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

	// 提供给奇虎360的点击事件上报接口
	mux := servant.NewTarsHttpMux()
	httpObjSrv := handler.NewAdPromotionHttpObj()
	mux.HandleFunc("/hello", httpObjSrv.Hello)
	mux.HandleFunc("/v1/qihu360/click", httpObjSrv.ClickAd)
	if err := servant.AddHttpServant(mux, "ADPromotionHttpObj"); err != nil {
		appzaplog.Error("AddHttpServant failed", zap.Error(err))
		return
	}

	//内网绑定
	if err := servant.AddHttpServant(mux, "ADPromotionHttpObjInternal"); err != nil {
		appzaplog.Error("ADPromotionHttpObjInternal failed", zap.Error(err))
		return
	}

	servant.Run()
}
