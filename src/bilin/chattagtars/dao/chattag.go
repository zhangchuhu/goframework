package dao

import (
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"github.com/jinzhu/gorm"
)

type ChatTag struct {
	gorm.Model
	TagName  string
	TagColor string
}

const tagtableoptions = "ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='用户聊天标签'"

func (t *ChatTag) Create() error {
	if !hujiaoChatTagDB.HasTable(&ChatTag{}) {
		if err := hujiaoChatTagDB.Set("gorm:table_options", tagtableoptions).CreateTable(&ChatTag{}).Error; err != nil {
			return err
		}
	}

	if err := hujiaoChatTagDB.Create(t).Error; err != nil {
		appzaplog.Error("ChatTag Create err", zap.Error(err))
		return err
	}
	return nil
}

func GetAll() ([]ChatTag, error) {
	var ret []ChatTag
	if err := hujiaoChatTagDB.Find(&ret).Error; err != nil {
		if IsTableNotExistErr(err) || gorm.IsRecordNotFoundError(err) {
			return ret, nil
		}
		appzaplog.Error("ChatTag Find err", zap.Error(err))
		return ret, err
	}
	return ret, nil
}

func (t *ChatTag) Update() error {
	if err := hujiaoChatTagDB.Model(t).Updates(t).Error; err != nil {
		appzaplog.Error("ChatTag Update err", zap.Error(err), zap.Any("req", t))
		return err
	}
	appzaplog.Debug("ChatTag Update", zap.Any("resp", t))
	return nil
}

func (t *ChatTag) Del() error {
	if err := hujiaoChatTagDB.Model(t).Delete(t).Error; err != nil {
		appzaplog.Error("ChatTag Del err", zap.Error(err), zap.Any("req", t))
		return err
	}
	appzaplog.Debug("ChatTag Del", zap.Any("resp", t))
	return nil
}
