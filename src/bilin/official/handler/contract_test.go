package handler

import (
	"encoding/json"
	"testing"
	"time"
)

func TestHandleContractReq(t *testing.T) {
	httpret := HttpRetComm{
		Desc: "success",
		Time: time.Now().Unix(),
		Data: &Contract{
			GuildID:              1100,
			HostUid:              200,
			ContractStartTime:    time.Now().Unix(),
			ContractEndTime:      time.Now().Add(time.Minute * 60 * 24).Unix(),
			GuildSharePercentage: 70,
			HostSharePercentage:  30,
		},
	}

	bin, err := json.Marshal(httpret)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(bin))
}
