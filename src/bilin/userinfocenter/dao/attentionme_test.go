package dao

import "testing"

func TestAttentionMeNum(t *testing.T) {
	count, err := AttentionMeNum(29481968)
	if err != nil {
		t.Error(err)
	}
	t.Log(count)
}
