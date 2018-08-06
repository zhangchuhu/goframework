package dao_test

import (
	"bilin/userinfocenter/dao"
	_ "code.yy.com/yytars/goframework/tars/servant"
	"testing"
)

func TestGetUserInfo(t *testing.T) {
	//err, info := dao.GetUserInfo(17794899)
	info, err := dao.GetUserInfo(17796525)
	if err != nil {
		t.Error("GetUserInfo error:" + err.Error())
	} else {
		t.Logf("GetUserInfo success, info: %v", info)
	}
}

func TestGetAvatatrUsers(t *testing.T) {
	info, err := dao.GetAvatatrUsers(201379, 100)
	if err != nil {
		t.Error("GetAvatatrUsers error:" + err.Error())
	} else {
		t.Logf("GetAvatatrUsers success, info: %v", info)
	}
}
