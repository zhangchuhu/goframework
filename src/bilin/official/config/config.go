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
	RateLimiteRate     int64  `json:"rate_limite_rate"`
	HuJiaoCallDb       string `json:"hu_jiao_call_db"`
	BIDataDb           string `json:"bi_data_db"`
	ContractThriftAddr string `json:"contract_thrift_addr"`
	TurnOverThriftAddr string `json:"turn_over_thrift_addr"`
	CORSACAO           string `json:"corsacao"` // 跨域请求的源路径
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
