package dao

import (
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"github.com/jinzhu/gorm"
	"time"
)

type Stickie struct {
	gorm.Model
	TypeId    int64     `gorm:"not null;column:TypeId"`
	RoomID    uint64    `gorm:"not null;column:RoomID" sql:"DEFAULT:0"`
	StartTime time.Time `gorm:"not null;column:START_TIME"`
	EndTime   time.Time `gorm:"not null;column:END_TIME"`
	Weight    int64     `gorm:"not null column:Weight" sql:"DEFAULT:0"`
}

func (Stickie) TableName() string {
	return "STICKYPOST"
}

func GetStickie() ([]Stickie, error) {
	var ret []Stickie
	db_ := hujiaoRecDB.Find(&ret)
	if db_.RecordNotFound() {
		return []Stickie{}, nil
	}
	return ret, db_.Error
}

const tableoptions = "ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='分类置顶推荐'"

func (p *Stickie) Create() error {
	if !hujiaoRecDB.HasTable(&Stickie{}) {
		if err := hujiaoRecDB.Set("gorm:table_options", tableoptions).CreateTable(&Stickie{}).Error; err != nil {
			return err
		}
	}

	if err := hujiaoRecDB.Create(p).Error; err != nil {
		return err
	}
	return nil
}

func (p *Stickie) Update() error {
	if err := hujiaoRecDB.Model(p).Updates(p).Error; err != nil {
		appzaplog.Error("Stickie Update err", zap.Error(err), zap.Any("req", p))
		return err
	}
	appzaplog.Debug("Stickie Update", zap.Any("resp", p))
	return nil
}

func (p *Stickie) Del() error {
	if err := hujiaoRecDB.Delete(p).Error; err != nil {
		appzaplog.Error("DelStickie Update err", zap.Error(err), zap.Any("stikcy", p))
		return err
	}
	appzaplog.Debug("DelStickie Update", zap.Any("stikcy", p))
	return nil
}
