package clientcenter_test

import (
	"bilin/clientcenter"
	"testing"
	"time"
)

const (
	testEnvUid = 17796525
	proEnvUid  = 29481968
)

func TestVerifyAccessToken(t *testing.T) {
	pass, err := clientcenter.VerifyAccessToken("whatapp", "100")
	if err != nil {
		t.Error(err)
	}
	t.Log(pass)
}

func TestGetUserInfoByUserId(t *testing.T) {
	info, err := clientcenter.GetUserInfoByUserId(17796250)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", info)
}

func TestGetUserHobbies(t *testing.T) {
	info, err := clientcenter.GetUserHobbies(proEnvUid)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", info)
}

func TestQueryDynamicAndDynamicAttrListByUserByPage(t *testing.T) {
	info, err := clientcenter.QueryDynamicAndDynamicAttrListByUserByPage(proEnvUid, proEnvUid, time.Now().Unix()*1000, 1)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", info)
}

func TestGetBilinKtvById(t *testing.T) {
	info, found := clientcenter.GetBilinKtvById(131)
	t.Logf("%+v    found %v", info, found)
}

func TestQueryFriendList(t *testing.T) {
	info := clientcenter.QueryFriendList(40373825)
	t.Logf("%+v", info)
}

func TestQueryAttentionList(t *testing.T) {
	info := clientcenter.QueryAttentionList(17795058)
	t.Logf("%+v", info)
}

func TestGetUserByBLId(t *testing.T) {
	info, found := clientcenter.GetUserByBLId(40373)
	t.Logf("%+v    found %v", info, found)
}
