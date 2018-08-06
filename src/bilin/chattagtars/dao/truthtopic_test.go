package dao

import (
	"strings"
	"testing"
)

func TestTruthTopic_Create(t *testing.T) {
	truth := TruthTopic{
		Topic: "爱一个女孩子，与其为了她的幸福而放弃她,\n不如留住她，为她的幸福而努力",
	}
	if err := truth.Create(); err != nil {
		t.Error(err)
	}
}

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
