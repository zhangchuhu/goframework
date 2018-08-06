package config

import (
	"encoding/json"
	"sync"

	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/tars"
)

var (
	appconfiglk sync.Mutex
	appconfig   *AppConfig
)

type AppConfig struct {
	OnlinePushURL  string
	OnlineQueryURL string
}

func InitAndSubConfig(filename string, cb func(appconfig *AppConfig)) error {
	//read appconfig first
	if err := loadconfig(filename); err != nil {
		appzaplog.Error("loadconfig failed", zap.Error(err))
		return err
	}
	cb(appconfig)
	go func() {
		for info := range tars.SubConfigPush() {
			if info.Filename == filename {
				if err := loadconfig(info.Filename); err != nil {
					appzaplog.Error("InitAndSubConfig reload failed", zap.String("filename", info.Filename), zap.Error(err))
					continue
				} else {
					appzaplog.Info("InitAndSubConfig reload", zap.String("filename", info.Filename))
					cb(appconfig)
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
	var localconfig AppConfig
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

func GetAppConfig() *AppConfig {
	return appconfig
}
