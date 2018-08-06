/*
 * Copyright (c) 2018-07-05.
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

type ChatTag struct {
	ID       int64  `json:"id"`
	TagName  string `json:"tag_name"`
	TagColor string `json:"tag_color"`
}

type ChatTagList struct {
	TotalPageSize int       `json:"total_page_size"`
	Items         []ChatTag `json:"items"`
}

func GetChatTag(c *gin.Context) {
	appzaplog.Debug("[+]GetChatTag")
	code := http.StatusOK
	guildRec := &ChatTagList{}

	for {
		pageinfo, err := pageInfo(c)
		if err != nil {
			appzaplog.Error("GetChatTag pageInfo err", zap.Error(err))
			code = http.StatusBadRequest
			break
		}
		info, err := clientcenter.ChatTagClient().RChatTag(context.TODO(), &bilin.RChatTagReq{})
		if err != nil {
			appzaplog.Error("GetChatTag GuildRoomS err", zap.Error(err))
			code = http.StatusInternalServerError
			break
		}
		guildRec.TotalPageSize = len(info.Chattag)
		for _, v := range info.Chattag {
			guildRec.Items = append(guildRec.Items, ChatTag{
				ID:       v.Id,
				TagName:  v.TagName,
				TagColor: v.TagColor,
			})
		}

		if len(guildRec.Items) > 0 {
			start, end, err := pageinfo.StartEnd(len(guildRec.Items))
			if err != nil {
				appzaplog.Error("GetChatTag StartEnd err", zap.Error(err))
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

	appzaplog.Debug("[-]GetChatTag", zap.Any("resp", jsonret))
	c.JSON(200, jsonret)
}

func CreateChatTag(c *gin.Context) {
	appzaplog.Debug("[+]CreateChatTag")
	code := http.StatusOK
	guildRec := &ChatTag{}

	for {
		if err := c.BindJSON(guildRec); err != nil {
			appzaplog.Error("CreateChatTag BindJSON err", zap.Error(err))
			code = http.StatusBadRequest
			break
		}
		appzaplog.Debug("CreateChatTag", zap.Any("req", guildRec))

		_, err := clientcenter.ChatTagClient().CChatTag(context.TODO(), &bilin.CChatTagReq{
			Chattag: &bilin.ChatTag{
				Id:       guildRec.ID,
				TagName:  guildRec.TagName,
				TagColor: guildRec.TagColor,
			},
		})
		if err != nil {
			appzaplog.Error("CreateChatTag CreateGuildRoom err", zap.Error(err))
			code = http.StatusInternalServerError
			break
		}
		break
	}

	jsonret := &HttpRetComm{
		Code: code,
		Data: guildRec,
	}

	appzaplog.Debug("[-]CreateChatTag", zap.Any("resp", jsonret))
	c.JSON(200, jsonret)
}

func DelChatTag(c *gin.Context) {
	appzaplog.Debug("[+]DelChatTag")
	code := http.StatusOK

	for {
		roomid := c.Param("id")
		if roomid == "" {
			appzaplog.Error("DelChatTag empty roomid err")
			code = http.StatusBadRequest
			break
		}
		appzaplog.Debug("DelChatTag", zap.String("roomid", roomid))
		roomidU64, err := strconv.ParseInt(roomid, 10, 64)
		if err != nil {
			appzaplog.Error("DelChatTag ParseUint err")
			code = http.StatusBadRequest
			break
		}
		_, err = clientcenter.ChatTagClient().DChatTag(context.TODO(), &bilin.DChatTagReq{
			Chattag: &bilin.ChatTag{
				Id: roomidU64,
			},
		})
		if err != nil {
			appzaplog.Error("DelChatTag DelCategoryGuildRec err", zap.Error(err))
			code = http.StatusInternalServerError
			break
		}
		break
	}

	jsonret := &HttpRetComm{
		Code: code,
	}

	appzaplog.Debug("[-]DelChatTag", zap.Any("resp", jsonret))
	c.JSON(200, jsonret)
}

func UpdateChatTag(c *gin.Context) {
	appzaplog.Debug("[+]UpdateChatTag")
	code := http.StatusOK
	guildRec := &ChatTag{}

	for {
		if err := c.BindJSON(guildRec); err != nil {
			appzaplog.Error("UpdateChatTag BindJSON err", zap.Error(err))
			code = http.StatusBadRequest
			break
		}
		appzaplog.Debug("UpdateChatTag", zap.Any("req", guildRec))

		_, err := clientcenter.ChatTagClient().UChatTag(context.TODO(), &bilin.UChatTagReq{
			Chattag: &bilin.ChatTag{
				Id:       guildRec.ID,
				TagName:  guildRec.TagName,
				TagColor: guildRec.TagColor,
			},
		})
		if err != nil {
			appzaplog.Error("UpdateChatTag BindJSON err", zap.Error(err))
			code = http.StatusInternalServerError
			break
		}
		break
	}

	jsonret := &HttpRetComm{
		Code: code,
		Data: guildRec,
	}

	appzaplog.Debug("[-]UpdateChatTag", zap.Any("resp", jsonret))
	c.JSON(200, jsonret)
}
