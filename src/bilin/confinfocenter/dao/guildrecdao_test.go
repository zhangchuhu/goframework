package dao

import "testing"

func TestGuildRec_Create(t *testing.T) {
	guildrec := &GuildRec{
		RoomID: 100,
		TypeId: 3,
	}
	if err := guildrec.Create(); err != nil {
		t.Error(err)
	}
	t.Log(guildrec)
}

func TestGetGuildRec(t *testing.T) {
	info, err := GetGuildRec()
	if err != nil {
		t.Error(err)
	}
	t.Log(info)
}

func TestDelGuildRec(t *testing.T) {
	err := DelGuildRec(1)
	if err != nil {
		t.Error(err)
	}
	t.Log("delete ok")
}

func TestUpdateGuildRec(t *testing.T) {
	uguild := &GuildRec{
		RoomID: 1,
		TypeId: 3,
	}
	if err := UpdateGuildRec(uguild); err != nil {
		t.Error(err)
	}
	t.Log(uguild)
}
