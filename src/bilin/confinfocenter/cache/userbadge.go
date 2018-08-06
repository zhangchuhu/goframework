package cache

import (
	"bilin/common/cacheprocessor"
	"bilin/confinfocenter/dao"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"errors"
	"sync/atomic"
	"time"
)

func InitCache() error {
	if err := cacheprocessor.CacheProcessor("userBadgeCache", 2*time.Minute, userBadgeCache); err != nil {
		appzaplog.Error("userBadgeCache start failed", zap.Error(err))
		return err
	}
	return nil
}

var userBadeg atomic.Value

func userBadgeCache() error {
	badge, err := dao.GetUserBadges()
	if err != nil {
		appzaplog.Error("GetUserBadges failed", zap.Error(err))
		return err
	}
	userBadeg.Store(badge)
	return nil
}

func GetUserBadge() ([]dao.UserBadge, error) {
	if badge := userBadeg.Load(); badge != nil {
		if ret, ok := badge.([]dao.UserBadge); ok {
			return ret, nil
		}
	}
	return nil, errors.New("not found")
}
