package main

import (
	"bilin/protocol/userinfocenter"
	"bilin/userinfocenter/config"
	"bilin/userinfocenter/dao"
	"bilin/userinfocenter/handler"
	"bilin/userinfocenter/nsqwrapper"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/tars"
)

func main() {
	appzaplog.Debug("Enter main")
	if err := config.InitAndSubConfig("appconfig.json"); err != nil {
		appzaplog.Error("InitAndSubConfig failed", zap.Error(err))
		return
	}

	if err := dao.InitMySqlDao(); err != nil {
		appzaplog.Error("InitMySqlDao failed", zap.Error(err))
		return
	}

	if err := dao.InitRedisDao(); err != nil {
		appzaplog.Error("InitMySqlDao failed", zap.Error(err))
		return
	}

	if err := dao.InitAvatarCache(); err != nil {
		appzaplog.Error("InitAvatarCache failed", zap.Error(err))
		return
	}

	if conf := config.GetAppConfig(); conf != nil && len(conf.NSQLOOKUPAddrs) > 0 {
		if err := nsqwrapper.InitNsq(conf.NSQLOOKUPAddrs); err != nil {
			appzaplog.Error("InitNsq failed", zap.Error(err))
			return
		}
	}

	if conf := config.GetAppConfig(); conf == nil || len(conf.AppleCheckThrift) == 0 {
		appzaplog.Error("get AppleCheckThrift fail")
		return
	} else {
		dao.InitThriftConnentPool(conf.AppleCheckThrift)
	}

	srvObj := handler.NewUserInfoCenterObj()
	dispObj := userinfocenter.NewUserInfoCenterObjDispatcher()
	if err := tars.AddServant(dispObj, srvObj, "UserInfoCenterObj"); err != nil {
		appzaplog.Error("AddPbServant failed", zap.Error(err))
		return
	}
	tars.Run()
}
