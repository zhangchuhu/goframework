package cache

import (
	"encoding/json"
	"sort"
	"testing"
)

func TestNewHotLineList(t *testing.T) {
	reclivinglist := []*RecommandLivingInfo{
		&RecommandLivingInfo{
			LiveID:    111,
			Title:     "第一家",
			UserCount: 110,
			LivingHostInfo: LivingHostInfo{
				UserID:   111,
				NickName: "人人",
				City:     "广州",
				SmallURL: "http://111.com/1.jpg",
				TagUrl:   []string{"http://tag.yy.com/1.jpg", "http://tag.yy.com/2.jpg"},
			},
		},
		&RecommandLivingInfo{
			LiveID:    121,
			Title:     "第er家",
			UserCount: 99,
			LivingHostInfo: LivingHostInfo{
				UserID:   120,
				NickName: "刚看过看",
				City:     "海口",
				SmallURL: "http://111.com/2.jpg",
			},
		},
	}

	homelist := &HttpRetComm{
		IsEncrypt: "false",
		Data: HttpRetDataComm{
			Code: 0,
			Msg:  "success",
			Body: RecommandLivingBody{
				LastPage:            "true",
				RecommandLivingList: reclivinglist,
			},
		},
	}

	byte, err := json.Marshal(homelist)
	if err != nil {
		t.Error("failed to marshal", err)
	}
	t.Logf("byte:%v", string(byte))
}

func TestToSmall(t *testing.T) {
	small := tosmall("https://img.inbilin.com/17795865/17795865_1523699808521.jpg-big")
	t.Logf("small:%v", small)
}

func TestSortByWeight(t *testing.T) {
	var normalhost []*RecommandLivingInfo
	normalhost = append(normalhost, &RecommandLivingInfo{
		SortWeight: 100,
	},
		&RecommandLivingInfo{
			SortWeight: 120,
		},
		&RecommandLivingInfo{
			SortWeight: 110,
		},
	)

	sort.Sort(RecommandLivingInfoSlice(normalhost))
	t.Log(normalhost)
}

func TestGetRec(t *testing.T) {
	info := []*RecommandLivingInfo{
		&RecommandLivingInfo{
			LiveID:            1,
			UserCount:         100,
			Last5MinRankValue: 10,
		},
		&RecommandLivingInfo{
			LiveID:    2,
			UserCount: 200,
		},
		&RecommandLivingInfo{
			LiveID:    3,
			UserCount: 300,
		},
		&RecommandLivingInfo{
			LiveID:    4,
			UserCount: 400,
		},
	}
	getRec(info)

}
