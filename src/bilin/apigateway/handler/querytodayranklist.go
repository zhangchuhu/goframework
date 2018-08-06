package handler

import (
	"bilin/apigateway/services/rank"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"encoding/json"
	"fmt"
	"net/http"
)

func HandleTodayRankList(rw http.ResponseWriter, r *http.Request) *HttpError {

	today_rank := rank.TadayRank{
		ContributeRank: rank.GetContributeRank(),
		AnchorRank:     rank.GetAnchorRank(),
		GuardRank:      rank.GetGuardRank(),
	}

	resp := &HttpRetComm{
		IsEncrypt: "false",
		Data: HttpRetDataComm{
			Code: 0,
			Msg:  "success",
			Body: rank.TadayRankBody{
				IsShow:    true,
				TadayRank: today_rank,
			},
		},
	}

	byte, err := json.Marshal(resp)
	if err != nil {
		appzaplog.Error("failed to marshal", zap.Error(err))
		return jsonMarshalHttpErr
	}

	fmt.Fprintf(rw, "%s", byte)
	return successHttp
}
