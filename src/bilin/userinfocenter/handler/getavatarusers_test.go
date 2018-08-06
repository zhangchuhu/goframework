package handler_test

import (
	"bilin/userinfocenter/handler"
	_ "code.yy.com/yytars/goframework/tars/servant"
	"testing"
)

func TestGetRondomAvatarUsers(t *testing.T) {

	if list, err := handler.GetRondomAvatarUsers(); err != nil {
		t.Error("GetRondomAvatarUsers error:" + err.Error())
	} else {
		t.Logf("GetRondomAvatarUsers success, list:%v", list)
	}
}
