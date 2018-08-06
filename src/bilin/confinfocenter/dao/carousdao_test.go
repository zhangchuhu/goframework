package dao

import (
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"github.com/jinzhu/gorm"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	var err error
	hujiaoCallDB, err = gorm.Open("mysql", "bilin_admin:avYsLkYwQ@tcp(58.215.143.9:6307)/Hujiao?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		appzaplog.Error("gorm open HuJiaoCallDb failed", zap.Error(err))
		return
	}
	hujiaoUserDB, err = gorm.Open("mysql", "bilin_admin:avYsLkYwQ@tcp(58.215.143.9:6307)/Hujiao?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		appzaplog.Error("gorm open HuJiaoUserDb failed", zap.Error(err))
		return
	}
	os.Exit(m.Run())
}

func TestGetCarousel(t *testing.T) {

	info, err := GetCarousel()
	if err != nil {
		t.Error("GetCarousel error:" + err.Error())
	} else {
		t.Logf("GetCarousel success, info: %+v", info)
	}
}

func TestLivingBanner(t *testing.T) {
	info, err := Banner(7)
	if err != nil {
		t.Error(err)
	}
	t.Log(info)
}
