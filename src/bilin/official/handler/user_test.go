package handler

import (
	"encoding/json"
	"testing"
	"time"
)

func TestHandleHostReq(t *testing.T) {
	ret := HttpRetComm{
		Desc: "success",
		Time: time.Now().Unix(),
		Data: HostInfo{
			UID:           100,
			BILINNumber:   200,
			NickName:      "晓晓",
			Avatar:        "https://test.jpg",
			TotalCharmNum: 100,
			TotalHeartNum: 1000,
			FansNum:       500,
			AttentionNum:  10,
		},
	}
	bin, err := json.Marshal(ret)
	if err != nil {
		t.Error("json marshal fialed", err)
	}
	t.Log(string(bin))
}
