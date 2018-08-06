package dao

import (
	"testing"
	"time"
)

func TestSyncLivingRoomInfos(t *testing.T) {
	info, err := SyncLivingRoomInfos()
	if err != nil {
		t.Error("SyncLivingRoomInfos failed", err)
	}
	t.Logf("%v", info)
}

//func init() {
//	RedisClient = redis.NewClient(&redis.Options{
//		Addr:     "183.36.122.50:4019",
//		Password: "", // no password set
//		DB:       0,  // use default DB
//	})
//}
func TestSyncLivingRoomInfosByScan(t *testing.T) {
	defer func(now time.Time) {
		t.Log("time spend:", time.Since(now))
	}(time.Now())
	info, err := SyncLivingRoomInfosByScan()
	if err != nil {
		t.Error(err)
	}
	t.Logf("len:%v,%v", len(info), info)
}
