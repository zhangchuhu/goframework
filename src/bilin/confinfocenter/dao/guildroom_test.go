package dao

import "testing"

func TestGuildRoom_Create(t *testing.T) {
	room := &GuildRoom{
		GuildID: 100,
		RoomID:  9018,
	}
	if err := room.Create(); err != nil {
		t.Error(err)
	}
	t.Log(room)
}

func TestGetGuildChannelS(t *testing.T) {
	info, err := GetGuildRoomS()
	if err != nil {
		t.Error(err)
	}
	t.Log(info)
}

func TestGuildRoom_Del(t *testing.T) {
	room := &GuildRoom{
		GuildID: 100,
		RoomID:  9018,
	}
	if err := room.Del(); err != nil {
		t.Error(err)
	}
	t.Log(room)
}

func TestGuildRoom_Get(t *testing.T) {
	room := &GuildRoom{}
	info, err := room.Get()
	if err != nil {
		t.Error(err)
	}
	t.Log(info)
}
