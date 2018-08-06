package dao

import (
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"github.com/jinzhu/gorm"
)

type GuildRec struct {
	gorm.Model
	TypeId uint64 `gorm:"not null;column:TypeId"`
	RoomID uint64 `gorm:"not null;column:RoomID" sql:"DEFAULT:0`
}

func (GuildRec) TableName() string {
	return "GUILDREC"
}

const guildrectableoptions = "ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='品类频道热门推荐表'"

func (p *GuildRec) Create() error {
	if !hujiaoRecDB.HasTable(p.TableName()) {
		if err := hujiaoRecDB.Set("gorm:table_options", guildrectableoptions).CreateTable(p).Error; err != nil {
			return err
		}
	}

	if err := hujiaoRecDB.Create(p).Error; err != nil {
		return err
	}
	return nil
}

func GetGuildRec() ([]GuildRec, error) {
	var ret []GuildRec
	db_ := hujiaoRecDB.Find(&ret)
	if db_.RecordNotFound() {
		return []GuildRec{}, nil
	}
	return ret, db_.Error
}

func UpdateGuildRec(p *GuildRec) error {
	gm := &GuildRec{}
	gm.ID = p.ID
	if err := hujiaoRecDB.Model(gm).Updates(p).Error; err != nil {
		appzaplog.Error("UpdateGuildRec update err", zap.Error(err), zap.Any("req", p))
		return err
	}
	return nil
}

func DelGuildRec(id uint) error {
	delguild := &GuildRec{}
	delguild.ID = id
	if err := hujiaoRecDB.Delete(delguild).Error; err != nil {
		return err
	}
	return nil
}
