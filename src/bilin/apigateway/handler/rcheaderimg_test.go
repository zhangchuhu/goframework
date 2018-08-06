package handler

import (
	"bilin/apigateway/cache"
	"encoding/json"
	"testing"
)

func TestHandleRCHeaderImg(t *testing.T) {
	homelist := &cache.HttpRetComm{
		IsEncrypt: "false",
		Data: cache.HttpRetDataComm{
			Code: 0,
			Msg:  "success",
			Body: RCHeaderImgBody{
				SmallImgUrls: []string{
					"http://1.jpg",
					"http://2.jpg",
					"http://3.jpg",
				},
			},
		},
	}

	byte, err := json.Marshal(homelist)
	if err != nil {
		t.Error("failed to marshal", err)
	}
	t.Logf("byte:%v", string(byte))
}
