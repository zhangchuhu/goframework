package handler

import (
	"encoding/json"
	"testing"
	"time"
)

func TestGetHostIncomingRecords(t *testing.T) {
	httprt := HttpRetComm{
		Desc: "success",
		Time: time.Now().Unix(),
		Data: &IncomingRecordS{
			CurMonthIncoming: 100,
			CashIncoming:     20,
			Records: []IncomingRecord{
				IncomingRecord{
					DateTime:        "2018-06-12",
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
