/*
 * Copyright (c) 2018-07-23.
 * Author: kordenlu
 * 功能描述:${<VARIABLE_NAME>}
 */

package dao

import (
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"github.com/jinzhu/gorm"
)

type ContractFlow struct {
	gorm.Model
	OperationUid int64
	ContractId   int64
	GuildId      int64
	HostUid      int64
	OwUid        int64
	Operation    int32 // 1 接收， 2 拒绝
}

const contractflowtableoptions = "ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='签约操作记录'"

func (c *ContractFlow) Create() error {
	if !hujiaoRecDB.HasTable(&ContractFlow{}) {
		if err := hujiaoRecDB.Set("gorm:table_options", contractflowtableoptions).CreateTable(&ContractFlow{}).Error; err != nil {
			appzaplog.Error("Create ContractFlow Table err", zap.Error(err))
			return err
		}
	}

	if db_ := hujiaoRecDB.Create(c); db_.Error != nil {
		appzaplog.Error("Create ContractFlow err", zap.Error(db_.Error))
		return db_.Error
	}

	return nil
}

func (c *ContractFlow) Get() ([]ContractFlow, error) {
	var ret []ContractFlow
	if err := hujiaoRecDB.Where(c).Find(&ret).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) || IsTableNotExistErr(err) {
			return ret, nil
		}
		appzaplog.Error("ContractFlow.Get err", zap.Error(err), zap.Any("filter", c))
		return ret, err
	}
	return ret, nil
}

func (c *ContractFlow) Update() error {
	if err := hujiaoRecDB.Model(c).Updates(c).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) || IsTableNotExistErr(err) {
			return nil
		}
		appzaplog.Error("ContractFlow.Updates err", zap.Error(err), zap.Any("filter", c))
		return err
	}
	return nil
}

func (c *ContractFlow) Delete() error {
	if err := hujiaoRecDB.Where(c).Delete(c).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) || IsTableNotExistErr(err) {
			return nil
		}
		appzaplog.Error("ContractFlow.Delete err", zap.Error(err), zap.Any("filter", c))
		return err
	}
	return nil
}
