package dao

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"testing"
)

func init() {
	var err error
	//hujiaoCallDB, err = gorm.Open("mysql", "bilin_admin:avYsLkYwQ@tcp(58.215.143.9:6307)/Hujiao?charset=utf8&parseTime=True&loc=Local")
	hujiaoRecDB, err = gorm.Open("mysql", "bilin_admin:avYsLkYwQ@tcp(58.215.143.9:6307)/Hujiao?charset=utf8&parseTime=True&loc=Local")
	//hujiaoRecDB, err = gorm.Open("mysql", "blguildadmin:1oCrqMZYwr@tcp(14.116.173.148:6305)/Hujiao?readTimeout=500ms&charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
	}
}

func TestHostRec_Create(t *testing.T) {
	hostrec := &HostRec{
		TypeId: 1,
		HostID: 100,
	}
	if err := hostrec.Create(); err != nil {
		t.Error(err)
	}
	t.Log(hostrec)
}

func TestGetHostRec(t *testing.T) {
	info, err := GetHostRec()
	if err != nil {
		t.Error(err)
	}
	t.Log(info)
}

func TestHostRec_Update(t *testing.T) {
	hostrec := HostRec{
		HostID: 200,
		TypeId: 1000,
	}
	hostrec.ID = 1
	if err := hostrec.Update(); err != nil {
		t.Error(err)
	}
}
