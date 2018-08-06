package dao

import (
	"bilin/userinfocenter/config"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"errors"
	"github.com/jinzhu/gorm"
)

var (
	UserDB       *gorm.DB
	UserAvatarDB *gorm.DB
	AttentionDB  *gorm.DB
)

func InitMySqlDao() error {
	var err error
	if conf := config.GetAppConfig(); conf != nil {
		UserDB, err = gorm.Open("mysql", conf.HuJiaoDb)
		if err != nil {
			appzaplog.Error("gorm open HuJiaoDb failed", zap.Error(err))
			return err
		}

		UserAvatarDB, err = gorm.Open("mysql", conf.AvatarDb)
		if err != nil {
			appzaplog.Error("gorm open UserAvatarDB failed", zap.Error(err))
			return err
		}

		AttentionDB, err = gorm.Open("mysql", conf.AttentionDb)
		if err != nil {
			appzaplog.Error("gorm open AttentionDB failed", zap.Error(err))
			return err
		}
	} else {
		return errors.New("no appconfig find")
	}

	return nil
}
