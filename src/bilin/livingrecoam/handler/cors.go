package handler

import (
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"github.com/gin-gonic/gin"
)

func CORSMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", c.GetHeader("Origin"))
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, X-Auth-Token")
		c.Header("Access-Control-Allow-Headers", "Content-Type,Access-Token,X-Token")
	}
}

func CORSOptionHandler(c *gin.Context) {
	appzaplog.Debug("[+]OPTIONS")
	jsonret := &HttpRetComm{
		Code: 200,
	}
	c.JSON(204, jsonret)
	appzaplog.Debug("[-]OPTIONS", zap.Any("resp", jsonret))
}
