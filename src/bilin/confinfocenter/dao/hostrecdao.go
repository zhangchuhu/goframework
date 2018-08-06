package dao

import (
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"github.com/jinzhu/gorm"
)

type HostRec struct {
	gorm.Model
	TypeId uint64 `gorm:"not null;column:TypeId"`
	HostID uint64 `gorm:"not null;column:HostID" sql:"DEFAULT:0`
}

func (HostRec) TableName() string {
	return "HOSTREC"
}

func GetHostRec() ([]HostRec, error) {

	var ret []HostRec
	db_ := hujiaoRecDB.Find(&ret)
	if db_.RecordNotFound() {
		return []HostRec{}, nil
	}
	return ret, db_.Error
}

const hostrectableoptions = "ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='品类主持推荐表'"

func (p *HostRec) Create() error {
	if !hujiaoRecDB.HasTable(p.TableName()) {
		if err := hujiaoRecDB.Set("gorm:table_options", hostrectableoptions).CreateTable(p).Error; err != nil {
			return err
		}
	}

	if err := hujiaoRecDB.Create(p).Error; err != nil {
		appzaplog.Error("Create hostRec err", zap.Error(err), zap.Any("req", p))
		return err
	}
	return nil
}

func (p *HostRec) Update() error {
	if err := hujiaoRecDB.Model(p).Updates(p).Error; err != nil {
		appzaplog.Error("Update hostRec err", zap.Error(err), zap.Any("req", p))
		return err
	}
	return nil
}

func (p *HostRec) Del() error {
	if err := hujiaoRecDB.Delete(p).Error; err != nil {
		appzaplog.Error("Del hostRec err", zap.Error(err), zap.Any("req", p))
		return err
	}
	return nil
}
