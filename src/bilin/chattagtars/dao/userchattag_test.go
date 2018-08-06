package dao

import "testing"

func TestUserChatTag_Create(t *testing.T) {
	utag := &UserChatTag{
		FromUserID:  200,
		ToUserID:    64,
		ChatTags:    "1,2,3",
		UpdateTimes: 1,
	}
	if err := utag.Create(); err != nil {
		t.Error(err)
	}
}

func TestUserChatTag_Get(t *testing.T) {
	utag := &UserChatTag{
		FromUserID: 200,
		ToUserID:   1,
	}
	if info, err := utag.Get(); err != nil {
		t.Error(err)
	} else {
		t.Logf("%+v", info)
	}
}

func TestUserChatTag_Update(t *testing.T) {

	to := &UserChatTag{
		FromUserID:  100,
		ToUserID:    64,
		UpdateTimes: 2,
	}
	to.ID = 1
	if err := to.Update(); err != nil {
		t.Error(err)
	}
}

func TestUserChatTag_GetAll(t *testing.T) {
	utag := &UserChatTag{
		ToUserID:   17795724,
		FromUserID: 100,
	}
	if info, err := utag.GetAll(); err != nil {
		t.Error(err)
	} else {
		t.Logf("%+v", info)
	}
}

func TestBatchUserChatTag(t *testing.T) {
	info, err := BatchUserChatTag([]int64{17795724, 17796342})
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", info)
}
