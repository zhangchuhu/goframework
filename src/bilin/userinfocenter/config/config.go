package config

import (
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/tars"
	"encoding/json"
	"sync"
)

var (
	appconfiglk sync.Mutex
	appconfig   *AppConfig
)

type AppConfig struct {
	UserInfoRedis string `json:"user_info_redis"`
	HuJiaoDb      string `json:"hu_jiao_db"`
	AvatarDb      string `json:"avatar_db"`
	AttentionDb   string `json:"attention_db"`
	//NSQDAddrs        []string `json:"nsqd_addrs"`
	NSQLOOKUPAddrs   []string `json:"nsqlookup_addrs"`
	AppleCheckThrift string   `json:"apple_check_thrift"`
}

func InitAndSubConfig(filename string) error {
	//read appconfig first
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

func GetAppConfig() (ac *AppConfig) {
	appconfiglk.Lock()
	ac = appconfig
	appconfiglk.Unlock()
	return
}
