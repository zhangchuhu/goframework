package config

import (
	"encoding/json"
	"sync"

	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/tars"
)

const (
	SearchURLProdEnv = "http://searchservice.yy.com/search" // 正式环境
	SearchURLTestEnv = "http://14.17.109.28:8088/search"    // 测试环境

	UpdateURL = "http://updateservice.yy.com"
)

var (
	appconfiglk sync.Mutex
	appconfig   *AppConfig
)

type AppConfig struct {
	SearchURL    string
	UpdateURL    string
	KafkaBrokers []string
	KafkaTopics  []string
}

func init() {
	appconfig = &AppConfig{
		SearchURL: SearchURLTestEnv,
		UpdateURL: UpdateURL,
		KafkaBrokers: []string{
			"kafkasz001-test001.yy.com:8101",
			"kafkasz001-test002.yy.com:8101",
			"kafkasz001-test003.yy.com:8101",
		},
		KafkaTopics: []string{
			"bilin_user_update_test",
			"bilin_room_update_test",
			"bilin_song_update_test",
		},
	}
}

func InitAndSubConfig(filename string, cb func(appconfig *AppConfig)) error {
	const FuncName = "InitAndSubConfig: "
	//read appconfig first
	if err := loadconfig(filename); err != nil {
		log.Error(FuncName+"load config fail", zap.Error(err))
		return err
	}
	cb(appconfig)
	go func() {
		for info := range tars.SubConfigPush() {
			if info.Filename == filename {
				if err := loadconfig(info.Filename); err != nil {
					log.Error(FuncName+"reload config fail", zap.String("filename", info.Filename), zap.Error(err))
					continue
				} else {
					log.Info(FuncName+"reload ok", zap.String("filename", info.Filename))
					cb(appconfig)
				}
			} else {
				log.Info(FuncName+"ignore not used file change", zap.String("filename", info.Filename))
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
