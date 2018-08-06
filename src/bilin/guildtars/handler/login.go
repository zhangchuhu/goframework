/*
 * Copyright (c) 2018-07-26.
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
	"time"
)

// 工会管理后台登录
func (this *GuildTarsObj) OAMLogin(ctx context.Context, r *bilin.OAMLoginReq) (*bilin.OAMLoginResp, error) {
	code := int64(0)
	defer func(now time.Time) {
		httpmetrics.DefReport("OAMLogin", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	resp := &bilin.OAMLoginResp{}
	oamuser := dao.OAMUser{
		Username: r.Username,
		Passwd:   r.Passwd,
	}
	info, err := oamuser.Get()
	if err != nil {
		appzaplog.Error("OAMLogin Get err", zap.Error(err), zap.Any("req", r))
		code = GetOamUserFailed
		return resp, err
	}
	for _, v := range info {
		switch v.Role {
		case 1:
			resp.Token = "admin"
		case 2:
			resp.Token = "livingrecop"
		case 3:
			resp.Token = "guildop"
		case 4:
			resp.Token = "rcop"
		}
		break
	}
	appzaplog.Debug("[-]RGuildRoom", zap.Any("resp", resp))
	return resp, nil
}
