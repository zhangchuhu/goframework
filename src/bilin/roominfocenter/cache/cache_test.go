package cache

import "testing"

func TestGetRoomCache(t *testing.T) {
	cache := GetRoomCache()
	if len(cache) == 0 {
		t.Error("empty cache")
	}
}
