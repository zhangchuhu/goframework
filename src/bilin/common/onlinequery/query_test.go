package onlinequery

import (
	"testing"
)

func TestGetUserRoom(t *testing.T) {
	uid := int64(88889999)
	rid, err := GetUserRoom(uid)
	if err != nil {
		t.Error(err)
	}
	t.Logf("uid: %v, rid: %v\n", uid, rid)
}

func TestGetRoomUser(t *testing.T) {
	rid := int64(400000414)
	uid, err := GetRoomUser(rid)
	if err != nil {
		t.Error(err)
	}
	t.Logf("uid: %v, rid: %v\n", uid, rid)
}

func TestUserCount(t *testing.T) {
	cnt, err := UserCount()
	if err != nil {
		t.Error(err)
	}
	t.Logf("count: %v\n", cnt)
}
