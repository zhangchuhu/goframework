/*
 * Copyright (c) 2018-07-26.
 * Author: kordenlu
 * 功能描述:${<VARIABLE_NAME>}
 */

package dao

import (
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"github.com/jinzhu/gorm"
)

type OAMUser struct {
	gorm.Model
	Username string
	Passwd   string
	Role     int64
}

const oamusertableoptions = "ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='管理后台用户信息'"

func (p *OAMUser) Create() error {
	if !hujiaoRecDB.HasTable(&OAMUser{}) {
		if err := hujiaoRecDB.Set("gorm:table_options", tableoptions).CreateTable(&OAMUser{}).Error; err != nil {
			return err
		}
	}

	if err := hujiaoRecDB.Create(p).Error; err != nil {
		return err
	}
	return nil
}

func (g *OAMUser) Get() ([]OAMUser, error) {
	var ret []OAMUser
	if err := hujiaoRecDB.Where(g).Find(&ret).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		appzaplog.Error("GuildRoom Find err", zap.Error(err), zap.Any("req", g))
		return ret, err
	}
	return ret, nil
}
