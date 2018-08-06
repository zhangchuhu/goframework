package dao

import (
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"fmt"
	"github.com/jinzhu/gorm"
)

type PuaTopic struct {
	gorm.Model
	Topic string
}

const puatableoptions = "ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='撩妹套话'"

func (t *PuaTopic) Create() error {
	if !hujiaoChatTagDB.HasTable(&PuaTopic{}) {
		if err := hujiaoChatTagDB.Set("gorm:table_options", puatableoptions).CreateTable(&PuaTopic{}).Error; err != nil {
			return err
		}
	}

	if err := hujiaoChatTagDB.Create(t).Error; err != nil {
		appzaplog.Error("ChatTag Create err", zap.Error(err))
		return err
	}
	return nil
}

func GetAllPuaTopic() ([]PuaTopic, error) {
	var ret []PuaTopic
	if err := hujiaoChatTagDB.Find(&ret).Error; err != nil {
		if IsTableNotExistErr(err) || gorm.IsRecordNotFoundError(err) {
			return ret, nil
		}
		appzaplog.Error("GetAllPuaTopic Find err", zap.Error(err))
		return ret, err
	}
	return ret, nil
}

func GetPuaTopicByPage(page, pagesize int64) ([]PuaTopic, int64, error) {
	var ret []PuaTopic
	var count int64
	if err := hujiaoChatTagDB.Model(&PuaTopic{}).Count(&count).Error; err != nil {
		if IsTableNotExistErr(err) || gorm.IsRecordNotFoundError(err) {
			return ret, 0, nil
		}
		appzaplog.Error("GetAllPuaTopic Find err", zap.Error(err))
		return ret, 0, err
	}
	if page < 1 {
		return ret, 0, fmt.Errorf("page too small")
	}
	offset := (page - 1) * pagesize
	if err := hujiaoChatTagDB.Offset(offset).Limit(pagesize).Find(&ret).Order("id DESC").Error; err != nil {
		appzaplog.Error("GetAllPuaTopic Find err", zap.Error(err))
		return ret, count, err
	}
	return ret, count, nil
}

func (t *PuaTopic) Update() error {
	if err := hujiaoChatTagDB.Model(t).Updates(t).Error; err != nil {
		appzaplog.Error("PuaTopic Update err", zap.Error(err), zap.Any("req", t))
		return err
	}
	appzaplog.Debug("PuaTopic Update", zap.Any("resp", t))
	return nil
}

func (t *PuaTopic) Del() error {
	if err := hujiaoChatTagDB.Model(t).Delete(t).Error; err != nil {
		appzaplog.Error("PuaTopic Del err", zap.Error(err), zap.Any("req", t))
		return err
	}
	appzaplog.Debug("PuaTopic Del", zap.Any("resp", t))
	return nil
}

func GetTopicNotIn(ids []int64, limit int) ([]PuaTopic, error) {
	var ret []PuaTopic
	if err := hujiaoChatTagDB.Not(ids).Limit(limit).Find(&ret).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) || IsTableNotExistErr(err) {
			return nil, nil
		}
		appzaplog.Error("GetTopicNotIn err", zap.Error(err), zap.Int64s("ids", ids))
		return nil, err
	}
	return ret, nil
}
