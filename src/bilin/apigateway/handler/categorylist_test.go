package handler

import (
	"bilin/apigateway/cache"
	"encoding/json"
	"testing"
)

func TestNewCategoryList(t *testing.T) {
	livecategory := []cache.LivingCategory{
		cache.LivingCategory{
			CategoryID:      100,
			CategoryName:    "热门",
			FontColor:       "#fffff",
			BackgroundImage: "http://live1.1/1.pic",
		},
		cache.LivingCategory{
			CategoryID:      101,
			CategoryName:    "音乐",
			FontColor:       "#fffff",
			BackgroundImage: "http://live1.1/1.pic",
		},
	}

	homelist := &cache.HttpRetComm{
		IsEncrypt: "false",
		Data: cache.HttpRetDataComm{
			Code: 0,
			Msg:  "success",
			Body: cache.LivingCategoryBody{
				LivingCategoryList: livecategory,
			},
		},
	}

	byte, err := json.Marshal(homelist)
	if err != nil {
		t.Error("failed to marshal", err)
	}
	t.Logf("byte:%v", string(byte))
}
