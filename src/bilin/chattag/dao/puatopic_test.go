package dao

import (
	"strings"
	"testing"
)

func TestGetTopicNotIn(t *testing.T) {
	info, err := GetTopicNotIn([]int64{2}, 1)
	if err != nil {
		t.Error(err)
	}
	t.Log(info)
	for _, v := range info {
		lines := strings.Split(v.Topic, "\n")
		t.Log(lines[1])
	}
}

func TestRandPuaTopic(t *testing.T) {
	topic, err := RandPuaTopic(1)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", topic)
}
