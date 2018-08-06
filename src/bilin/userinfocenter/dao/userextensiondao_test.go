package dao

import "testing"

func TestGetUserExtension(t *testing.T) {
	info, err := GetUserExtension(100)
	if err != nil {
		t.Errorf("GetUserExtension failed")
	}
	cc, err := GetCommonCountry(info.CITY)
	if err != nil {
		t.Error("GetCommonCountry failed", err)
	}
	t.Log("info:", info, cc.Name)
}
