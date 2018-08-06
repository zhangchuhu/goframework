package dao

import (
	"testing"
)

func TestGetUserBadges(t *testing.T) {
	info, err := GetUserBadges()
	if err != nil {
		t.Error(err)
	}
	t.Log(info)
}
