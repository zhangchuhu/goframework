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

type GuildRec struct {
	ID     uint64 `json:"id,omitempty"`
	Roomid uint64 `json:"roomid,omitempty"`
	Typeid int64  `json:"typeid,omitempty"`
}

type GuildRecList struct {
	TotalPageSize int        `json:"total_page_size"`
	Items         []GuildRec `json:"items"`
}

func GetGuildRec(c *gin.Context) {
	appzaplog.Debug("[+]GetGuildRec")
	code := http.StatusOK
	guildRec := &GuildRecList{}

	for {
		pageinfo, err := pageInfo(c)
		if err != nil {
			appzaplog.Error("GetGuildRec pageInfo err", zap.Error(err))
			code = http.StatusBadRequest
			break
		}
		info, err := clientcenter.ConfClient().CategoryGuildRec(context.TODO(), &bilin.CategoryGuildRecReq{})
		if err != nil {
			appzaplog.Error("CategoryGuildRec err", zap.Error(err))
			code = http.StatusInternalServerError
			break
		}
		guildRec.TotalPageSize = len(info.Cateogryguildinfos)
		for _, v := range info.Cateogryguildinfos {
			guildRec.Items = append(guildRec.Items, GuildRec{
				Roomid: v.Roomid,
				Typeid: int64(v.Typeid),
				ID:     v.Id,
			})
		}

		if len(guildRec.Items) > 0 {
			start, end, err := pageinfo.StartEnd(len(guildRec.Items))
			if err != nil {
				appzaplog.Error("GetHostRec StartEnd err", zap.Error(err))
				code = http.StatusBadRequest
				break
			}
			guildRec.Items = guildRec.Items[start:end]
		}
		break
	}

	jsonret := &HttpRetComm{
		Code: code,
		Data: guildRec,
	}

	appzaplog.Debug("[-]GetGuildRec", zap.Any("resp", jsonret))
	c.JSON(200, jsonret)
}

func UpdateGuildRec(c *gin.Context) {
	appzaplog.Debug("[+]UpdateGuildRec")
	code := http.StatusOK
	guildRec := &GuildRec{}

	for {
		if err := c.BindJSON(guildRec); err != nil {
			appzaplog.Error("UpdateGuildRec BindJSON err", zap.Error(err))
			code = http.StatusBadRequest
			break
		}
		appzaplog.Debug("UpdateGuildRec", zap.Any("req", guildRec))
		_, err := clientcenter.ConfClient().UpdateCategoryGuildRec(context.TODO(), &bilin.UpdateCategoryGuildRecReq{
			Info: &bilin.CategoryGuildRecInfo{
				Id:     guildRec.ID,
				Roomid: guildRec.Roomid,
				Typeid: uint64(guildRec.Typeid),
			},
		})
		if err != nil {
			appzaplog.Error("UpdateGuildRec UpdateCategoryGuildRec err", zap.Error(err))
			code = http.StatusInternalServerError
			break
		}
		break
	}

	jsonret := &HttpRetComm{
		Code: code,
		Data: guildRec,
	}

	appzaplog.Debug("[-]UpdateGuildRec", zap.Any("resp", jsonret))
	c.JSON(200, jsonret)
}

func DelGuildRec(c *gin.Context) {
	appzaplog.Debug("[+]DelGuildRec")
	code := http.StatusOK

	for {
		roomid := c.Param("id")
		if roomid == "" {
			appzaplog.Error("empty roomid err")
			code = http.StatusBadRequest
			break
		}
		appzaplog.Debug("DelGuildRec", zap.String("roomid", roomid))
		roomidU64, err := strconv.ParseUint(roomid, 10, 64)
		if err != nil {
			appzaplog.Error("DelGuildRec ParseUint err")
			code = http.StatusBadRequest
			break
		}
		_, err = clientcenter.ConfClient().DelCategoryGuildRec(context.TODO(), &bilin.DelCategoryGuildRecReq{
			Info: &bilin.CategoryGuildRecInfo{
				Id: roomidU64,
			},
		})
		if err != nil {
			appzaplog.Error("DelGuildRec DelCategoryGuildRec err", zap.Error(err))
			code = http.StatusInternalServerError
			break
		}
		break
	}

	jsonret := &HttpRetComm{
		Code: code,
	}

	appzaplog.Debug("[-]DelGuildRec", zap.Any("resp", jsonret))
	c.JSON(200, jsonret)
}

func CreateGuildRec(c *gin.Context) {
	appzaplog.Debug("[+]CreateGuildRec")
	code := http.StatusOK
	guildRec := &GuildRec{}

	for {
		if err := c.BindJSON(guildRec); err != nil {
			appzaplog.Error("CreateGuildRec BindJSON err", zap.Error(err))
			code = http.StatusBadRequest
			break
		}
		appzaplog.Debug("CreateGuildRec", zap.Any("req", guildRec))

		_, err := clientcenter.ConfClient().CreateCategoryGuildRec(context.TODO(), &bilin.CreateCategoryGuildRecReq{
			Info: &bilin.CategoryGuildRecInfo{
				Roomid: guildRec.Roomid,
				Typeid: uint64(guildRec.Typeid),
			},
		})
		if err != nil {
			appzaplog.Error("CreateGuildRec CreateCategoryGuildRec err", zap.Error(err))
			code = http.StatusInternalServerError
			break
		}
		break
	}

	jsonret := &HttpRetComm{
		Code: code,
		Data: guildRec,
	}

	appzaplog.Debug("[-]CreateGuildRec", zap.Any("resp", jsonret))
	c.JSON(200, jsonret)
}
