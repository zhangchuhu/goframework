package handler

import (
	"encoding/json"
	"testing"
	"time"
)

func TestHandleLivingRecordReq(t *testing.T) {
	httpret := HttpRetComm{
		Desc: "success",
		Time: time.Now().Unix(),
		Data: LivingRecords{
			Records: []LivingRecord{
				{
					HOSTBilinID:        100,
					LivingStartTime:    "2015-07-08",
					LivingTime:         50,
					AudienceNum:        100,
					MikeUserNum:        10,
					OneMinuteOutInRate: 70,
					AverageStayTime:    50,
					RoomID:             9100,
					//GiftHeartNum:       100,
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
