package dao

import "testing"

func TestGetLivingRecordByHostID(t *testing.T) {
	info, err := GetLivingRecordByHostID([]uint64{8766212}, "2018-06-13", "2018-06-15")
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", info)
}

func TestGetLivingRecordByHostIDAndRoomID(t *testing.T) {
	info, err := GetLivingRecordByHostIDAndRoomID([]uint64{8766212}, []uint64{1100}, "2018-06-13", "2018-06-15")
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", info)
}
