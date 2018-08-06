package dao

import (
	"testing"
	"time"
)

func TestGetStickie(t *testing.T) {
	info, err := GetStickie()
	if err != nil {
		t.Errorf("GetStickie failed", err)
	}
	t.Logf("%v", info)
}

func TestCreateStickie(t *testing.T) {
	startTime, err := time.ParseInLocation(TimeLayoutOthers, "2018-06-15 18:20:00", time.Local)
	if err != nil {
		t.Error(err)
		return
	}
	endtime, err := time.ParseInLocation(TimeLayoutOthers, "2018-06-24 18:20:20", time.Local)
	if err != nil {
		t.Error(err)
		return
	}
	stickie := &Stickie{
		TypeId:    1000,
		RoomID:    32953,
		Weight:    10,
		StartTime: startTime,
		EndTime:   endtime,
	}
	if err := stickie.Create(); err != nil {
		t.Error(err)
	}
	t.Log(*stickie)
}

const TimeLayoutOthers string = "2006-01-02 15:04:05"

func TestStickie_Update(t *testing.T) {
	updatedids := []uint{1}
	//otheruids := []uint{2,3,4,5,6}
	startTime, err := time.ParseInLocation(TimeLayoutOthers, "2018-06-17 18:20:00", time.Local)
	if err != nil {
		t.Error(err)
		return
	}
	endtime, err := time.ParseInLocation(TimeLayoutOthers, "2018-06-18 18:20:20", time.Local)
	if err != nil {
		t.Error(err)
		return
	}

	sp := &Stickie{
		StartTime: startTime,
		EndTime:   endtime,
	}
	for _, v := range updatedids {
		sp.ID = v
		err = sp.Update()
		if err != nil {
			t.Error(err)
		}
	}

	t.Log(sp)
}

func TestStickie_Del(t *testing.T) {
	sticky := Stickie{}
	sticky.ID = 1
	if err := sticky.Del(); err != nil {
		t.Error(err)
	}
}
