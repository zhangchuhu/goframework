package handler

import (
	"bilin/clientcenter"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/kissgo/httpmetrics"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

type HttpRetDataComm struct {
	Code int64       `json:"code"`
	Msg  string      `json:"msg"`
	Body interface{} `json:"body"`
}

type HttpRetComm struct {
	IsEncrypt string          `json:"isEncrypt"`
	Data      HttpRetDataComm `json:"data"`
}

var (
	successHttp             = &HttpError{"success", 0}
	authHttpErr             = &HttpError{"auth failed", 1}
	ratelimitHttpErr        = &HttpError{"rate limited", 2}
	parseFormHttpErr        = &HttpError{"parseform failed", 3}
	daoGetHttpErr           = &HttpError{"dao get failed", 4}
	delUsedTopicHttpErr     = &HttpError{"del used topic failed", 5}
	notSupportMethodHttpErr = &HttpError{"not support method", 6}
	parseURLHttpErr         = &HttpError{"parse url failed", 7}
	daoPutHttpErr           = &HttpError{"dao put failed", 8}
	TagNotAllowedErr        = &HttpError{desc: "tag not allowed", code: 9}
	daoUpdateHttpErr        = &HttpError{"dao update failed", 10}
	getUserInfoHttpErr      = &HttpError{"getuserinfo failed", 11}
)

type HandlerFuncWithErro func(*gin.Context) *HttpError

type HttpError struct {
	desc string
	code int64
}

func AuthMiddleWare(fn HandlerFuncWithErro) HandlerFuncWithErro {
	return func(c *gin.Context) *HttpError {
		//auth
		if ok, err := auth(c); !ok || err != nil {
			appzaplog.Warn("user auth failed")
			c.JSON(401, gin.H{
				"code": 1,
				"desc": "auth failed",
			})
			return authHttpErr
		}
		return fn(c)
	}
}

func MetricsMiddleWare(uri string, fn HandlerFuncWithErro) func(c *gin.Context) {
	return func(c *gin.Context) {
		var err *HttpError
		defer func(now time.Time) {
			httpmetrics.DefReport(uri, err.code, now, httpmetrics.DefaultSuccessFun)
		}(time.Now())
		err = fn(c)
	}
}

const (
	uidContextKey = "uid"
)

func auth(c *gin.Context) (bool, error) {
	uid := c.Query("userId")
	if uid == "" {
		appzaplog.Warn("auth no userId")
		return false, fmt.Errorf("no userId")
	}
	uidInt, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		appzaplog.Error("auth ParseInt err", zap.Error(err), zap.String("uid", uid))
		return false, err
	}
	c.Set(uidContextKey, uidInt)

	token := c.Query("accessToken")
	if token == "" {
		appzaplog.Warn("auth no accessToken", zap.String("uid", uid))
		return false, fmt.Errorf("auth no accessToken")
	}
	return clientcenter.VerifyAccessToken(token, uid)
}

type QueryStringParam struct {
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

func (rp *QueryStringParam) Unmarshal(c *gin.Context) error {
	var err error
	rp.Platform = c.Query("platform")
	rp.Version = c.Query("version")
	rp.DeviceId = c.Query("deviceId")
	rp.ClientType = c.Query("clientType")
	if rp.UserId, err = strconv.ParseInt(c.Query("userId"), 10, 64); err != nil {
		return err
	}
	return err
}
