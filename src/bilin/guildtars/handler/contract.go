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

// 2038/01/19 03:14:07
//超过这个值在mysql不能保存，改为这个值
const maxTimeStampUnixSec = 2147454847

// 签约信息CRUD
func (this *GuildTarsObj) CContract(ctx context.Context, r *bilin.CContractReq) (*bilin.CContractResp, error) {
	appzaplog.Debug("[+]CContract", zap.Any("req", r))
	code := int64(0)
	defer func(now time.Time) {
		httpmetrics.DefReport("CContract", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	resp := &bilin.CContractResp{}
	tag := pb2daoContract(r.Info)
	if tag == nil {
		appzaplog.Warn("req info null")
		return resp, fmt.Errorf("req info null")
	}
	if err := tag.Create(); err != nil {
		appzaplog.Error("CContract Create err", zap.Error(err))
		code = CreateContractFailed
		return resp, err
	}
	appzaplog.Debug("[-]CContract", zap.Any("resp", resp))
	return resp, nil
}

func (this *GuildTarsObj) RContract(ctx context.Context, r *bilin.RContractReq) (*bilin.RContractResp, error) {
	code := int64(0)
	defer func(now time.Time) {
		httpmetrics.DefReport("RContract", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	resp := &bilin.RContractResp{}
	var (
		info []dao.Contract
		err  error
	)
	c := pb2daoContract(r.Filter)
	if c == nil {
		info, err = dao.GetAllContract()
	} else {
		info, err = c.Get()
	}

	if err != nil {
		appzaplog.Error("RContract Get err", zap.Error(err), zap.Any("req", r))
		code = GetContractFailed
		return resp, err
	}
	for _, v := range info {
		resp.Info = append(resp.Info, &bilin.Contract{
			Id:                   int64(v.ID),
			Guildid:              v.GuildID,
			Hostuid:              v.HostUid,
			Contractstarttime:    v.ContractStartTime.Unix(),
			Contractendtime:      v.ContractEndTime.Unix(),
			Guildsharepercentage: v.GuildSharePercentage,
			Hostsharepercentage:  v.HostSharePercentage,
			Contractstate:        v.ContractState,
		})
	}
	appzaplog.Debug("[-]RContract", zap.Any("resp", resp))
	return resp, nil
}

func (this *GuildTarsObj) UContract(ctx context.Context, r *bilin.UContractReq) (*bilin.UContractResp, error) {
	code := int64(0)
	defer func(now time.Time) {
		httpmetrics.DefReport("UContract", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	resp := &bilin.UContractResp{}
	c := pb2daoContract(r.Info)
	if c == nil {
		appzaplog.Warn("req info null")
		return resp, fmt.Errorf("req info null")
	}
	if err := c.Update(); err != nil {
		appzaplog.Error("UContract Update err", zap.Error(err), zap.Any("req", r))
		code = UpdateContractFailed
		return resp, err
	}
	appzaplog.Debug("[-]UContract", zap.Any("resp", resp))
	return resp, nil
}
func (this *GuildTarsObj) DContract(ctx context.Context, r *bilin.DContractReq) (*bilin.DContractResp, error) {
	code := int64(0)
	defer func(now time.Time) {
		httpmetrics.DefReport("DContract", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	resp := &bilin.DContractResp{}
	c := pb2daoContract(r.Info)
	if c == nil {
		appzaplog.Warn("DContract req info null")
		return resp, fmt.Errorf("DContract req info null")
	}
	if err := c.Delete(); err != nil {
		appzaplog.Error("DContract Delete err", zap.Error(err), zap.Any("req", r))
		code = DelContractFailed
		return resp, err
	}
	appzaplog.Debug("[-]DContract", zap.Any("resp", resp))
	return resp, nil
}

func pb2daoContract(c *bilin.Contract) *dao.Contract {
	if c == nil {
		return nil
	}
	ret := &dao.Contract{
		GuildID:              c.Guildid,
		HostUid:              c.Hostuid,
		GuildSharePercentage: c.Guildsharepercentage,
		HostSharePercentage:  c.Hostsharepercentage,
		ContractState:        c.Contractstate,
	}
	ret.ID = uint(c.Id)
	if c.Contractstarttime > 0 {
		ret.ContractStartTime = time.Unix(c.Contractstarttime, 0)
	}
	if c.Contractendtime > 0 {
		//todo 使用unix时间戳代替mysql的timestamp
		if c.Contractendtime > maxTimeStampUnixSec {
			c.Contractendtime = maxTimeStampUnixSec
			appzaplog.Warn("timestamp overflow", zap.Int64("endtime", c.Contractendtime))
		}
		ret.ContractEndTime = time.Unix(c.Contractendtime, 0)
	}
	appzaplog.Debug("[-]pb2daoContract", zap.Any("resp", ret))
	return ret
}
