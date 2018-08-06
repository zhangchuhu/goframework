package pushproxy

import (
	"testing"

	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
)

func init() {
	log.InitAppLog()
}

func TestPushToUser(t *testing.T) {
	err := PushToUser(100, 100, "test push to user")
	if err != nil {
		t.Error(err)
	}
}

func TestPushToRoom(t *testing.T) {
	err := PushToRoom(100, 100, "test push to room")
	if err != nil {
		t.Error(err)
	}
}
