package dao

import (
	"strings"
	"testing"
)

func TestGetTruthTopicNotIn(t *testing.T) {
	info, err := GetTruthTopicNotIn([]int64{2}, 1)
	if err != nil {
		t.Error(err)
	}
	t.Log(info)
	for _, v := range info {
		lines := strings.Split(v.Topic, "\n")
		t.Log(lines[1])
	}
}
