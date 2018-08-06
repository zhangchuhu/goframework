package dao

import (
	"bilin/guildtars/config"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"strings"
)

var (
	hujiaoRecDB  *gorm.DB
	dbNotInitErr = errors.New("db not init")
)

func InitMysqlDao() error {
	var err error
	if conf := config.GetAppConfig(); conf != nil {
		hujiaoRecDB, err = gorm.Open("mysql", conf.HuJiaoRecDb)
		if err != nil {
			appzaplog.Error("gorm open HuJiaoRecDb failed", zap.Error(err))
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
