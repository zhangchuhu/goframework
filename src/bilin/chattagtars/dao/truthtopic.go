package dao

import (
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"fmt"
	"github.com/jinzhu/gorm"
)

type TruthTopic struct {
	gorm.Model
	Topic string
}

const truthtableoptions = "ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='真心话'"

func (t *TruthTopic) Create() error {
	if !hujiaoChatTagDB.HasTable(&TruthTopic{}) {
		if err := hujiaoChatTagDB.Set("gorm:table_options", truthtableoptions).CreateTable(&TruthTopic{}).Error; err != nil {
			return err
		}
	}

	if err := hujiaoChatTagDB.Create(t).Error; err != nil {
		appzaplog.Error("ChatTag Create err", zap.Error(err))
		return err
	}
	return nil
}

func GetAllTruthTopic() ([]TruthTopic, error) {
	var ret []TruthTopic
	if err := hujiaoChatTagDB.Find(&ret).Error; err != nil {
		if IsTableNotExistErr(err) || gorm.IsRecordNotFoundError(err) {
			return ret, nil
		}
		appzaplog.Error("GetAllPuaTopic Find err", zap.Error(err))
		return ret, err
	}
	return ret, nil
}

func GetAllTruthTopicByPage(page, pagesize int64) ([]TruthTopic, int64, error) {
	var ret []TruthTopic
	var count int64
	if err := hujiaoChatTagDB.Model(&TruthTopic{}).Count(&count).Error; err != nil {
		if IsTableNotExistErr(err) || gorm.IsRecordNotFoundError(err) {
			return ret, 0, nil
		}
		appzaplog.Error("GetAllTruthTopicByPage Find err", zap.Error(err))
		return ret, 0, err
	}
	if page < 1 {
		return ret, 0, fmt.Errorf("page too small")
	}
	offset := (page - 1) * pagesize
	if err := hujiaoChatTagDB.Offset(offset).Limit(pagesize).Find(&ret).Order("id DESC").Error; err != nil {
		appzaplog.Error("GetAllTruthTopicByPage Find err", zap.Error(err))
		return ret, count, err
	}
	return ret, count, nil
}

func (t *TruthTopic) Update() error {
	if err := hujiaoChatTagDB.Model(t).Updates(t).Error; err != nil {
		appzaplog.Error("TruthTopic Update err", zap.Error(err), zap.Any("req", t))
		return err
	}
	appzaplog.Debug("TruthTopic Update", zap.Any("resp", t))
	return nil
}

func (t *TruthTopic) Del() error {
	if err := hujiaoChatTagDB.Model(t).Delete(t).Error; err != nil {
		appzaplog.Error("TruthTopic Del err", zap.Error(err), zap.Any("req", t))
		return err
	}
	appzaplog.Debug("TruthTopic Del", zap.Any("resp", t))
	return nil
}

func GetTruthTopicNotIn(ids []int64, limit int) ([]TruthTopic, error) {
	var ret []TruthTopic
	if err := hujiaoChatTagDB.Not(ids).Limit(limit).Find(&ret).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) || IsTableNotExistErr(err) {
			return nil, nil
		}
		appzaplog.Error("GetTopicNotIn err", zap.Error(err), zap.Int64s("ids", ids))
		return nil, err
	}
	return ret, nil
}
