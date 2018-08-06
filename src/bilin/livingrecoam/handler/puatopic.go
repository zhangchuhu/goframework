/*
 * Copyright (c) 2018-07-04.
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
	"strings"
)

type PUATopic struct {
	ID    int64    `json:"id"`
	Topic []string `json:"topic"`
}

type PUATopicList struct {
	TotalPageSize int64      `json:"total_page_size"`
	Items         []PUATopic `json:"items"`
}

func GetPUATopic(c *gin.Context) {
	appzaplog.Debug("[+]GetPUATopic")
	code := http.StatusOK
	guildRec := &PUATopicList{}

	for {
		pageinfo, err := pageInfo(c)
		if err != nil {
			appzaplog.Error("GetPUATopic pageInfo err", zap.Error(err))
			code = http.StatusBadRequest
			break
		}
		info, err := clientcenter.ChatTagClient().RPUATopic(context.TODO(), &bilin.RPUATopicReq{
			Page: &bilin.PageInfo{
				Pagenum:  int64(pageinfo.pageNum),
				Pagesize: int64(pageinfo.pageSize),
			},
		})
		if err != nil {
			appzaplog.Error("GetPUATopic GuildRoomS err", zap.Error(err))
			code = http.StatusInternalServerError
			break
		}
		guildRec.TotalPageSize = info.Totalpagenum
		for _, v := range info.Info {
			guildRec.Items = append(guildRec.Items, PUATopic{
				Topic: strings.Split(v.Topic, "\n"),
				ID:    v.Id,
			})
		}

		//if len(guildRec.Items) > 0 {
		//	start, end, err := pageinfo.StartEnd(len(guildRec.Items))
		//	if err != nil {
		//		appzaplog.Error("GetGuildRoom StartEnd err", zap.Error(err))
		//		code = http.StatusBadRequest
		//		break
		//	}
		//	guildRec.Items = guildRec.Items[start:end]
		//}
		break
	}

	jsonret := &HttpRetComm{
		Code: code,
		Data: guildRec,
	}

	appzaplog.Debug("[-]GetGuildRoom", zap.Any("resp", jsonret))
	c.JSON(200, jsonret)
}

func CreatePUATopic(c *gin.Context) {
	appzaplog.Debug("[+]CreatePUATopic")
	code := http.StatusOK
	guildRec := &PUATopic{}

	for {
		if err := c.BindJSON(guildRec); err != nil {
			appzaplog.Error("CreatePUATopic BindJSON err", zap.Error(err))
			code = http.StatusBadRequest
			break
		}
		appzaplog.Debug("CreatePUATopic", zap.Any("req", guildRec))

		_, err := clientcenter.ChatTagClient().CPUATopic(context.TODO(), &bilin.CPUATopicReq{
			Info: &bilin.PUATopic{
				Id:    guildRec.ID,
				Topic: strings.Join(guildRec.Topic, "\n"),
			},
		})
		if err != nil {
			appzaplog.Error("CreatePUATopic CreateGuildRoom err", zap.Error(err))
			code = http.StatusInternalServerError
			break
		}
		break
	}

	jsonret := &HttpRetComm{
		Code: code,
		Data: guildRec,
	}

	appzaplog.Debug("[-]CreatePUATopic", zap.Any("resp", jsonret))
	c.JSON(200, jsonret)
}

func UpdatePUATopic(c *gin.Context) {
	appzaplog.Debug("[+]UpdatePUATopic")
	code := http.StatusOK
	guildRec := &PUATopic{}

	for {
		if err := c.BindJSON(guildRec); err != nil {
			appzaplog.Error("UpdatePUATopic BindJSON err", zap.Error(err))
			code = http.StatusBadRequest
			break
		}
		appzaplog.Debug("UpdatePUATopic", zap.Any("req", guildRec))

		_, err := clientcenter.ChatTagClient().UPUATopic(context.TODO(), &bilin.UPUATopicReq{
			Info: &bilin.PUATopic{
				Id:    guildRec.ID,
				Topic: strings.Join(guildRec.Topic, "\n"),
			},
		})
		if err != nil {
			appzaplog.Error("UpdatePUATopic UPUATopic err", zap.Error(err))
			code = http.StatusInternalServerError
			break
		}
		break
	}

	jsonret := &HttpRetComm{
		Code: code,
		Data: guildRec,
	}

	appzaplog.Debug("[-]UpdatePUATopic", zap.Any("resp", jsonret))
	c.JSON(200, jsonret)
}

func DelPUATopic(c *gin.Context) {
	appzaplog.Debug("[+]DelPUATopic")
	code := http.StatusOK

	for {
		roomid := c.Param("id")
		if roomid == "" {
			appzaplog.Error("DelPUATopic empty roomid err")
			code = http.StatusBadRequest
			break
		}
		appzaplog.Debug("DelPUATopic", zap.String("roomid", roomid))
		roomidU64, err := strconv.ParseInt(roomid, 10, 64)
		if err != nil {
			appzaplog.Error("DelPUATopic ParseUint err")
			code = http.StatusBadRequest
			break
		}
		_, err = clientcenter.ChatTagClient().DPUATopic(context.TODO(), &bilin.DPUATopicReq{
			Info: &bilin.PUATopic{
				Id: roomidU64,
			},
		})
		if err != nil {
			appzaplog.Error("DelPUATopic DelCategoryGuildRec err", zap.Error(err))
			code = http.StatusInternalServerError
			break
		}
		break
	}

	jsonret := &HttpRetComm{
		Code: code,
	}

	appzaplog.Debug("[-]DelPUATopic", zap.Any("resp", jsonret))
	c.JSON(200, jsonret)
}
