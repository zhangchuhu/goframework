package dao

import (
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"strconv"
	"time"
)

const (
//CAROUSEL = 6
//ANCHOR_TYPE = 19
)

type Carousel struct {
	ID           int64  `gorm:"AUTO_INCREMENT;primary_key;column:ID"`
	BackgroudURL string `gorm:"not null;column:BACKGROUD_URL"`
	TargetType   int32  `gorm:"not null;column:TARGET_TYPE"`
	TargetURL    string `gorm:"not null;column:TARGET_URL"`
	StartTime    int64  `gorm:"not null;column:START_TIME"`
	EndTime      int64  `gorm:"not null;column:END_TIME"`
	Channel      string `gorm:"column:CHANNEL" sql:"DEFAULT:''"`
	Version      string `gorm:"column:VERSION" sql:"DEFAULT:''"`
	ForUserType  int32  `gorm:"column:FOR_USER_TYPE" sql:"DEFAULT:0"`
	Sort         int32  `gorm:"column:SORT" sql:"DEFAULT:0"`
	Width        int32  `gorm:"column:WIDTH"`
	Height       int32  `gorm:"column:HEIGHT"`
	Position     int32  `gorm:"column:POSITION"`
	HotLineType  string `gorm:"column:HOT_LINE_TYPE"`
}

func (Carousel) TableName() string {
	return "BANNER"
}

func GetCarousel() ([]*Carousel, error) {
	var ret []*Carousel
	cur_time := time.Now().Unix() * 1000
	cur_time_str := strconv.FormatInt(cur_time, 10)
	condition := "IS_DELETE = 0 and TYPE = 6 and START_TIME < " + cur_time_str + " and END_TIME >" + cur_time_str
	db_ := hujiaoUserDB.Where(condition).Find(&ret)
	if db_.RecordNotFound() {
		return nil, nil
	}
	return ret, db_.Error
}

func Banner(typid int64) ([]*Carousel, error) {
	var ret []*Carousel
	cur_time := time.Now().Unix() * 1000
	condition := fmt.Sprintf("IS_DELETE = 0 and TYPE = %d and START_TIME < %d and END_TIME >%d", typid, cur_time, cur_time)
	//condition := "IS_DELETE = 0 and TYPE = 6 and START_TIME < " + cur_time_str + " and END_TIME >" + cur_time_str
	appzaplog.Debug("LivingBanner", zap.String("condition", condition))
	db_ := hujiaoUserDB.Where(condition).Find(&ret)
	if db_.RecordNotFound() {
		return nil, nil
	}
	return ret, db_.Error
}

func LivingBanner() ([]*Carousel, error) {
	return Banner(10)
}
