package handler

import (
	"encoding/json"
	"testing"
	"time"
)

func TestGetGuildIncomingDetail(t *testing.T) {
	httprt := HttpRetComm{
		Desc: "success",
		Time: time.Now().Unix(),
		Data: &GuildIncomingDetailS{
			TotalPagesize: 10,
			Records: []GuildIncomingDetail{
				{
					DateTime:              "2018-06-12 23:00",
					ContributeNickName:    "123",
					ContributeBilinNumber: 100,
					HostNickName:          "主播昵称",
					HostBilinNumber:       111,
					PropName:              "不知道",
					PropNum:               10,
					TotalPropValue:        100,
				},
			},
		},
	}
	bin, err := json.Marshal(httprt)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(bin))
}
