/*
 * Copyright (c) 2018-07-18.
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

// 工会信息CRUD
func (this *GuildTarsObj) CGuild(ctx context.Context, r *bilin.CGuildReq) (*bilin.CGuildResp, error) {
	code := int64(0)
	defer func(now time.Time) {
		httpmetrics.DefReport("CGuild", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	resp := &bilin.CGuildResp{}
	tag := pb2daoGuild(r.Info)
	if tag == nil {
		appzaplog.Warn("req info null")
		return resp, fmt.Errorf("req info null")
	}
	if err := tag.Create(); err != nil {
		appzaplog.Error("CGuild Create err", zap.Error(err))
		code = CreateGuildFailed
		return resp, err
	}
	appzaplog.Debug("[-]CGuild", zap.Any("resp", resp))
	return resp, nil
}
func (this *GuildTarsObj) RGuild(ctx context.Context, r *bilin.RGuildReq) (*bilin.RGuildResp, error) {
	code := int64(0)
	defer func(now time.Time) {
		httpmetrics.DefReport("RGuild", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	resp := &bilin.RGuildResp{}
	c := pb2daoGuild(r.Filter)
	var (
		info []dao.Guild
		err  error
	)
	if c == nil {
		info, err = dao.GetAllGuild()
	} else {
		info, err = c.Get()
	}
	if err != nil {
		appzaplog.Error("RGuild Get err", zap.Error(err), zap.Any("req", r))
		code = GetGuildFailed
		return resp, err
	}
	for _, v := range info {
		resp.Info = append(resp.Info, &bilin.Guild{
			Id:        int64(v.ID),
			Title:     v.Title,
			Ow:        v.OW,
			Mobile:    v.Mobile,
			Describle: v.Describle,
			Guildlog:  v.GuildLogo,
		})
	}
	appzaplog.Debug("[-]RGuild", zap.Any("resp", resp))
	return resp, nil
}
func (this *GuildTarsObj) UGuild(ctx context.Context, r *bilin.UGuildReq) (*bilin.UGuildResp, error) {
	code := int64(0)
	defer func(now time.Time) {
		httpmetrics.DefReport("UGuild", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	resp := &bilin.UGuildResp{}
	c := pb2daoGuild(r.Info)
	if c == nil {
		appzaplog.Warn("UGuild req info null")
		return resp, fmt.Errorf("req info null")
	}
	if err := c.Update(); err != nil {
		appzaplog.Error("UGuild Update err", zap.Error(err), zap.Any("req", r))
		code = UpdateGuildFailed
		return resp, err
	}
	appzaplog.Debug("[-]UGuild", zap.Any("resp", resp))
	return resp, nil
}

func (this *GuildTarsObj) DGuild(ctx context.Context, r *bilin.DGuildReq) (*bilin.DGuildResp, error) {
	code := int64(0)
	defer func(now time.Time) {
		httpmetrics.DefReport("DGuild", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	resp := &bilin.DGuildResp{}
	c := pb2daoGuild(r.Info)
	if c == nil {
		appzaplog.Warn("DGuild req info null")
		return resp, fmt.Errorf("req info null")
	}
	if err := c.Delete(); err != nil {
		appzaplog.Error("DGuild Delete err", zap.Error(err), zap.Any("req", r))
		code = DelGuildFailed
		return resp, err
	}
	appzaplog.Debug("[-]DGuild", zap.Any("resp", resp))
	return resp, nil
}

func pb2daoGuild(guild *bilin.Guild) *dao.Guild {
	if guild == nil {
		return nil
	}
	ret := &dao.Guild{
		OW:        guild.Ow,
		Title:     guild.Title,
		Mobile:    guild.Mobile,
		Describle: guild.Describle,
		GuildLogo: guild.Guildlog,
	}
	ret.ID = uint(guild.Id)
	return ret
}
