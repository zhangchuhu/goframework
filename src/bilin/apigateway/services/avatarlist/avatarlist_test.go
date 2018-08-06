package avatarlist_test

import (
	"bilin/apigateway/services/avatarlist"
	_ "code.yy.com/yytars/goframework/tars/servant"
	"testing"
	//"time"
)

func TestGetAvatarList(t *testing.T) {
	info := avatarlist.GetAvatarList()
	if info == nil {
		t.Logf("GetAvatarList nil")
	} else {
		t.Logf("GetAvatarList success, info: %v", info)
	}
}
