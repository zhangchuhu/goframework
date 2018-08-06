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
	"time"
)

type Contract struct {
	ID                   int64  `json:"id"`
	GuildID              int64  `json:"guild_id"`
	HostUid              int64  `json:"host_uid"`
	ContractStartTime    string `json:"contract_start_time"`
	ContractEndTime      string `json:"contract_end_time"`
	GuildSharePercentage int64  `json:"guild_share_percentage"`
	HostSharePercentage  int64  `json:"host_share_percentage"`
	ContractState        int32  `json:"contract_state"`
}

type ContractList struct {
	TotalPageSize int        `json:"total_page_size"`
	Items         []Contract `json:"items"`
}

func GetContract(c *gin.Context) {
	appzaplog.Debug("[+]GetContract")
	code := http.StatusOK
	guildRec := &ContractList{}

	for {
		pageinfo, err := pageInfo(c)
		if err != nil {
			appzaplog.Error("GetContract pageInfo err", zap.Error(err))
			code = http.StatusBadRequest
			break
		}
		info, err := clientcenter.GuildTarsClient().RContract(context.TODO(), &bilin.RContractReq{})
		if err != nil {
			appzaplog.Error("GetContract GuildRoomS err", zap.Error(err))
			code = http.StatusInternalServerError
			break
		}
		guildRec.TotalPageSize = len(info.Info)
		for _, v := range info.Info {
			guildRec.Items = append(guildRec.Items, Contract{
				ID:                   v.Id,
				GuildID:              v.Guildid,
				HostUid:              v.Hostuid,
				ContractStartTime:    time.Unix(v.Contractstarttime, 0).Format(time.RFC3339Nano),
				ContractEndTime:      time.Unix(v.Contractendtime, 0).Format(time.RFC3339Nano),
				GuildSharePercentage: v.Guildsharepercentage,
				HostSharePercentage:  v.Hostsharepercentage,
				ContractState:        v.Contractstate,
			})
		}

		if len(guildRec.Items) > 0 {
			start, end, err := pageinfo.StartEnd(len(guildRec.Items))
			if err != nil {
				appzaplog.Error("GetContract StartEnd err", zap.Error(err))
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

	appzaplog.Debug("[-]GetContract", zap.Any("resp", jsonret))
	c.JSON(200, jsonret)
}

func CreateContract(c *gin.Context) {
	appzaplog.Debug("[+]CreateContract")
	code := http.StatusOK
	guildRec := &Contract{}

	for {
		if err := c.BindJSON(guildRec); err != nil {
			appzaplog.Error("CreateContract BindJSON err", zap.Error(err))
			code = http.StatusBadRequest
			break
		}
		appzaplog.Debug("CreateContract", zap.Any("req", guildRec))
		if guildRec.HostUid <= 0 {
			appzaplog.Error("not hostuid dedicate", zap.Any("req", guildRec))
			code = http.StatusBadRequest
			break
		}
		startTime, err := time.ParseInLocation(time.RFC3339Nano, guildRec.ContractStartTime, time.Local)
		if err != nil {
			appzaplog.Error("CreateContract ParseInLocation err", zap.Error(err), zap.Any("req", guildRec))
			code = http.StatusBadRequest
			break
		}

		endTime, err := time.ParseInLocation(time.RFC3339Nano, guildRec.ContractEndTime, time.Local)
		if err != nil {
			appzaplog.Error("CreateContract ParseInLocation err", zap.Error(err), zap.Any("req", guildRec))
			code = http.StatusBadRequest
			break
		}
		mycontract, err := clientcenter.GuildTarsClient().RContract(context.TODO(), &bilin.RContractReq{
			Filter: &bilin.Contract{
				Hostuid: guildRec.HostUid,
			},
		})
		if err != nil {
			appzaplog.Error("CreateContract CreateGuildRoom err", zap.Error(err))
			code = http.StatusInternalServerError
			break
		}
		if mycontract.Info != nil {
			appzaplog.Warn("already contract host", zap.Any("req", guildRec))
			code = http.StatusNotAcceptable
			break
		}
		_, err = clientcenter.GuildTarsClient().CContract(context.TODO(), &bilin.CContractReq{
			Info: &bilin.Contract{
				Id:                   guildRec.ID,
				Guildid:              guildRec.GuildID,
				Hostuid:              guildRec.HostUid,
				Contractstarttime:    startTime.Unix(),
				Contractendtime:      endTime.Unix(),
				Guildsharepercentage: guildRec.GuildSharePercentage,
				Hostsharepercentage:  guildRec.HostSharePercentage,
				Contractstate:        guildRec.ContractState,
			},
		})
		if err != nil {
			appzaplog.Error("CreateContract CreateGuildRoom err", zap.Error(err))
			code = http.StatusInternalServerError
			break
		}
		break
	}

	jsonret := &HttpRetComm{
		Code: code,
		Data: guildRec,
	}

	appzaplog.Debug("[-]CreateContract", zap.Any("resp", jsonret))
	c.JSON(200, jsonret)
}

func UpdateContract(c *gin.Context) {
	appzaplog.Debug("[+]UpdateContract")
	code := http.StatusOK
	guildRec := &Contract{}

	for {
		if err := c.BindJSON(guildRec); err != nil {
			appzaplog.Error("UpdateContract BindJSON err", zap.Error(err))
			code = http.StatusBadRequest
			break
		}
		appzaplog.Debug("UpdateContract", zap.Any("req", guildRec))
		startTime, err := time.ParseInLocation(time.RFC3339Nano, guildRec.ContractStartTime, time.Local)
		if err != nil {
			appzaplog.Error("CreateContract ParseInLocation err", zap.Error(err), zap.Any("req", guildRec))
			code = http.StatusBadRequest
			break
		}

		endTime, err := time.ParseInLocation(time.RFC3339Nano, guildRec.ContractEndTime, time.Local)
		if err != nil {
			appzaplog.Error("CreateContract ParseInLocation err", zap.Error(err), zap.Any("req", guildRec))
			code = http.StatusBadRequest
			break
		}

		_, err = clientcenter.GuildTarsClient().UContract(context.TODO(), &bilin.UContractReq{
			Info: &bilin.Contract{
				Id:                   guildRec.ID,
				Guildid:              guildRec.GuildID,
				Hostuid:              guildRec.HostUid,
				Contractstarttime:    startTime.Unix(),
				Contractendtime:      endTime.Unix(),
				Guildsharepercentage: guildRec.GuildSharePercentage,
				Hostsharepercentage:  guildRec.HostSharePercentage,
				Contractstate:        guildRec.ContractState,
			},
		})
		if err != nil {
			appzaplog.Error("UpdateContract BindJSON err", zap.Error(err))
			code = http.StatusInternalServerError
			break
		}
		break
	}

	jsonret := &HttpRetComm{
		Code: code,
		Data: guildRec,
	}

	appzaplog.Debug("[-]UpdateContract", zap.Any("resp", jsonret))
	c.JSON(200, jsonret)
}

func DelContract(c *gin.Context) {
	appzaplog.Debug("[+]DelContract")
	code := http.StatusOK

	for {
		roomid := c.Param("id")
		if roomid == "" {
			appzaplog.Error("DelContract empty roomid err")
			code = http.StatusBadRequest
			break
		}
		appzaplog.Debug("DelContract", zap.String("roomid", roomid))
		roomidU64, err := strconv.ParseInt(roomid, 10, 64)
		if err != nil {
			appzaplog.Error("DelContract ParseUint err")
			code = http.StatusBadRequest
			break
		}
		_, err = clientcenter.GuildTarsClient().DContract(context.TODO(), &bilin.DContractReq{
			Info: &bilin.Contract{
				Id: roomidU64,
			},
		})
		if err != nil {
			appzaplog.Error("DelContract DelCategoryGuildRec err", zap.Error(err))
			code = http.StatusInternalServerError
			break
		}
		break
	}

	jsonret := &HttpRetComm{
		Code: code,
	}

	appzaplog.Debug("[-]DelContract", zap.Any("resp", jsonret))
	c.JSON(200, jsonret)
}
