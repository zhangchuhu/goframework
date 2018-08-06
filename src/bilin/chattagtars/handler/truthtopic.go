/*
 * Copyright (c) 2018-07-05.
 * Author: kordenlu
 * 功能描述:真心话CRUD
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

//
// 真心话CRUD
func (p *ChatTagTarsObj) CTruthTopic(ctx context.Context, r *bilin.CTruthTopicReq) (*bilin.CTruthTopicResp, error) {
	appzaplog.Debug("[+]CTruthTopic")
	code := int64(MetricCodeSuccess)
	defer func(now time.Time) {
		httpmetrics.DefReport("CTruthTopic", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	resp := &bilin.CTruthTopicResp{}
	tag := dao.TruthTopic{
		Topic: r.Info.Topic,
	}
	if err := tag.Create(); err != nil {
		appzaplog.Error("CTruthTopic Create err", zap.Error(err))
		code = MetricCodeCreateErr
		return resp, err
	}
	appzaplog.Debug("[-]CTruthTopic", zap.Any("resp", resp))
	return resp, nil
}

func (p *ChatTagTarsObj) RTruthTopic(ctx context.Context, r *bilin.RTruthTopicReq) (*bilin.RTruthTopicResp, error) {
	appzaplog.Debug("[+]RTruthTopic")
	var (
		info       []dao.TruthTopic
		totalCount int64
		err        error
	)
	code := int64(MetricCodeSuccess)
	defer func(now time.Time) {
		httpmetrics.DefReport("RTruthTopic", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	resp := &bilin.RTruthTopicResp{}
	if r.Page != nil {
		page := int64(1)
		pagesize := int64(20)
		if r != nil && r.Page != nil {
			page = r.Page.Pagenum
			pagesize = r.Page.Pagesize
		}
		info, totalCount, err = dao.GetAllTruthTopicByPage(page, pagesize)
	} else {
		info, err = dao.GetAllTruthTopic()
		totalCount = int64(len(info))
	}

	if err != nil {
		appzaplog.Error("RTruthTopic dao.GetAllTruthTopicByPage err", zap.Error(err))
		code = MetricCodeReadErr
		return resp, err
	}

	resp.Totalpagenum = totalCount
	for _, v := range info {
		resp.Info = append(resp.Info, &bilin.TruthTopic{
			Id:    int64(v.ID),
			Topic: v.Topic,
		})
	}
	appzaplog.Debug("[-]RTruthTopic", zap.Any("resp", resp))
	return resp, nil
}
func (p *ChatTagTarsObj) UTruthTopic(ctx context.Context, r *bilin.UTruthTopicReq) (*bilin.UTruthTopicResp, error) {
	appzaplog.Debug("[+]UTruthTopic")
	code := int64(MetricCodeSuccess)
	defer func(now time.Time) {
		httpmetrics.DefReport("UTruthTopic", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	resp := &bilin.UTruthTopicResp{}
	tag := dao.TruthTopic{
		Topic: r.Info.Topic,
	}
	tag.ID = uint(r.Info.Id)
	if err := tag.Update(); err != nil {
		appzaplog.Error("UTruthTopic Update err", zap.Error(err))
		code = MetricCodeUpdateErr
		return resp, err
	}
	appzaplog.Debug("[-]UTruthTopic", zap.Any("resp", resp))
	return resp, nil
}
func (p *ChatTagTarsObj) DTruthTopic(ctx context.Context, r *bilin.DTruthTopicReq) (*bilin.DTruthTopicResp, error) {
	appzaplog.Debug("[+]DTruthTopic")
	code := int64(MetricCodeSuccess)
	defer func(now time.Time) {
		httpmetrics.DefReport("DTruthTopic", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	resp := &bilin.DTruthTopicResp{}
	tag := dao.TruthTopic{}
	tag.ID = uint(r.Info.Id)
	if tag.ID == 0 {
		appzaplog.Warn("DTruthTopic delall not allowed")
		code = MetricCodeDelAllWarn
		return resp, fmt.Errorf("delall not allowed")
	}
	if err := tag.Del(); err != nil {
		appzaplog.Error("DTruthTopic Del err", zap.Error(err))
		code = MetricCodeDelErr
		return resp, err
	}
	appzaplog.Debug("[-]DTruthTopic", zap.Any("resp", resp))
	return resp, nil
}
