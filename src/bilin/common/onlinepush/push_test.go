package onlinepush

import (
	"bilin/protocol"
	"testing"
)

func TestPushToUser(t *testing.T) {
	var mpush bilin.MultiPush
	mpush.UserIDs = []int64{88889999}
	mpush.Msg = &bilin.ServerPush{
		MessageType: 10086,
		PushBuffer:  []byte("China Mobile"),
	}
	offline, err := PushToUser(mpush)
	if err != nil {
		t.Error(err)
	}
	t.Logf("offline: %v\n", offline)
}

func TestPushToRoom(t *testing.T) {
	push := bilin.ServerPush{
		MessageType: 10086,
		PushBuffer:  []byte("China Mobile X"),
	}
	err := PushToRoom(10, push)
	if err != nil {
		t.Error(err)
	}
}
