package config

import (
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/tars"
	"encoding/json"
	"sync"
	"bilin/protocol"
)

var (
	appconfiglk sync.Mutex
	appconfig *bilin.GetEmotionConfigRes
)

func InitAndSubConfig(filename string) error {
	if err := loadconfig(filename); err != nil {
		appzaplog.Error("loadconfig failed", zap.Error(err))
		return err
	}
	go func() {
		for info := range tars.SubConfigPush() {
			if info.Filename == filename {
				if err := loadconfig(info.Filename); err != nil {
					appzaplog.Error("InitAndSubConfig reload failed", zap.String("filename", info.Filename), zap.Error(err))
					continue
				} else {
					appzaplog.Info("InitAndSubConfig reload", zap.String("filename", info.Filename))
				}
			} else {
				appzaplog.Info("InitAndSubConfig ignore not used file change", zap.String("filename", info.Filename))
				continue
			}
		}
	}()
	return nil
}

func loadconfig(filename string) error {
	var localconfig bilin.GetEmotionConfigRes
	bin, err := tars.ReadConf(filename)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bin, &localconfig); err != nil {
		appzaplog.Error("loadconfig failed", zap.Error(err))
		return err
	}
	appconfiglk.Lock()
	appconfig = &localconfig
	appconfiglk.Unlock()
	return nil
}

func GetAppConfig() (ac *bilin.GetEmotionConfigRes) {
	appconfiglk.Lock()
	ac = appconfig
	appconfiglk.Unlock()
	return
}
