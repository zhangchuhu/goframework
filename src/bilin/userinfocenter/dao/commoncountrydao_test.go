package dao

import (
	"github.com/jinzhu/gorm"
	"testing"
)

func init() {
	var err error
	//UserDB, err = gorm.Open("mysql", "hujiao@HujiaoUser:MaUYa8w1Z@tcp(221.228.79.244:8066)/Hujiao?charset=utf8&parseTime=True&loc=Local")
	UserDB, err = gorm.Open("mysql", "bilin_admin:avYsLkYwQ@tcp(58.215.143.9:6307)/Hujiao?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic("init db failed")
	}
}
func TestGetCommonCountry(t *testing.T) {
	info, err := GetCommonCountry(10)
	if err != nil {
		t.Error("GetCommonCountry failed", err)
	}
	t.Logf("info:%+v", info)
}
