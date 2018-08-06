/*
 * Copyright (c) 2018-07-20.
 * Author: kordenlu
 * 功能描述:${<VARIABLE_NAME>}
 */

package handler

import (
	"bilin/guildtars/dao"
	"bilin/protocol"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/kissgo/httpmetrics"
	"context"
	"fmt"
	"time"
)

// 工会房间CRUD
func (this *GuildTarsObj) CGuildRoom(ctx context.Context, r *bilin.CGuildRoomReq) (*bilin.CGuildRoomResp, error) {
	code := int64(0)
	defer func(now time.Time) {
		httpmetrics.DefReport("CGuildRoom", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	resp := &bilin.CGuildRoomResp{}
	tag := pb2daoGuildRoom(r.Info)
	if tag == nil {
		appzaplog.Warn("CGuildRoom req info null")
		return resp, fmt.Errorf("req info null")
	}
	if err := tag.Create(); err != nil {
		appzaplog.Error("CGuildRoom Create err", zap.Error(err))
		code = CreateGuildRoomFailed
		return resp, err
	}
	appzaplog.Debug("[-]CGuildRoom", zap.Any("resp", resp))
	return resp, nil
}

func (this *GuildTarsObj) RGuildRoom(ctx context.Context, r *bilin.RGuildRoomReq) (*bilin.RGuildRoomResp, error) {
	code := int64(0)
	defer func(now time.Time) {
		httpmetrics.DefReport("RGuildRoom", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	resp := &bilin.RGuildRoomResp{}
	c := pb2daoGuildRoom(r.Filter)
	var (
		info []dao.GuildRoom
		err  error
	)
	if c == nil {
		info, err = dao.GetGuildRoomS()
	} else {
		info, err = c.Get()
	}
	if err != nil {
		appzaplog.Error("RGuildRoom Get err", zap.Error(err), zap.Any("req", r))
		code = GetGuildRoomFailed
		return resp, err
	}
	for _, v := range info {
		resp.Info = append(resp.Info, &bilin.GuildRoom{
			Id:      int64(v.ID),
			Guildid: v.GuildID,
			Roomid:  v.RoomID,
		})
	}
	appzaplog.Debug("[-]RGuildRoom", zap.Any("resp", resp))
	return resp, nil
}

func (this *GuildTarsObj) UGuildRoom(ctx context.Context, r *bilin.UGuildRoomReq) (*bilin.UGuildRoomResp, error) {
	code := int64(0)
	defer func(now time.Time) {
		httpmetrics.DefReport("UGuildRoom", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	resp := &bilin.UGuildRoomResp{}
	c := pb2daoGuildRoom(r.Info)
	if c == nil {
		appzaplog.Warn("UGuildRoom req info null")
		return resp, fmt.Errorf("req info null")
	}
	if err := c.Update(); err != nil {
		appzaplog.Error("UGuildRoom Update err", zap.Error(err), zap.Any("req", r))
		code = UpdateGuildRoomFailed
		return resp, err
	}
	appzaplog.Debug("[-]UGuildRoom", zap.Any("resp", resp))
	return resp, nil
}

func (this *GuildTarsObj) DGuildRoom(ctx context.Context, r *bilin.DGuildRoomReq) (*bilin.DGuildRoomResp, error) {
	code := int64(0)
	defer func(now time.Time) {
		httpmetrics.DefReport("DGuildRoom", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	resp := &bilin.DGuildRoomResp{}
	c := pb2daoGuildRoom(r.Info)
	if c == nil {
		appzaplog.Warn("DGuildRoom req info null")
		return resp, fmt.Errorf("req info null")
	}
	if err := c.Delete(); err != nil {
		appzaplog.Error("DGuildRoom Delete err", zap.Error(err), zap.Any("req", r))
		code = DelGuildFailed
		return resp, err
	}
	appzaplog.Debug("[-]DGuildRoom", zap.Any("resp", resp))
	return resp, nil
}

func pb2daoGuildRoom(guild *bilin.GuildRoom) *dao.GuildRoom {
	if guild == nil {
		return nil
	}
	ret := &dao.GuildRoom{
		GuildID: guild.Guildid,
		RoomID:  guild.Roomid,
	}
	ret.ID = uint(guild.Id)
	return ret
}
