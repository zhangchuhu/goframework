package handler

import (
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginBody struct {
	User     string `form:"user" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	appzaplog.Debug("post login", zap.String("username", c.Param("username")))
	var json LoginBody
	var resp HttpRetComm = HttpRetComm{
		Code: http.StatusOK,
	}
	if err := c.BindJSON(&json); err != nil {
		resp.Code = http.StatusBadRequest
		c.JSON(http.StatusBadRequest, &resp)
		return
	}
	if json.User == "admin" && json.Password == "admin" {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"token": "admin"}})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
	}
}

func UserInfo(c *gin.Context) {
	appzaplog.Debug("[+]UserInfo")
	c.JSON(200, HttpRetComm{
		Code: 200,
		Data: &User{
			Roles:  []string{"admin"},
			Name:   "admin",
			Avatar: "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
		},
	})
}

type User struct {
	Roles  []string `json:"roles"`
	Name   string   `json:"name"`
	Avatar string   `json:"avatar"`
}

func LoginOut(c *gin.Context) {
	jsonret := &HttpRetComm{
		Code: 200,
	}
	c.JSON(204, jsonret)
	appzaplog.Debug("[-]LoginOut", zap.Any("resp", jsonret))
}
