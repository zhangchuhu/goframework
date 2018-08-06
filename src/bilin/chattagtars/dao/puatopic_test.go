package dao

import (
	"strings"
	"testing"
)

func TestPuaTopic_Create(t *testing.T) {
	pua := PuaTopic{
		Topic: "爱一个女孩子，与其为了她的幸福而放弃她,\n不如留住她，为她的幸福而努力",
	}
	if err := pua.Create(); err != nil {
		t.Error(err)
	}
}

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
