package dao

import "testing"

func TestGetUserBiID(t *testing.T) {
	if info, err := GetUserBiID(100); err != nil {
		t.Error(err)
	} else {
		t.Log(info)
	}
}

func TestBatchUserBiID(t *testing.T) {
	if info, err := BatchUserBiID([]uint64{}); err != nil {
		t.Error(err)
	} else {
		t.Log(info)
	}
}
