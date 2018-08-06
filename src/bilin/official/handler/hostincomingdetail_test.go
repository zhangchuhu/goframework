package handler

import (
	"encoding/json"
	"testing"
	"time"
)

func TestGetHostIncomingDetail(t *testing.T) {
	httprt := HttpRetComm{
		Desc: "success",
		Time: time.Now().Unix(),
		Data: &HostIncomingDetailS{
			TotalPagesize: 10,
			Records: []HostIncomingDetail{
				HostIncomingDetail{
					DateTime:              "2018-06-12 23:00",
					ContributeNickName:    "123",
					ContributeBilinNumber: 100,
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
