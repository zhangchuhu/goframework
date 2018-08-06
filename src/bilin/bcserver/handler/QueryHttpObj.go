package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
)

const (
	OK          = 0
	InternalErr = -1
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

// retPWrite marshal the result and write to client(post).
func retPWrite(w http.ResponseWriter, r *http.Request, res map[string]interface{}, body *string, start time.Time) {
	data, err := json.Marshal(res)
	if err != nil {
		log.Error(r.URL.String(), zap.Any("err", err), zap.Any("res", res))
		return
	}
	dataStr := string(data)
	if _, err := w.Write([]byte(dataStr)); err != nil {
		log.Error(r.URL.String(), zap.Any("err", err), zap.Any("data", dataStr))
		return
	}
	log.Debug(r.URL.String(), zap.Any("ip", r.RemoteAddr), zap.Any("time", time.Now().Sub(start).Seconds()),
		zap.Any("res", dataStr), zap.Any("post", *body))
}

type QueryHttpObj struct {
}

func NewQueryHttpObj() (o *QueryHttpObj) {
	o = &QueryHttpObj{}
	return
}

func (o *QueryHttpObj) QueryRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var (
		res = make(map[string]interface{})
	)
	defer retWrite(w, r, res, time.Now())
	res["QueryRoom"] = "success"
}

type OfficialMikeReq struct {
	UserID int64  `json:"uid,omitempty"`
	RoomID int64  `json:"roomid,omitempty"`
	Data   string `json:"data,omitempty"`
}

func (o *QueryHttpObj) OfficialOnMike(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var (
		err       error
		req       OfficialMikeReq
		bodyBytes []byte
		body      string
		res       = make(map[string]interface{})
	)
	defer retPWrite(w, r, res, &body, time.Now())
	if bodyBytes, err = ioutil.ReadAll(r.Body); err != nil {
		goto failed
	}
	body = string(bodyBytes)
	if err = json.Unmarshal(bodyBytes, &req); err != nil {
		goto failed
	}

	log.Debug(r.URL.String(), zap.Any("req", req))
	// 排班主播迟到，而旧主播还没有下麦，有可能没有清麦上旧主播信息
	// 判断用户是否在房间
	// 用户不在房间，回调java
	// 主播相同，不用切换，直接返回成功，广播麦上信息
	// 主播不相同，开始切换

failed:
	if err != nil {
		log.Error(r.URL.String(), zap.Any("err", err))
		res["ret"] = InternalErr
	} else {
		res["ret"] = OK
	}
}

func (o *QueryHttpObj) OfficialOffMike(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var (
		err       error
		req       OfficialMikeReq
		bodyBytes []byte
		body      string
		res       = make(map[string]interface{})
	)
	defer retPWrite(w, r, res, &body, time.Now())
	if bodyBytes, err = ioutil.ReadAll(r.Body); err != nil {
		goto failed
	}
	body = string(bodyBytes)
	if err = json.Unmarshal(bodyBytes, &req); err != nil {
		goto failed
	}

	log.Debug(r.URL.String(), zap.Any("req", req))
	// 1.强制下麦, 清频道数据
	// 2.send shutdown mike msg
	// 3.删除麦上用户
	// 4.回调java

failed:
	if err != nil {
		log.Error(r.URL.String(), zap.Any("err", err))
		res["ret"] = InternalErr
	} else {
		res["ret"] = OK
	}
}
