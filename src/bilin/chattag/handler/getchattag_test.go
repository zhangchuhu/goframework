package handler

import "testing"

func TestTopNTagListByUser(t *testing.T) {
	in := &ChatTagList{
		ChatTags: []ChatTag{
			{TagName: "top1", TotalTagNum: 2000},
			{TagName: "top3", TotalTagNum: 1090},
			{TagName: "top2", TotalTagNum: 1999},
			{TagName: "top4", TotalTagNum: 1009},
		},
	}
	info := topNTagList(in, 10)
	t.Log(info)
}
