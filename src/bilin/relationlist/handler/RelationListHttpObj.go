package handler

import (
	"bilin/bcserver/bccommon"
	"bilin/relationlist/controller"
	"bilin/relationlist/entity"
	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/kissgo/httpmetrics"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// retWrite marshal the result and write to client(get).
func retWrite(w http.ResponseWriter, r *http.Request, res map[string]interface{}, start time.Time) {
	data, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		log.Error(r.URL.String(), zap.Any("err", err), zap.Any("res", res))
		return
	}
	dataStr := string(data) + "\n"
	if _, err := w.Write([]byte(dataStr)); err != nil {
		log.Error(r.URL.String(), zap.Any("err", err), zap.Any("data", dataStr))
		return
	}
	log.Debug(r.URL.String(), zap.Any("ip", r.RemoteAddr), zap.Any("time", time.Now().Sub(start).Seconds()))
}

type RelationListHttpObj struct {
}

type RelationListResp struct {
	Code      int32                      `json:"code"`
	ErrorDesc string                     `json:"error"`
	Data      *entity.RelationStatistics `json:"data"`
}

func NewRelationListHttpObj() (o *RelationListHttpObj) {
	return &RelationListHttpObj{}
}

func FakeResp() *RelationListResp {
	RelationListResp := &RelationListResp{Code: 0, Data: &entity.RelationStatistics{}}
	anchorInfo := &entity.UserRelationInfo{UserID: 17795535, Nick: "konakona",
		Avatar: "https://img.inbilin.com/17795556/17795556_1521709188889.jpg-small", RelationVal: 88888}
	RelationListResp.Data.AnchorInfo = anchorInfo

	for i := 0; i < 10; i++ {
		item := &entity.UserRelationInfo{UserID: uint64(17795535 + i), Nick: fmt.Sprintf("guest%d", i),
			Avatar: "https://img.inbilin.com/17795556/17795556_1521709188889.jpg-small", RelationVal: int64(1000 + i*100)}

		if i < 3 {
			item.MedalUrl = "https://vipweb.bs2cdn.yy.com/vipinter_f5b19f9974fd477f9d012240ba869801.png"
			item.MedalText = fmt.Sprintf("勋章%d", i)
		}

		RelationListResp.Data.RelationList = append(RelationListResp.Data.RelationList, item)
	}

	return RelationListResp
}

func GetValueFromReq(form map[string][]string, pramer string) (result int) {
	if form[pramer] == nil {
		result = 0
	} else {
		var err error
		result, err = strconv.Atoi(form[pramer][0])
		if err != nil {
			result = 0
		}
	}

	return
}

func (o *RelationListHttpObj) GetRelationList(resp http.ResponseWriter, req *http.Request) {
	const prefix = "GetRelationList "

	var metrics_ret int64 = 0
	defer func(now time.Time) {
		httpmetrics.DefReport(strings.TrimSpace(prefix), metrics_ret, now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	if req.Method != "GET" {
		http.Error(resp, "Method Not Allowed", http.StatusMethodNotAllowed)
		log.Error(prefix+"Method Not Allowed", zap.Any("req.Method", req.Method))
		metrics_ret = -1
		return
	}

	req.ParseForm()
	if req.Form["uid"] == nil {
		log.Error(prefix + "uid not given in query string")
		metrics_ret = -1
		return
	}

	owner := GetValueFromReq(req.Form, "uid")
	reqType := GetValueFromReq(req.Form, "type")
	start := GetValueFromReq(req.Form, "start")
	rows := GetValueFromReq(req.Form, "rows")
	if rows == 0 {
		rows = 50
	}

	log.Info(prefix+"begin", zap.Any("owner", owner), zap.Any("reqType", reqType), zap.Any("start", start), zap.Any("rows", rows))

	RelationListResp := &RelationListResp{Code: 0}
	if reqType == 0 {
		RelationListResp.Data = controller.GetDailyRelationList(uint64(owner), start, rows)
	} else if reqType == 1 {
		RelationListResp.Data = controller.GetWeeklyRelationList(uint64(owner), start, rows)
	} else {
		RelationListResp.Data = controller.GetTotalRelationList(uint64(owner), start, rows)
	}

	ret, err := json.Marshal(RelationListResp)
	if err != nil {
		log.Error("[-]"+prefix+"json.Marshal failed", zap.Any("err", err))
		metrics_ret = -2
		return
	}
	fmt.Fprintf(resp, string(ret))

	log.Info("[+]"+prefix+"end", zap.Any("RelationListResp", RelationListResp))
}

//http 接口，供客户端获取亲密榜榜单
func (o *RelationListHttpObj) GetRelationListByJsonP(resp http.ResponseWriter, req *http.Request) {
	const prefix = "GetRelationListByJsonP "

	var metrics_ret int64 = 0
	defer func(now time.Time) {
		httpmetrics.DefReport(strings.TrimSpace(prefix), metrics_ret, now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	if req.Method != "GET" {
		http.Error(resp, "Method Not Allowed", http.StatusMethodNotAllowed)
		log.Error(prefix+"Method Not Allowed", zap.Any("req.Method", req.Method))
		metrics_ret = -1
		return
	}

	log.Info("[+]" + prefix + "begin")

	callbackName := req.URL.Query().Get("callback")
	if callbackName == "" {
		log.Error(prefix + "Please give callback name in query string")
		metrics_ret = -1
		return
	}

	req.ParseForm()
	if req.Form["uid"] == nil {
		log.Error(prefix + "uid not given in query string")
		metrics_ret = -1
		return
	}

	owner := GetValueFromReq(req.Form, "uid")
	reqType := GetValueFromReq(req.Form, "type")
	start := GetValueFromReq(req.Form, "start")
	rows := GetValueFromReq(req.Form, "rows")
	if rows == 0 {
		rows = 50
	}

	log.Info(prefix+"begin", zap.Any("owner", owner), zap.Any("reqType", reqType), zap.Any("start", start), zap.Any("rows", rows))

	RelationListResp := &RelationListResp{Code: 0}
	if reqType == 0 {
		RelationListResp.Data = controller.GetDailyRelationList(uint64(owner), start, rows)
	} else if reqType == 1 {
		RelationListResp.Data = controller.GetWeeklyRelationList(uint64(owner), start, rows)
	} else {
		RelationListResp.Data = controller.GetTotalRelationList(uint64(owner), start, rows)
	}

	ret, err := json.Marshal(RelationListResp)
	if err != nil {
		log.Error("[-]"+prefix+"json.Marshal failed", zap.Any("err", err))
		metrics_ret = -2
		return
	}
	resp.Header().Set("Content-Type", "application/javascript")
	fmt.Fprintf(resp, "%s(%s);", callbackName, string(ret))

	log.Info("[+]"+prefix+"end", zap.Any("RelationListResp", RelationListResp))
}
