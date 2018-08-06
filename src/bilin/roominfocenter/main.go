/*
从bcserver的redis定时同步直播间的基础数据。
当bcserver的数据结构发生改变时，需要做相应调整。
可以共用bcserver的代码
*/
package main

import (
	"bilin/common/cacheprocessor"
	"bilin/protocol"
	"bilin/roominfocenter/cache"
	"bilin/roominfocenter/config"
	"bilin/roominfocenter/dao"
	"bilin/roominfocenter/handler"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/tars"
	"time"
)

func main() {
	appzaplog.Debug("Enter main")

	if err := config.InitAndSubConfig("appconfig.json"); err != nil {
		appzaplog.Error("InitAndSubConfig failed", zap.Error(err))
		return
	}

	if err := dao.InitRedisDao(); err != nil {
		appzaplog.Error("InitRedisDao failed", zap.Error(err))
		return
	}

	if err := cacheprocessor.CacheProcessor("RefreshRoomCache", 3*time.Second, cache.RefreshRoomCache); err != nil {
		appzaplog.Error("CacheProcessor failed", zap.Error(err))
		return
	}

	srvObj := handler.NewRoomInfoCenterServantObj()
	dispObj := bilin.NewRoomInfoServantDispatcher()
	if err := tars.AddServant(dispObj, srvObj, "RoomInfoCenterServantObj"); err != nil {
		appzaplog.Error("AddPbServant failed", zap.Error(err))
		return
	}
	tars.Run()
}
