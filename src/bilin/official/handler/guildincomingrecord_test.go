package handler

import (
	"encoding/json"
	"testing"
	"time"
)

func TestGetGuildIncomingRecords(t *testing.T) {
	httprt := HttpRetComm{
		Desc: "success",
		Time: time.Now().Unix(),
		Data: &GuildIncomingRecordS{
			CurMonthIncoming: 100,
			CashIncoming:     20,
			Records: []GuildIncomingRecord{
				GuildIncomingRecord{
					HostBilinID:     100,
					HostNickName:    "ttt",
					HeartNum:        100,
					GuildPercentage: 70,
					HOSTIncoming:    20,
					GuildIncoming:   10,
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
