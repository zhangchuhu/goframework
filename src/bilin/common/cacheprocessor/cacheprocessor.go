package cacheprocessor

import (
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/kissgo/httpmetrics"
	"time"
)

func CacheProcessor(key string, duration time.Duration, refreshfun func() error) error {
	if err := refreshfun(); err != nil {
		return err
	}
	go func() {
		ticker := time.NewTicker(duration)
		defer ticker.Stop()
		for _ = range ticker.C {
			if err := refreshfun(); err != nil {
				httpmetrics.CounterMetric(key+"Fail", 1)
				appzaplog.Error("cacheProcessor refresh failed", zap.Error(err))
			}
		}
	}()
	return nil
}
