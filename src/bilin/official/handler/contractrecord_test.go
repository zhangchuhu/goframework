package handler

import (
	"encoding/json"
	"testing"
	"time"
)

func TestHandleContractRecordReq(t *testing.T) {
	httpret := HttpRetComm{
		Desc: "success",
		Time: time.Now().Unix(),
		Data: &ContractRecords{
			Records: []ContractRecord{
				{
					BILINNumber:          100,
					NickName:             "TEST",
					ContractEndTime:      1000,
					ContractStartTime:    100,
					GuildSharePercentage: 70,
					HostSharePercentage:  30,
				},
			},
		},
	}

	bin, err := json.Marshal(httpret)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(bin))
}
