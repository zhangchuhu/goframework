package cache

import (
	"bilin/clientcenter"
	"bilin/common/cacheprocessor"
	"bilin/protocol"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"context"
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
	info, err := clientcenter.ChatTagClient().RChatTag(context.TODO(), &bilin.RChatTagReq{}) //dao.GetAll()
	if err != nil {
		appzaplog.Error("refreshChatTag dao.GetAll err", zap.Error(err))
		return err
	}
	tagmap := make(map[int64]*bilin.ChatTag)
	for _, v := range info.Chattag {
		tagmap[int64(v.Id)] = v
	}
	gCache.Store("chattag", tagmap)
	return nil
}

func TakeChatTagCache() map[int64]*bilin.ChatTag {
	info, ok := gCache.Load("chattag")
	if ok {
		return info.(map[int64]*bilin.ChatTag)
	}
	return nil
}
