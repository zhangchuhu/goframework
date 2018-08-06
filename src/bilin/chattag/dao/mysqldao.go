package dao

import (
	"bilin/chattag/config"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"strings"
)

var (
	hujiaoChatTagDB *gorm.DB
	dbNotInitErr    = errors.New("db not init")
)

func InitMysqlDao() error {
	var err error
	if conf := config.GetAppConfig(); conf != nil {
		hujiaoChatTagDB, err = gorm.Open("mysql", conf.ChatTagDb)
		if err != nil {
			appzaplog.Error("gorm open ChatTagDb failed", zap.Error(err))
			return err
		}
	} else {
		return errors.New("no appconfig find")
	}

	return nil
}

func IsTableNotExistErr(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "Table") && strings.Contains(errStr, "doesn't exist")
}
