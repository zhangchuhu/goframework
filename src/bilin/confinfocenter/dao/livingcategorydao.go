package dao

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

type AppModel struct {
	ID        uint      `gorm:"AUTO_INCREMENT;primary_key;column:ID"`
	CREATE_ON time.Time `gorm:"column:CREATE_ON"`
	UPDATE_ON time.Time `gorm:"column:UPDATE_ON"`
	CREATE_BY string    `gorm:"column:CREATE_BY"`
	UPDATE_BY string    `gorm:"column:UPDATE_BY"`
}
type LivingCategory struct {
	AppModel
	TYPE_ID          int    `gorm:"unique;not null;column:TYPE_ID"`
	TYPE_NAME        string `gorm:"not null;column:TYPE_NAME"`
	FONT_COLOR       string `gorm:"not null;column:FONT_COLOR"`
	BACKGROUND_IMAGE string `gorm:"not null;column:BACKGROUND_IMAGE"`
	SORT             int    `gorm:"column:SORT" sql:"DEFAULT:0"`
	DEL_FLAG         int    `gorm:"column:DEL_FLAG" sql:"DEFAULT:0"`
}

func (LivingCategory) TableName() string {
	return "HOTLINE_DIRECT_TYPE"
}

func GetLivingCategorys() ([]LivingCategory, error) {
	if hujiaoCallDB == nil {
		return nil, dbNotInitErr
	}
	var ret []LivingCategory
	db_ := hujiaoCallDB.Where("DEL_FLAG = 0").Find(&ret)
	if db_.RecordNotFound() {
		return []LivingCategory{}, nil
	}
	return ret, db_.Error
}
