package dao

import (
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"fmt"
	"github.com/jinzhu/gorm"
)

type UserChatTag struct {
	gorm.Model
	FromUserID  int64 `gorm:"index"`
	ToUserID    int64 `gorm:"index"`
	ChatTags    string
	UpdateTimes int64
	TalkSecond  int64
	TagStatus   int64
}

const usertagtableoptions = "ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='用户聊天标签'"

func tableName(userid int64) string {
	return fmt.Sprintf("user_chat_tag_%d", userid%100)
}

func (t *UserChatTag) Create() error {
	tablename := tableName(t.ToUserID)
	if !hujiaoChatTagDB.HasTable(tablename) {
		if err := hujiaoChatTagDB.Table(tablename).Set("gorm:table_options", usertagtableoptions).CreateTable(&UserChatTag{}).Error; err != nil {
			return err
		}
	}

	if err := hujiaoChatTagDB.Table(tablename).Create(t).Error; err != nil {
		return err
	}
	return nil
}

func (t *UserChatTag) Update() error {
	tablename := tableName(t.ToUserID)
	if err := hujiaoChatTagDB.Table(tablename).Model(t).Updates(t).Error; err != nil {
		appzaplog.Error("UserChatTag Update err", zap.Error(err), zap.Any("req", t))
		return err
	}
	appzaplog.Debug("UserChatTag Update", zap.Any("resp", t))
	return nil
}

func (t *UserChatTag) Get() (*UserChatTag, error) {
	tablename := tableName(t.ToUserID)
	var ret UserChatTag
	if err := hujiaoChatTagDB.Table(tablename).Where(t).First(&ret).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) || IsTableNotExistErr(err) {
			return nil, nil
		}
		appzaplog.Error("UserChatTag Find err", zap.Error(err), zap.Any("req", t))
		return nil, err
	}
	return &ret, nil
}

func (t *UserChatTag) GetAll() ([]UserChatTag, error) {
	tablename := tableName(t.ToUserID)
	var ret []UserChatTag
	if err := hujiaoChatTagDB.Table(tablename).Where(t).Find(&ret).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) || IsTableNotExistErr(err) {
			return ret, nil
		}
		appzaplog.Error("UserChatTag Find err", zap.Error(err), zap.Any("req", t))
		return ret, err
	}
	return ret, nil
}

func BatchUserChatTag(toUserIds []int64) (map[int64][]UserChatTag, error) {
	var retmap = make(map[int64][]UserChatTag)
	for _, v := range toUserIds {
		tablename := tableName(v)
		var ret []UserChatTag
		if err := hujiaoChatTagDB.Table(tablename).Where(&UserChatTag{ToUserID: v}).Find(&ret).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) || IsTableNotExistErr(err) {
				continue
			}
			appzaplog.Error("UserChatTag Find err", zap.Error(err), zap.Int64s("uids", toUserIds))
			return retmap, err
		}
		retmap[v] = ret
	}
	return retmap, nil
}

func (t *UserChatTag) Del() error {
	if err := hujiaoChatTagDB.Model(t).Delete(t).Error; err != nil {
		appzaplog.Error("ChatTag Del err", zap.Error(err), zap.Any("req", t))
		return err
	}
	appzaplog.Debug("ChatTag Del", zap.Any("resp", t))
	return nil
}
