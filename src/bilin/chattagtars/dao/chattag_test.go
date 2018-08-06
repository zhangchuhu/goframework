package dao

import "testing"

func TestChatTag_Create(t *testing.T) {
	tags := []ChatTag{
		{TagName: "小鲜肉", TagColor: "#0000FF"},
		{TagName: "大叔", TagColor: "#0000FF"},
		{TagName: "闷骚", TagColor: "#0000FF"},
		{TagName: "搞笑", TagColor: "#0000FF"},
		{TagName: "土豪", TagColor: "#0000FF"},
		{TagName: "声优", TagColor: "#0000FF"},
		{TagName: "撩妹大神", TagColor: "#0000FF"},
		{TagName: "聊骚专家", TagColor: "#0000FF"},
	}
	for _, tag := range tags {
		if err := tag.Create(); err != nil {
			t.Error(err)
		}
	}
}

func TestGetAll(t *testing.T) {
	info, err := GetAll()
	if err != nil {
		t.Error(err)
	}
	t.Log(info)
}

func TestChatTag_Update(t *testing.T) {
	from := &ChatTag{
		TagColor: "#0000FF",
	}
	from.ID = 1
	if err := from.Update(); err != nil {
		t.Error(err)
	}
}
