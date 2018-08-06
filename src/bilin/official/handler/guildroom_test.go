package handler

import (
	"encoding/json"
	"testing"
	"time"
)

func TestGetGuildRoom(t *testing.T) {
	httprt := HttpRetComm{
		Desc: "success",
		Time: time.Now().Unix(),
		Data: &GuildRoom{
			GuildId: 100,
			RoomIds: []int64{911, 912, 918},
		},
	}
	bin, err := json.Marshal(httprt)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(bin))
}
