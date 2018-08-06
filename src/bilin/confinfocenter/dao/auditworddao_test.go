package dao

import (
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"github.com/jinzhu/gorm"
	"testing"
)

func init() {
	var err error
	hujiaoDB, err = gorm.Open("mysql", "bilin_admin:avYsLkYwQ@tcp(58.215.143.9:6307)/HujiaoMsg?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		appzaplog.Error("hujiaoDB open failed", zap.Error(err))
	}
}

func TestGetAuditWorld(t *testing.T) {
	if word, err := GetAuditWorld(); err != nil {
		t.Error(err)
	} else {
		t.Log(word)
	}
}
