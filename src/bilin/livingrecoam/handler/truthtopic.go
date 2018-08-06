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
	"strings"
)

type TruthTopic struct {
	ID    int64    `json:"id"`
	Topic []string `json:"topic"`
}

type TruthTopicList struct {
	TotalPageSize int64        `json:"total_page_size"`
	Items         []TruthTopic `json:"items"`
}

func GetTruthTopic(c *gin.Context) {
	appzaplog.Debug("[+]GetTruthTopic")
	code := http.StatusOK
	guildRec := &TruthTopicList{}

	for {
		pageinfo, err := pageInfo(c)
		if err != nil {
			appzaplog.Error("GetTruthTopic pageInfo err", zap.Error(err))
			code = http.StatusBadRequest
			break
		}
		info, err := clientcenter.ChatTagClient().RTruthTopic(context.TODO(), &bilin.RTruthTopicReq{
			Page: &bilin.PageInfo{
				Pagenum:  int64(pageinfo.pageNum),
				Pagesize: int64(pageinfo.pageSize),
			},
		})
		if err != nil {
			appzaplog.Error("GetTruthTopic GuildRoomS err", zap.Error(err))
			code = http.StatusInternalServerError
			break
		}
		guildRec.TotalPageSize = info.Totalpagenum
		for _, v := range info.Info {
			guildRec.Items = append(guildRec.Items, TruthTopic{
				Topic: strings.Split(v.Topic, "\n"),
				ID:    v.Id,
			})
		}

		//if len(guildRec.Items) > 0 {
		//	start, end, err := pageinfo.StartEnd(len(guildRec.Items))
		//	if err != nil {
		//		appzaplog.Error("GetTruthTopic StartEnd err", zap.Error(err))
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

	appzaplog.Debug("[-]GetTruthTopic", zap.Any("resp", jsonret))
	c.JSON(200, jsonret)
}

func CreateTruthTopic(c *gin.Context) {
	appzaplog.Debug("[+]CreateTruthTopic")
	code := http.StatusOK
	guildRec := &TruthTopic{}

	for {
		if err := c.BindJSON(guildRec); err != nil {
			appzaplog.Error("CreateTruthTopic BindJSON err", zap.Error(err))
			code = http.StatusBadRequest
			break
		}
		appzaplog.Debug("CreateTruthTopic", zap.Any("req", guildRec))

		_, err := clientcenter.ChatTagClient().CTruthTopic(context.TODO(), &bilin.CTruthTopicReq{
			Info: &bilin.TruthTopic{
				Id:    guildRec.ID,
				Topic: strings.Join(guildRec.Topic, "\n"),
			},
		})
		if err != nil {
			appzaplog.Error("CreateTruthTopic CreateGuildRoom err", zap.Error(err))
			code = http.StatusInternalServerError
			break
		}
		break
	}

	jsonret := &HttpRetComm{
		Code: code,
		Data: guildRec,
	}

	appzaplog.Debug("[-]CreateTruthTopic", zap.Any("resp", jsonret))
	c.JSON(200, jsonret)
}

func UpdateTruthTopic(c *gin.Context) {
	appzaplog.Debug("[+]UpdateTruthTopic")
	code := http.StatusOK
	guildRec := &TruthTopic{}

	for {
		if err := c.BindJSON(guildRec); err != nil {
			appzaplog.Error("UpdateTruthTopic BindJSON err", zap.Error(err))
			code = http.StatusBadRequest
			break
		}
		appzaplog.Debug("UpdateTruthTopic", zap.Any("req", guildRec))

		_, err := clientcenter.ChatTagClient().UTruthTopic(context.TODO(), &bilin.UTruthTopicReq{
			Info: &bilin.TruthTopic{
				Id:    guildRec.ID,
				Topic: strings.Join(guildRec.Topic, "\n"),
			},
		})
		if err != nil {
			appzaplog.Error("UpdateTruthTopic UTruthTopic err", zap.Error(err))
			code = http.StatusInternalServerError
			break
		}
		break
	}

	jsonret := &HttpRetComm{
		Code: code,
		Data: guildRec,
	}

	appzaplog.Debug("[-]UpdateTruthTopic", zap.Any("resp", jsonret))
	c.JSON(200, jsonret)
}

func DelTruthTopic(c *gin.Context) {
	appzaplog.Debug("[+]DelTruthTopic")
	code := http.StatusOK

	for {
		roomid := c.Param("id")
		if roomid == "" {
			appzaplog.Error("DelTruthTopic empty roomid err")
			code = http.StatusBadRequest
			break
		}
		appzaplog.Debug("DelTruthTopic", zap.String("roomid", roomid))
		roomidU64, err := strconv.ParseInt(roomid, 10, 64)
		if err != nil {
			appzaplog.Error("DelTruthTopic ParseUint err")
			code = http.StatusBadRequest
			break
		}
		_, err = clientcenter.ChatTagClient().DTruthTopic(context.TODO(), &bilin.DTruthTopicReq{
			Info: &bilin.TruthTopic{
				Id: roomidU64,
			},
		})
		if err != nil {
			appzaplog.Error("DelTruthTopic DelCategoryGuildRec err", zap.Error(err))
			code = http.StatusInternalServerError
			break
		}
		break
	}

	jsonret := &HttpRetComm{
		Code: code,
	}

	appzaplog.Debug("[-]DelTruthTopic", zap.Any("resp", jsonret))
	c.JSON(200, jsonret)
}
