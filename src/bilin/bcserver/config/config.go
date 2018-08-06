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
	JavaThriftAddr      []string `json:"java_thrift_addr"`
	ActTaskThriftAddr   string   `json:"act_task_thrift_addr"`
	MsgFilterThriftAddr string   `json:"msg_filter_thrift_addr"`
	RedisAddr           string   `json:"redis_addr"`
	PushProxyAddr       string   `json:"push_proxy"`
	MysqlAddr           string   `json:"mysql_addr"`
	InvisibleUids       []uint64 `json:"invisible_uids"`
	KafkaAddr           []string `json:"kafka_addr"`
	KafkaTopic          string   `json:"kafka_topic"`
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
