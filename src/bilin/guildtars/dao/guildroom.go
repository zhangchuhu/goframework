package dao

import (
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"github.com/jinzhu/gorm"
)

type GuildRoom struct {
	gorm.Model
	GuildID int64
	RoomID  int64
}

const gctableoptions = "ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='工会频道'"

func (g *GuildRoom) Create() error {
	if !hujiaoRecDB.HasTable(&GuildRoom{}) {
		if err := hujiaoRecDB.Set("gorm:table_options", gctableoptions).CreateTable(&GuildRoom{}).Error; err != nil {
			return err
		}
	}

	if err := hujiaoRecDB.Create(g).Error; err != nil {
		return err
	}
	return nil
}

func (g *GuildRoom) Del() error {
	if err := hujiaoRecDB.Delete(g).Error; err != nil {
		appzaplog.Error("GuildRoom Delete err", zap.Error(err), zap.Any("guildchannel", g))
		return err
	}
	appzaplog.Debug("GuildRoom Del", zap.Any("guildchannel", g))
	return nil
}

func GetGuildRoomS() ([]GuildRoom, error) {
	var ret []GuildRoom
	db_ := hujiaoRecDB.Find(&ret)
	if db_.RecordNotFound() {
		return []GuildRoom{}, nil
	}
	return ret, db_.Error
}

func (g *GuildRoom) Get() ([]GuildRoom, error) {
	var ret []GuildRoom
	if err := hujiaoRecDB.Where(g).Find(&ret).Error; err != nil {
		appzaplog.Error("GuildRoom Find err", zap.Error(err), zap.Any("req", g))
		return ret, err
	}
	return ret, nil
}

func (c *GuildRoom) Update() error {
	if err := hujiaoRecDB.Model(c).Updates(c).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) || IsTableNotExistErr(err) {
			return nil
		}
		appzaplog.Error("GuildRoom.Updates err", zap.Error(err), zap.Any("filter", c))
		return err
	}
	return nil
}

func (c *GuildRoom) Delete() error {
	if err := hujiaoRecDB.Where(c).Delete(c).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) || IsTableNotExistErr(err) {
			return nil
		}
		appzaplog.Error("GuildRoom.Delete err", zap.Error(err), zap.Any("filter", c))
		return err
	}
	return nil
}
