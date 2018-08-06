package dao

import (
	"bilin/official/config"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	hujiaoCallDB *gorm.DB
	BIDataDB     *gorm.DB

	dbNotInitErr = errors.New("db not init")
)

func InitMysqlDao() error {
	var err error
	if conf := config.GetAppConfig(); conf != nil {

		hujiaoCallDB, err = gorm.Open("mysql", conf.HuJiaoCallDb)
		if err != nil {
			appzaplog.Error("gorm open HuJiaoCallDb failed", zap.Error(err))
			return err
		}

		BIDataDB, err = gorm.Open("mysql", conf.BIDataDb)
		if err != nil {
			appzaplog.Error("gorm open BIDataDB failed", zap.Error(err))
			return err
		}
	} else {
		return errors.New("no appconfig find")
	}

	return nil
}
