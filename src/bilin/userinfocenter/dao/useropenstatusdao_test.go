package dao_test

import (
	"bilin/userinfocenter/dao"
	_ "code.yy.com/yytars/goframework/tars/servant"
	"testing"
	//"time"
)

func TestGetUserOpenStaus(t *testing.T) {
	dao.InitThriftConnentPool("221.228.91.178:9090;")
	//dao.InitThriftConnentPool("172.26.64.16:9090;")
	info, err := dao.GetUserOpenStaus(17795316,"5.0.0","ios","127.0.0.1")
	if err != nil {
		t.Error("GetUserOpenStaus error:" + err.Error())
	} else {
		t.Logf("GetUserOpenStaus success, info: %v", info)
	}
}
