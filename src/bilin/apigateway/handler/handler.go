package handler

import (
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/httpmetrics"
	"github.com/juju/ratelimit"
	"net/http"
	"strconv"
	"time"
)

var (
	successHttp             = &HttpError{"success", 0}
	authHttpErr             = &HttpError{"auth failed", 1}
	ratelimitHttpErr        = &HttpError{"rate limited", 2}
	parseFormHttpErr        = &HttpError{"parseform failed", 3}
	noCacheHttpErr          = &HttpError{"no cache found", 4}
	jsonMarshalHttpErr      = &HttpError{"json marshal failed", 5}
	notSupportMethodHttpErr = &HttpError{"not support method", 6}
)

type HttpRetDataComm struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Body interface{} `json:"body"`
}

type HttpRetComm struct {
	IsEncrypt string          `json:"isEncrypt"`
	Data      HttpRetDataComm `json:"data"`
}

type HandlerFuncWithErro func(w http.ResponseWriter, r *http.Request) *HttpError

type HttpError struct {
	desc string
	code int64
}

func AuthMiddleWare(fn http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		//auth
		if ok, err := auth(r); !ok || err != nil {
			appzaplog.Warn("user auth failed")
			http.Error(rw, "user auth failed", http.StatusUnauthorized)
			return
		}
		fn(rw, r)
	}
}

func AuthMiddleWareV2(fn HandlerFuncWithErro) HandlerFuncWithErro {
	return func(rw http.ResponseWriter, r *http.Request) *HttpError {
		//auth
		if ok, err := auth(r); !ok || err != nil {
			appzaplog.Warn("user auth failed")
			http.Error(rw, "user auth failed", http.StatusUnauthorized)
			return authHttpErr
		}
		return fn(rw, r)
	}
}

func RateLimiteMiddleWare(rl *ratelimit.Bucket, fn http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if rl != nil && rl.TakeAvailable(1) == 0 {
			appzaplog.Warn("rate limited triggered")
			http.Error(rw, "rate limited", http.StatusTooManyRequests)
			return
		}
		fn(rw, r)
	}
}

func RateLimiteMiddleWareV2(rl *ratelimit.Bucket, fn HandlerFuncWithErro) HandlerFuncWithErro {
	return func(rw http.ResponseWriter, r *http.Request) *HttpError {
		if rl != nil && rl.TakeAvailable(1) == 0 {
			appzaplog.Warn("rate limited triggered")
			http.Error(rw, "rate limited", http.StatusTooManyRequests)
			return ratelimitHttpErr
		}
		return fn(rw, r)
	}
}

func MetricsMiddleWare(uri string, fn HandlerFuncWithErro) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var err *HttpError
		defer func(now time.Time) {
			httpmetrics.DefReport(uri, err.code, now, httpmetrics.DefaultSuccessFun)
		}(time.Now())
		err = fn(rw, r)
	}
}

func auth(r *http.Request) (bool, error) {
	return true, nil
}

type ReqParam struct {
	Platform     string
	DeviceId     string
	ClientType   string
	Keytimestamp int64
	HiidoId      string
	UserId       int64
	Signature    string
	Version      string
	AccessToken  string
	Ctimestamp   int64
	Cnonce       string
	NetType      string
}

func (rp *ReqParam) ParseURL(r *http.Request) error {
	var err error
	params := r.URL.Query()
	rp.Platform = params.Get("platform")
	rp.Version = params.Get("version")
	rp.DeviceId = params.Get("deviceId")
	rp.ClientType = params.Get("clientType")
	if rp.UserId, err = strconv.ParseInt(params.Get("userId"), 10, 64); err != nil {
		return err
	}
	return err
}
