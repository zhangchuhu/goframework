package handler

import (
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"github.com/gin-gonic/gin"
	"strconv"
)

type HttpRetComm struct {
	Code int64       `json:"code"`
	Desc string      `json:"desc"`
	Time int64       `json:"time"`
	Data interface{} `json:"data,omitempty"`
}

func UInt64Param(key string, c *gin.Context) (uint64, error) {
	id := c.Param(key)
	appzaplog.Debug("[+]GetHost", zap.String("id", id))
	idInt, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		appzaplog.Error("GetHostContract ParseUint err", zap.Error(err))
		return idInt, err
	}
	return idInt, nil
}
