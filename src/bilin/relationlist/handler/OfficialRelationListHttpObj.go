package handler

import (
	"bilin/bcserver/bccommon"
	"bilin/relationlist/controller"
	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/kissgo/httpmetrics"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type OfficialRelationListHttpObj struct {
}

type OfficialRelationListRequest struct {
	RoomId    uint64   `json:"roomid"`
	OldHost   uint64   `json:"old_host"`
	NewHost   uint64   `json:"new_host"`
	GuestUids []uint64 `json:"guest_uids"`
}

type OfficialRelationListResp struct {
	Code      int32  `json:"code"`
	ErrorDesc string `json:"error"`
}

func NewOfficialRelationListHttpObj() (o *OfficialRelationListHttpObj) {
	return &OfficialRelationListHttpObj{}
}

func (o *OfficialRelationListHttpObj) OfficialRoomChangeOwner(resp http.ResponseWriter, req *http.Request) {
	const prefix = "OfficialRoomChangeOwner "

	if req.Method != "POST" {
		http.Error(resp, "Method Not Allowed", http.StatusMethodNotAllowed)
		log.Error(prefix+"Method Not Allowed", zap.Any("req.Method", req.Method))
		return
	}

	officialResp := &OfficialRelationListResp{Code: 0}
	officialRequest := &OfficialRelationListRequest{}
	var err error

	defer func(now time.Time) {
		httpmetrics.DefReport(strings.TrimSpace(prefix), int64(officialResp.Code), now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	body, _ := ioutil.ReadAll(req.Body)
	if err = json.Unmarshal(body, officialRequest); err != nil {
		log.Error(prefix+"json.Unmarshal failed", zap.Error(err), zap.Any("body", string(body)))
		officialResp.Code = -1
		goto RETURN
	}

	log.Debug("[+]"+prefix+"begin", zap.Any("officialRequest", officialRequest))

	for _, item := range officialRequest.GuestUids {
		//旧的host下麦了，需要停止亲密度计算
		controller.UserOffMike(officialRequest.OldHost, item)
		//新的host上麦，需要开始计算亲密度
		if officialRequest.NewHost != 0 && item != officialRequest.NewHost {
			controller.UserOnMike(officialRequest.NewHost, item, time.Now().Unix())
		}
	}

RETURN:

	ret, _ := json.Marshal(officialResp)
	fmt.Fprintf(resp, string(ret))

	log.Debug("[+]"+prefix+"end", zap.Any("officialRequest", officialRequest), zap.Any("officialResp", officialResp))

}
