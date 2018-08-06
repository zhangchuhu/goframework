/*
 * Copyright (c) 2018-07-20.
 * Author: kordenlu
 * 功能描述:${<VARIABLE_NAME>}
 */

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

type Guild struct {
	ID        int64  `json:"id"`
	OW        int64  `json:"ow"`
	Title     string `json:"title"`
	Mobile    string `json:"mobile"`
	Describle string `json:"describle"`
	GuildLogo string `json:"guild_logo"`
}

type GuildList struct {
	TotalPageSize int     `json:"total_page_size"`
	Items         []Guild `json:"items"`
}

func GetGuild(c *gin.Context) {
	appzaplog.Debug("[+]GetGuild")
	code := http.StatusOK
	guildRec := &GuildList{}

	for {
		pageinfo, err := pageInfo(c)
		if err != nil {
			appzaplog.Error("GetGuildRoom pageInfo err", zap.Error(err))
			code = http.StatusBadRequest
			break
		}
		info, err := clientcenter.GuildTarsClient().RGuild(context.TODO(), &bilin.RGuildReq{})
		if err != nil {
			appzaplog.Error("GetGuild GuildRoomS err", zap.Error(err))
			code = http.StatusInternalServerError
			break
		}
		guildRec.TotalPageSize = len(info.Info)
		for _, v := range info.Info {
			guildRec.Items = append(guildRec.Items, Guild{
				ID:        v.Id,
				OW:        v.Ow,
				Title:     v.Title,
				Mobile:    v.Mobile,
				Describle: v.Describle,
				GuildLogo: v.Guildlog,
			})
		}

		if len(guildRec.Items) > 0 {
			start, end, err := pageinfo.StartEnd(len(guildRec.Items))
			if err != nil {
				appzaplog.Error("GetGuild StartEnd err", zap.Error(err))
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

	appzaplog.Debug("[-]GetGuild", zap.Any("resp", jsonret))
	c.JSON(200, jsonret)
}

func CreateGuild(c *gin.Context) {
	appzaplog.Debug("[+]CreateGuild")
	code := http.StatusOK
	guildRec := &Guild{}

	for {
		if err := c.BindJSON(guildRec); err != nil {
			appzaplog.Error("CreateGuild BindJSON err", zap.Error(err))
			code = http.StatusBadRequest
			break
		}
		appzaplog.Debug("CreateGuild", zap.Any("req", guildRec))

		_, err := clientcenter.GuildTarsClient().CGuild(context.TODO(), &bilin.CGuildReq{
			Info: &bilin.Guild{
				Id:        guildRec.ID,
				Ow:        guildRec.OW,
				Title:     guildRec.Title,
				Mobile:    guildRec.Mobile,
				Describle: guildRec.Describle,
				Guildlog:  guildRec.GuildLogo,
			},
		})
		if err != nil {
			appzaplog.Error("CreateGuild CreateGuildRoom err", zap.Error(err))
			code = http.StatusInternalServerError
			break
		}
		break
	}

	jsonret := &HttpRetComm{
		Code: code,
		Data: guildRec,
	}

	appzaplog.Debug("[-]CreateGuild", zap.Any("resp", jsonret))
	c.JSON(200, jsonret)
}

func UpdateGuild(c *gin.Context) {
	appzaplog.Debug("[+]UpdateGuild")
	code := http.StatusOK
	guildRec := &Guild{}

	for {
		if err := c.BindJSON(guildRec); err != nil {
			appzaplog.Error("UpdateGuild BindJSON err", zap.Error(err))
			code = http.StatusBadRequest
			break
		}
		appzaplog.Debug("UpdateGuild", zap.Any("req", guildRec))

		_, err := clientcenter.GuildTarsClient().UGuild(context.TODO(), &bilin.UGuildReq{
			Info: &bilin.Guild{
				Id:        guildRec.ID,
				Ow:        guildRec.OW,
				Title:     guildRec.Title,
				Mobile:    guildRec.Mobile,
				Describle: guildRec.Describle,
				Guildlog:  guildRec.GuildLogo,
			},
		})
		if err != nil {
			appzaplog.Error("UpdateGuild BindJSON err", zap.Error(err))
			code = http.StatusInternalServerError
			break
		}
		break
	}

	jsonret := &HttpRetComm{
		Code: code,
		Data: guildRec,
	}

	appzaplog.Debug("[-]UpdateGuild", zap.Any("resp", jsonret))
	c.JSON(200, jsonret)
}

func DelGuild(c *gin.Context) {
	appzaplog.Debug("[+]DelGuild")
	code := http.StatusOK

	for {
		roomid := c.Param("id")
		if roomid == "" {
			appzaplog.Error("DelGuild empty roomid err")
			code = http.StatusBadRequest
			break
		}
		appzaplog.Debug("DelGuild", zap.String("roomid", roomid))
		roomidU64, err := strconv.ParseInt(roomid, 10, 64)
		if err != nil {
			appzaplog.Error("DelGuild ParseUint err")
			code = http.StatusBadRequest
			break
		}
		_, err = clientcenter.GuildTarsClient().DGuild(context.TODO(), &bilin.DGuildReq{
			Info: &bilin.Guild{
				Id: roomidU64,
			},
		})
		if err != nil {
			appzaplog.Error("DelGuild DelCategoryGuildRec err", zap.Error(err))
			code = http.StatusInternalServerError
			break
		}
		break
	}

	jsonret := &HttpRetComm{
		Code: code,
	}

	appzaplog.Debug("[-]DelGuild", zap.Any("resp", jsonret))
	c.JSON(200, jsonret)
}
