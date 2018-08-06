/*
 * Copyright (c) 2018-07-10.
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
	"sort"
	"strconv"
	"strings"
	"time"
)

// 用户聊天标签CRUD
func (p *ChatTagTarsObj) CUserChatTag(ctx context.Context, r *bilin.CUserChatTagReq) (*bilin.CUserChatTagResp, error) {
	appzaplog.Debug("[+]CUserChatTag")
	code := int64(MetricCodeSuccess)
	defer func(now time.Time) {
		httpmetrics.DefReport("CUserChatTag", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	resp := &bilin.CUserChatTagResp{}
	tag := dao.UserChatTag{
		FromUserID:  r.Info.Fromuserid,
		ToUserID:    r.Info.Touserid,
		ChatTags:    r.Info.Chattags,
		UpdateTimes: r.Info.Updatetimes,
		TalkSecond:  r.Info.Talksecond,
		TagStatus:   r.Info.Tagstatus,
	}
	if err := tag.Create(); err != nil {
		appzaplog.Error("CUserChatTag Create err", zap.Error(err))
		code = MetricCodeCreateErr
		return resp, err
	}
	appzaplog.Debug("[-]CUserChatTag", zap.Any("resp", resp))
	return resp, nil
}
func (p *ChatTagTarsObj) RUserChatTag(ctx context.Context, r *bilin.RUserChatTagReq) (*bilin.RUserChatTagResp, error) {
	appzaplog.Debug("[+]RUserChatTag")
	code := int64(MetricCodeSuccess)
	defer func(now time.Time) {
		httpmetrics.DefReport("RUserChatTag", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	resp := &bilin.RUserChatTagResp{}
	if r.Info == nil || r.Info.Touserid == 0 {
		appzaplog.Warn("[-]RUserChatTag no userid")
		return resp, nil
	}
	uerchattag := dao.UserChatTag{
		FromUserID:  r.Info.Fromuserid,
		ToUserID:    r.Info.Touserid,
		ChatTags:    r.Info.Chattags,
		UpdateTimes: r.Info.Updatetimes,
		TalkSecond:  r.Info.Talksecond,
		TagStatus:   r.Info.Tagstatus,
	}
	info, err := uerchattag.GetAll()
	if err != nil {
		appzaplog.Error("RUserChatTag dao.GetAll err", zap.Any("info", r.Info), zap.Error(err))
		code = MetricCodeReadErr
		return resp, err
	}
	for _, v := range info {
		resp.Info = append(resp.Info, &bilin.UserChatTag{
			Id:          int64(v.ID),
			Fromuserid:  v.FromUserID,
			Touserid:    v.ToUserID,
			Chattags:    v.ChatTags,
			Updatetimes: v.UpdateTimes,
			Talksecond:  v.TalkSecond,
			Tagstatus:   v.TagStatus,
		})
	}
	appzaplog.Debug("[-]RUserChatTag", zap.Any("info", r.Info), zap.Any("resp", resp))
	return resp, nil
}

func convertTagId(tagid string) (tagids []int64) {
	if tagid == "" {
		return []int64{}
	}
	tagids_ := strings.Split(tagid, ",")
	for _, v := range tagids_ {
		if tagIdInt, err := strconv.ParseInt(v, 10, 64); err == nil {
			tagids = append(tagids, tagIdInt)
		}
	}
	return
}

// 根据标签个数排序的前N个标签信息
func (p *ChatTagTarsObj) RTopNUserChatTagSummary(ctx context.Context, r *bilin.RTopNUserChatTagSummaryReq) (*bilin.RTopNUserChatTagSummaryResp, error) {
	appzaplog.Debug("[+]RTopNUserChatTag")
	code := int64(MetricCodeSuccess)
	defer func(now time.Time) {
		httpmetrics.DefReport("RTopNUserChatTagSummary", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	resp := &bilin.RTopNUserChatTagSummaryResp{}
	if r.Topuser == nil {
		appzaplog.Warn("RTopNUserChatTagSummary no Topuser")
		return resp, nil
	}
	info, err := userChatTag(r.Topuser.Touserid)
	if err != nil {
		appzaplog.Error("RTopNUserChatTag userChatTag err", zap.Error(err), zap.Int64("uid", r.Topuser.Touserid))
		code = MetricCodeReadErr
		return resp, err
	}
	resp.Summary = &bilin.UserChatTagSummaryS{
		userChatTagSummary(info, cache.TakeChatTagCache(), r.Topuser.Topn),
	}

	appzaplog.Debug("[-]RTopNUserChatTag", zap.Any("resp", resp))
	return resp, nil
}

const maxbatchuchattagnum = 10

func (p *ChatTagTarsObj) BatchRTopNUserChatTagSummary(ctx context.Context, r *bilin.BatchRTopNUserChatTagSummaryReq) (*bilin.BatchRTopNUserChatTagSummaryResp, error) {
	appzaplog.Debug("[+]BatchRTopNUserChatTagSummary")
	code := int64(MetricCodeSuccess)
	defer func(now time.Time) {
		httpmetrics.DefReport("BatchRTopNUserChatTagSummary", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	resp := &bilin.BatchRTopNUserChatTagSummaryResp{
		Summary: make(map[int64]*bilin.UserChatTagSummaryS),
	}
	if len(r.Topnuser) > maxbatchuchattagnum {
		appzaplog.Error("BatchRTopNUserChatTagSummary too many elements one time", zap.Int("threshold", maxbatchuchattagnum))
		return resp, fmt.Errorf("over %d user one time", maxbatchuchattagnum)
	}

	cachetag := cache.TakeChatTagCache()
	for _, v := range r.Topnuser {
		info, err := userChatTag(v.Touserid)
		if err != nil {
			appzaplog.Error("BatchRTopNUserChatTagSummary userChatTag err", zap.Error(err), zap.Int64("uid", v.Touserid))
			code = MetricCodeReadErr
			return resp, err
		}
		resp.Summary[v.Touserid] = &bilin.UserChatTagSummaryS{
			userChatTagSummary(info, cachetag, v.Topn),
		}
	}
	appzaplog.Debug("[-]BatchRTopNUserChatTagSummary", zap.Any("resp", resp))
	return resp, nil
}

func userChatTag(userId int64) ([]dao.UserChatTag, error) {
	if userId == 0 {
		return nil, fmt.Errorf("Touserid not set")
	}
	uctag := dao.UserChatTag{
		ToUserID: userId,
	}
	return uctag.GetAll()
}

func userChatTagSummary(info []dao.UserChatTag, chattag map[int64]dao.ChatTag, topn int64) []*bilin.UserChatTagSummary {
	tagmap := make(map[int64]*bilin.UserChatTagSummary)
	var summary []*bilin.UserChatTagSummary

	for _, v := range info {
		tagids := convertTagId(v.ChatTags)
		for _, tagid := range tagids {
			if summary, ok := tagmap[tagid]; ok {
				summary.Totaltagnum++
			} else {
				tagmap[tagid] = &bilin.UserChatTagSummary{
					Touserid:    v.ToUserID,
					Tagid:       tagid,
					Totaltagnum: 1,
				}
			}
		}

	}

	for k, v := range tagmap {
		if chattaginfo, ok := chattag[k]; ok {
			v.Tagcolor = chattaginfo.TagColor
			v.Tagname = chattaginfo.TagName
		}
		summary = append(summary, v)
	}
	if len(summary) == 0 {
		return nil
	}

	sort.SliceStable(summary, func(i, j int) bool {
		return summary[i].Totaltagnum > summary[j].Totaltagnum
	})

	if topn > int64(len(summary)) {
		topn = int64(len(summary))
	}
	summary = summary[0:topn]
	return summary
}

func (p *ChatTagTarsObj) UUserChatTag(ctx context.Context, r *bilin.UUserChatTagReq) (*bilin.UUserChatTagResp, error) {
	appzaplog.Debug("[+]UUserChatTag")
	code := int64(MetricCodeSuccess)
	defer func(now time.Time) {
		httpmetrics.DefReport("UUserChatTag", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	resp := &bilin.UUserChatTagResp{}
	tag := dao.UserChatTag{
		FromUserID:  r.Info.Fromuserid,
		ToUserID:    r.Info.Touserid,
		ChatTags:    r.Info.Chattags,
		UpdateTimes: r.Info.Updatetimes,
		TalkSecond:  r.Info.Talksecond,
		TagStatus:   r.Info.Tagstatus,
	}
	tag.ID = uint(r.Info.Id)
	if err := tag.Update(); err != nil {
		appzaplog.Error("UUserChatTag Update err", zap.Error(err))
		code = MetricCodeUpdateErr
		return resp, err
	}
	appzaplog.Debug("[-]UUserChatTag", zap.Any("resp", resp))
	return resp, nil
}

func (p *ChatTagTarsObj) DUserChatTag(ctx context.Context, r *bilin.DUserChatTagReq) (*bilin.DUserChatTagResp, error) {
	appzaplog.Debug("[+]DUserChatTag")
	code := int64(MetricCodeSuccess)
	defer func(now time.Time) {
		httpmetrics.DefReport("DUserChatTag", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	resp := &bilin.DUserChatTagResp{}
	tag := dao.UserChatTag{}
	tag.ID = uint(r.Info.Id)
	if tag.ID == 0 {
		appzaplog.Warn("delall not allowed")
		code = MetricCodeDelAllWarn
		return resp, fmt.Errorf("delall not allowed")
	}
	if err := tag.Del(); err != nil {
		appzaplog.Error("DUserChatTag Del err", zap.Error(err))
		code = MetricCodeDelErr
		return resp, err
	}
	appzaplog.Debug("[-]DUserChatTag", zap.Any("resp", resp))
	return resp, nil
}
