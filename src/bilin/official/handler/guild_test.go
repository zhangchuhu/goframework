package handler

import (
	"encoding/json"
	"testing"
	"time"
)

func TestHandleGuildReq(t *testing.T) {
	httpret := HttpRetComm{
		Desc: "success",
		Time: time.Now().Unix(),
		Data: &GuildInfo{
			GuildID:   1100,
			Title:     "以工会",
			Mobile:    "123456789",
			Desc:      "最牛工会",
			GuildLogo: "http://1.jpg",
		},
	}

	bin, err := json.Marshal(httpret)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(bin))
}
