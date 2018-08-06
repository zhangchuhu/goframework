package cache

import (
	"bilin/chattagtars/dao"
	"bilin/common/cacheprocessor"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"sync"
	"time"
)

var gCache sync.Map

func InitCache() error {

	//开启品类分类缓存
	if err := cacheprocessor.CacheProcessor("refreshChatTag", 30*time.Second, refreshChatTag); err != nil {
		return err
	}

	return nil
}

func refreshChatTag() error {
	info, err := dao.GetAll()
	if err != nil {
		appzaplog.Error("refreshChatTag dao.GetAll err", zap.Error(err))
		return err
	}

	tagmap := make(map[int64]dao.ChatTag)
	for _, v := range info {
		tagmap[int64(v.ID)] = v
	}
	gCache.Store("chattag", tagmap)
	return nil
}

func TakeChatTagCache() map[int64]dao.ChatTag {
	info, ok := gCache.Load("chattag")
	if ok {
		return info.(map[int64]dao.ChatTag)
	}
	return nil
}
