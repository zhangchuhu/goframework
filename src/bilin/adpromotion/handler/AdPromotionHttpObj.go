package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"bilin/adpromotion/service"
	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"strings"
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

type AdPromotionHttpObj struct {
}

func NewAdPromotionHttpObj() (o *AdPromotionHttpObj) {
	service.MysqlInit()
	go service.ConsumerKafka()
	o = &AdPromotionHttpObj{}
	return
}

func (o *AdPromotionHttpObj) Hello(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var (
		res = make(map[string]interface{})
	)
	defer retWrite(w, r, res, time.Now())
	res["Hello"] = "success"
}

// http://bilin.adpromotion.yy.com/v1/qihu360/click?UniqueID=__UniqueID__&clicktime=__clicktime__&IP=__IP__&OS=__OS__&devicetype=__devicetype__&imei_md5=__imei_md5__&IDFA=__IDFA__&MAC_MD5=__MAC_MD5__&callback_url=__callback_url__
func (o *AdPromotionHttpObj) ClickAd(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	client_ip := strings.Split(r.RemoteAddr, ":")[0]
	log.Debug(r.URL.String(), zap.Any("ip", client_ip), zap.Any("Query", r.Form))

	//存db
	service.MysqlStorageQihu360Clicks(r.Form, client_ip)
}
