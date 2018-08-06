package config

import (
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/tars"
	"encoding/json"
)

var (
	appconfig *AppConfig
)

type AppConfig struct {
	RabbitMqAddr []string `json:"rabbitmq_addr"`
	MysqlAddr    string   `json:"mysql_addr"`
	RedisAddr    string   `json:"redis_addr"`
	SentinelAddr []string `json:"sentinel_addr"`
}

func InitAndSubConfig(filename string) error {
	//read appconfig first
	if err := loadconfig(filename); err != nil {
		appzaplog.Error("loadconfig failed", zap.Error(err))
		return err
	}
	return nil
}

func loadconfig(filename string) error {
	var localconfig AppConfig
	bin, err := tars.ReadConf(filename)
	if err != nil {
		appzaplog.Error("ReadConf failed", zap.Error(err))
		return err
	}
	if err := json.Unmarshal(bin, &localconfig); err != nil {
		appzaplog.Error("loadconfig failed", zap.Error(err))
		return err
	}
	appconfig = &localconfig
	return nil
}

func GetAppConfig() *AppConfig {
	return appconfig
}

func SetTestAppConfig(test *AppConfig) {
	appconfig = test
}
