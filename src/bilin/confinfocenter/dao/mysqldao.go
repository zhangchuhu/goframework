package dao

import (
	"bilin/confinfocenter/config"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"errors"
	"github.com/jinzhu/gorm"
)

var (
	hujiaoUserDB *gorm.DB
	hujiaoCallDB *gorm.DB
	hujiaoDB     *gorm.DB
	hujiaoRecDB  *gorm.DB
	dbNotInitErr = errors.New("db not init")
)

func InitMysqlDao() error {
	var err error
	if conf := config.GetAppConfig(); conf != nil {
		hujiaoUserDB, err = gorm.Open("mysql", conf.HuJiaoUserDb)
		if err != nil {
			appzaplog.Error("gorm open HuJiaoUserDb failed", zap.Error(err))
			return err
		}

		hujiaoCallDB, err = gorm.Open("mysql", conf.HuJiaoCallDb)
		if err != nil {
			appzaplog.Error("gorm open HuJiaoCallDb failed", zap.Error(err))
			return err
		}

		hujiaoDB, err = gorm.Open("mysql", conf.HuJiaoDb)
		if err != nil {
			appzaplog.Error("gorm open HuJiaoDb failed", zap.Error(err))
			return err
		}

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
