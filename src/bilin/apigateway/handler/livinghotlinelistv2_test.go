package handler

import (
	"bilin/apigateway/cache"
	"encoding/json"
	"strconv"
	"testing"
)

func TestHandleGetLivingHotLineListV2(t *testing.T) {

	body := &HotLineListV2{
		LastPage: "true",
	}

	body.RecommandLivingList = append(body.RecommandLivingList,
		&cache.RecommandLivingInfo{
			TypeID:    1,
			LiveID:    111,
			Title:     "第一家",
			UserCount: 110,
			LivingHostInfo: cache.LivingHostInfo{
				UserID:   111,
				NickName: "人人",
				City:     "广州",
				SmallURL: "http://111.com/1.jpg",
				TagUrl:   []string{"http://tag.yy.com/1.jpg", "http://tag.yy.com/2.jpg"},
			},
		},
	)
	body.RecommandLivingList = append(body.RecommandLivingList,
		&RecBanner{
			TypeID:     2,
			BGURL:      "http://1.txt",
			TargetType: 10,
			TargetURL:  "9100",
		},
	)
	body.RecommandLivingList = append(body.RecommandLivingList,
		&cache.RecommandLivingInfo{
			TypeID:    1,
			LiveID:    121,
			Title:     "第er家",
			UserCount: 99,
			LivingHostInfo: cache.LivingHostInfo{
				UserID:   120,
				NickName: "刚看过看",
				City:     "海口",
				SmallURL: "http://111.com/2.jpg",
			},
		},
	)
	homelist := &HttpRetComm{
		IsEncrypt: "false",
		Data: HttpRetDataComm{
			Code: 0,
			Msg:  "success",
			Body: body,
		},
	}

	byte, err := json.Marshal(homelist)
	if err != nil {
		t.Error("failed to marshal", err)
	}
	t.Logf("byte:%v", string(byte))
}

func TestRoomId(t *testing.T) {
	realroomid := int64(100)
	targetUrl := "inbilin://live/hotline?hotlineId=" + strconv.FormatInt(realroomid, 10)
	roomid := RoomId(targetUrl)
	if roomid != realroomid {
		t.Error("roomid should be 100")
	}
	t.Log(roomid)
}
