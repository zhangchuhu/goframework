/*
 * Copyright (c) 2018-07-04.
 * Author: kordenlu
 * 功能描述:${<VARIABLE_NAME>}
 */

package handler

import (
	"bilin/chattagtars/cache"
	"bilin/chattagtars/dao"
	"bilin/protocol"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/kissgo/httpmetrics"
	"context"
	"fmt"
	"time"
)

func (p *ChatTagTarsObj) CChatTag(ctx context.Context, r *bilin.CChatTagReq) (*bilin.CChatTagResp, error) {
	appzaplog.Debug("[+]CChatTag")
	code := int64(MetricCodeSuccess)
	defer func(now time.Time) {
		httpmetrics.DefReport("CChatTag", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	resp := &bilin.CChatTagResp{}
	tag := dao.ChatTag{
		TagName:  r.Chattag.TagName,
		TagColor: r.Chattag.TagColor,
	}
	if err := tag.Create(); err != nil {
		appzaplog.Error("CChatTag Create err", zap.Error(err))
		code = MetricCodeCreateErr
		return resp, err
	}
	appzaplog.Debug("[-]CChatTag", zap.Any("resp", resp))
	return resp, nil
}

func (p *ChatTagTarsObj) RChatTag(ctx context.Context, r *bilin.RChatTagReq) (*bilin.RChatTagResp, error) {
	appzaplog.Debug("[+]RChatTag")
	code := int64(MetricCodeSuccess)
	defer func(now time.Time) {
		httpmetrics.DefReport("RChatTag", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	resp := &bilin.RChatTagResp{}
	info := cache.TakeChatTagCache()
	if info == nil {
		appzaplog.Warn("RChatTag cache.TakeChatTagCache empty")
		code = MetricCodeReadErr
		return resp, nil
	}
	for _, v := range info {
		resp.Chattag = append(resp.Chattag, &bilin.ChatTag{
			Id:       int64(v.ID),
			TagName:  v.TagName,
			TagColor: v.TagColor,
		})
	}
	appzaplog.Debug("[-]RChatTag", zap.Any("resp", resp))
	return resp, nil
}

func (p *ChatTagTarsObj) UChatTag(ctx context.Context, r *bilin.UChatTagReq) (*bilin.UChatTagResp, error) {
	appzaplog.Debug("[+]UChatTag")
	code := int64(MetricCodeSuccess)
	defer func(now time.Time) {
		httpmetrics.DefReport("UChatTag", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	resp := &bilin.UChatTagResp{}
	tag := dao.ChatTag{
		TagName:  r.Chattag.TagName,
		TagColor: r.Chattag.TagColor,
	}
	tag.ID = uint(r.Chattag.Id)
	if err := tag.Update(); err != nil {
		appzaplog.Error("UChatTag Update err", zap.Error(err))
		code = MetricCodeUpdateErr
		return resp, err
	}
	appzaplog.Debug("[-]UChatTag", zap.Any("resp", resp))
	return resp, nil
}

func (p *ChatTagTarsObj) DChatTag(ctx context.Context, r *bilin.DChatTagReq) (*bilin.DChatTagResp, error) {
	appzaplog.Debug("[+]DChatTag")
	code := int64(MetricCodeSuccess)
	defer func(now time.Time) {
		httpmetrics.DefReport("DChatTag", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	resp := &bilin.DChatTagResp{}
	tag := dao.ChatTag{}
	tag.ID = uint(r.Chattag.Id)
	if tag.ID == 0 {
		appzaplog.Warn("delall not allowed")
		code = MetricCodeDelAllWarn
		return resp, fmt.Errorf("delall not allowed")
	}
	if err := tag.Del(); err != nil {
		appzaplog.Error("DChatTag Del err", zap.Error(err))
		code = MetricCodeDelErr
		return resp, err
	}
	appzaplog.Debug("[-]DChatTag", zap.Any("resp", resp))
	return resp, nil
}
