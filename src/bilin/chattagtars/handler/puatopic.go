/*
 * Copyright (c) 2018-07-05.
 * Author: kordenlu
 * 功能描述:PUA聊妹套话CRUD
 */

package handler

import (
	"bilin/chattagtars/dao"
	"bilin/protocol"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/kissgo/httpmetrics"
	"context"
	"fmt"
	"time"
)

func (p *ChatTagTarsObj) CPUATopic(ctx context.Context, r *bilin.CPUATopicReq) (*bilin.CPUATopicResp, error) {
	appzaplog.Debug("[+]CPUATopic")
	code := int64(MetricCodeSuccess)
	defer func(now time.Time) {
		httpmetrics.DefReport("CPUATopic", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	resp := &bilin.CPUATopicResp{}
	tag := dao.PuaTopic{
		Topic: r.Info.Topic,
	}
	if err := tag.Create(); err != nil {
		appzaplog.Error("CPUATopic Create err", zap.Error(err))
		code = MetricCodeCreateErr
		return resp, err
	}
	appzaplog.Debug("[-]CPUATopic", zap.Any("resp", resp))
	return resp, nil
}

func (p *ChatTagTarsObj) RPUATopic(ctx context.Context, r *bilin.RPUATopicReq) (*bilin.RPUATopicResp, error) {
	appzaplog.Debug("[+]RPUATopic")
	code := int64(MetricCodeSuccess)
	defer func(now time.Time) {
		httpmetrics.DefReport("RPUATopic", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	resp := &bilin.RPUATopicResp{}
	var (
		info       []dao.PuaTopic
		totalcount int64
		err        error
	)
	if r.Page != nil {
		page := int64(1)
		pagesize := int64(20)
		if r != nil && r.Page != nil {
			page = r.Page.Pagenum
			pagesize = r.Page.Pagesize
		}
		info, totalcount, err = dao.GetPuaTopicByPage(page, pagesize)
	} else {
		info, err = dao.GetAllPuaTopic()
		totalcount = int64(len(info))
	}
	if err != nil {
		appzaplog.Error("RPUATopic dao.GetPuaTopicByPage err", zap.Error(err))
		code = MetricCodeReadErr
		return resp, err
	}

	resp.Totalpagenum = totalcount
	for _, v := range info {
		resp.Info = append(resp.Info, &bilin.PUATopic{
			Id:    int64(v.ID),
			Topic: v.Topic,
		})
	}
	appzaplog.Debug("[-]RPUATopic", zap.Any("resp", resp))
	return resp, nil
}
func (p *ChatTagTarsObj) UPUATopic(ctx context.Context, r *bilin.UPUATopicReq) (*bilin.UPUATopicResp, error) {
	appzaplog.Debug("[+]UPUATopic")
	code := int64(MetricCodeSuccess)
	defer func(now time.Time) {
		httpmetrics.DefReport("UPUATopic", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	resp := &bilin.UPUATopicResp{}
	tag := dao.PuaTopic{
		Topic: r.Info.Topic,
	}
	tag.ID = uint(r.Info.Id)
	if err := tag.Update(); err != nil {
		appzaplog.Error("UPUATopic Update err", zap.Error(err))
		code = MetricCodeUpdateErr
		return resp, err
	}
	appzaplog.Debug("[-]UPUATopic", zap.Any("resp", resp))
	return resp, nil
}
func (p *ChatTagTarsObj) DPUATopic(ctx context.Context, r *bilin.DPUATopicReq) (*bilin.DPUATopicResp, error) {
	appzaplog.Debug("[+]DPUATopic")
	code := int64(MetricCodeSuccess)
	defer func(now time.Time) {
		httpmetrics.DefReport("DPUATopic", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	resp := &bilin.DPUATopicResp{}
	tag := dao.PuaTopic{}
	tag.ID = uint(r.Info.Id)
	if tag.ID == 0 {
		appzaplog.Warn("DPUATopic delall not allowed")
		code = MetricCodeDelAllWarn
		return resp, fmt.Errorf("delall not allowed")
	}
	if err := tag.Del(); err != nil {
		appzaplog.Error("DPUATopic Del err", zap.Error(err))
		code = MetricCodeDelErr
		return resp, err
	}
	appzaplog.Debug("[-]DChatTag", zap.Any("resp", resp))
	return resp, nil
}
