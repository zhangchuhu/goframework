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
	hujiaoChatTagDB, err = gorm.Open("mysql", "bilin_admin:avYsLkYwQ@tcp(58.215.143.9:6307)/Hujiao?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		appzaplog.Error("gorm open ChatTagDb failed", zap.Error(err))
		return
	}
	hujiaoChatTagDB.LogMode(true)

	os.Exit(m.Run())
}
