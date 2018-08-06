package handler

import (
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/httpmetrics"
	"github.com/gin-gonic/gin"
	"time"
)

var (
	successHttp = &HttpError{"success", 200}
	authHttpErr = &HttpError{"auth failed", 1}
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

func auth(c *gin.Context) (bool, error) {
	return true, nil
}
