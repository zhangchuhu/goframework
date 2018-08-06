package handler

import (
	"bilin/clientcenter"
	"bilin/protocol"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type HostRec struct {
	ID     int64  `json:"id"`
	Hostid uint64 `json:"hostid,omitempty"`
	Typeid int64  `json:"typeid,omitempty"`
}

type HostRecList struct {
	TotalPageSize int       `json:"total_page_size"`
	Items         []HostRec `json:"items"`
}

func GetHostRec(c *gin.Context) {
	appzaplog.Debug("[+]GetHostRec")
	var (
		code    = http.StatusOK
		hostRec = &HostRecList{}
	)

	for {
		pageinfo, err := pageInfo(c)
		if err != nil {
			appzaplog.Error("GetHostRec pageInfo err", zap.Error(err))
			code = http.StatusBadRequest
			break
		}
		info, err := clientcenter.ConfClient().CategoryHostRec(context.TODO(), &bilin.CategoryHostRecReq{})
		if err != nil {
			appzaplog.Error("GetHostRec CategoryGuildRec err", zap.Error(err))
			code = http.StatusInternalServerError
			break
		}
		for _, v := range info.Cateogryinfos {
			hostRec.Items = append(hostRec.Items, HostRec{
				ID:     v.Id,
				Hostid: v.Hostid,
				Typeid: int64(v.Typeid),
			})
		}
		hostRec.TotalPageSize = len(hostRec.Items)
		if len(hostRec.Items) > 0 {
			//hostRec.TotalPageSize = len(hostRec.Items)/pageinfo.pageSize + 1
			start, end, err := pageinfo.StartEnd(len(hostRec.Items))
			if err != nil {
				appzaplog.Error("GetHostRec StartEnd err", zap.Error(err))
				code = http.StatusBadRequest
				break
			}
			hostRec.Items = hostRec.Items[start:end]
		}

		break
	}

	jsonret := &HttpRetComm{
		Code: code,
		Data: hostRec,
	}

	appzaplog.Debug("[-]GetHostRec", zap.Any("resp", jsonret))
	c.JSON(200, jsonret)
}

func CreateHostRec(c *gin.Context) {
	appzaplog.Debug("[+]CreateHostRec")
	code := http.StatusOK
	hostRec := &HostRec{}

	for {
		if err := c.BindJSON(hostRec); err != nil {
			appzaplog.Error("CreateHostRec BindJSON err", zap.Error(err))
			code = http.StatusBadRequest
			break
		}
		appzaplog.Debug("CreateHostRec", zap.Any("req", hostRec))

		_, err := clientcenter.ConfClient().CreateCategoryHostRec(context.TODO(), &bilin.CreateCategoryHostRecReq{
			Info: &bilin.CategoryHostRecInfo{
				Hostid: hostRec.Hostid,
				Typeid: uint64(hostRec.Typeid),
			},
		})
		if err != nil {
			appzaplog.Error("CreateHostRec CreateCategoryHostRec err", zap.Error(err))
			code = http.StatusInternalServerError
			break
		}
		break
	}

	jsonret := &HttpRetComm{
		Code: code,
		Data: hostRec,
	}

	appzaplog.Debug("[-]CreateHostRec", zap.Any("resp", jsonret))
	c.JSON(200, jsonret)
}

func DelHostRec(c *gin.Context) {
	appzaplog.Debug("[+]DelHostRec")
	code := http.StatusOK

	for {
		roomid := c.Param("id")
		if roomid == "" {
			appzaplog.Error("DelHostRec empty roomid err")
			code = http.StatusBadRequest
			break
		}
		appzaplog.Debug("DelHostRec", zap.String("roomid", roomid))
		roomidU64, err := strconv.ParseInt(roomid, 10, 64)
		if err != nil {
			appzaplog.Error("DelHostRec ParseUint err")
			code = http.StatusBadRequest
			break
		}
		_, err = clientcenter.ConfClient().DelCategoryHostRec(context.TODO(), &bilin.DelCategoryHostRecReq{
			Info: &bilin.CategoryHostRecInfo{
				Id: roomidU64,
			},
		})
		if err != nil {
			appzaplog.Error("DelHostRec DelCategoryGuildRec err", zap.Error(err))
			code = http.StatusInternalServerError
			break
		}
		break
	}

	jsonret := &HttpRetComm{
		Code: code,
	}

	appzaplog.Debug("[-]DelHostRec", zap.Any("resp", jsonret))
	c.JSON(200, jsonret)
}

func UpdateHostRec(c *gin.Context) {
	appzaplog.Debug("[+]UpdateHostRec")
	code := http.StatusOK
	guildRec := &HostRec{}

	for {
		if err := c.BindJSON(guildRec); err != nil {
			appzaplog.Error("UpdateHostRec BindJSON err", zap.Error(err))
			code = http.StatusBadRequest
			break
		}
		appzaplog.Debug("UpdateHostRec", zap.Any("req", guildRec))
		_, err := clientcenter.ConfClient().UpdateCategoryHostRec(context.TODO(), &bilin.UpdateCategoryHostRecReq{
			Info: &bilin.CategoryHostRecInfo{
				Id:     guildRec.ID,
				Hostid: guildRec.Hostid,
				Typeid: uint64(guildRec.Typeid),
			},
		})
		if err != nil {
			appzaplog.Error("UpdateHostRec UpdateCategoryGuildRec err", zap.Error(err))
			code = http.StatusInternalServerError
			break
		}
		break
	}

	jsonret := &HttpRetComm{
		Code: code,
		Data: guildRec,
	}

	appzaplog.Debug("[-]UpdateHostRec", zap.Any("resp", jsonret))
	c.JSON(200, jsonret)
}
