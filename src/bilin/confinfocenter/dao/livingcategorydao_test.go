package dao

import (
	"testing"
)

func TestGetLivingCategorys(t *testing.T) {
	info, err := GetLivingCategorys()
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", info)
}
