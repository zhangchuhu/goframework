package handler

import (
	"bilin/confinfocenter/dao"
	"bilin/protocol"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/kissgo/httpmetrics"
	"context"
	"errors"
	"time"
)

// 工会品类推荐
func (this *ConfInfoServantObj) CategoryGuildRec(ctx context.Context, r *bilin.CategoryGuildRecReq) (*bilin.CategoryGuildRecResp, error) {
	code := int64(0)
	defer func(now time.Time) {
		httpmetrics.DefReport("CategoryGuildRec", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())

	rec, err := dao.GetGuildRec()
	if err != nil {
		code = GetGuildRecFailed
		return nil, err
	}
	resp := &bilin.CategoryGuildRecResp{}
	for _, v := range rec {
		resp.Cateogryguildinfos = append(resp.Cateogryguildinfos, &bilin.CategoryGuildRecInfo{
			Roomid: v.RoomID,
			Typeid: v.TypeId,
			Id:     uint64(v.ID),
		})
	}
	return resp, nil
}

func (this *ConfInfoServantObj) UpdateCategoryGuildRec(ctx context.Context, r *bilin.UpdateCategoryGuildRecReq) (*bilin.UpdateCategoryGuildRecResp, error) {
	code := int64(0)
	defer func(now time.Time) {
		httpmetrics.DefReport("UpdateCategoryGuildRec", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	resp := &bilin.UpdateCategoryGuildRecResp{}
	uguild := &dao.GuildRec{
		RoomID: r.Info.Roomid,
		TypeId: r.Info.Typeid,
	}
	uguild.ID = uint(r.Info.Id)

	err := dao.UpdateGuildRec(uguild)
	if err != nil {
		code = UpdateGuildRecFailed
		appzaplog.Error("UpdateGuildRec err", zap.Error(err), zap.Any("req", r))
		return resp, err
	}
	resp.Info = &bilin.CategoryGuildRecInfo{
		Id:     uint64(uguild.ID),
		Roomid: uguild.RoomID,
		Typeid: uguild.TypeId,
	}

	appzaplog.Debug("UpdateCategoryGuildRec", zap.Any("resp", resp))
	return resp, nil
}

func (this *ConfInfoServantObj) DelCategoryGuildRec(ctx context.Context, r *bilin.DelCategoryGuildRecReq) (*bilin.DelCategoryGuildRecResp, error) {
	code := int64(0)
	defer func(now time.Time) {
		httpmetrics.DefReport("DelCategoryGuildRec", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	resp := &bilin.DelCategoryGuildRecResp{}
	if r.Info.Id == 0 {
		appzaplog.Warn("DelCategoryGuildRec all", zap.Any("req", r))
		return resp, errors.New("DelCategoryGuildRec all")
	}
	err := dao.DelGuildRec(uint(r.Info.Id))
	if err != nil {
		code = DelGuildRecFailed
		appzaplog.Error("DelCategoryGuildRec err", zap.Error(err), zap.Any("req", r))
		return resp, err
	}
	return resp, nil
}

func (this *ConfInfoServantObj) CreateCategoryGuildRec(ctx context.Context, r *bilin.CreateCategoryGuildRecReq) (*bilin.CreateCategoryGuildRecResp, error) {
	code := int64(0)
	defer func(now time.Time) {
		httpmetrics.DefReport("CreateCategoryGuildRec", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	resp := &bilin.CreateCategoryGuildRecResp{}
	uguild := &dao.GuildRec{
		RoomID: r.Info.Roomid,
		TypeId: r.Info.Typeid,
	}
	err := uguild.Create()
	if err != nil {
		code = CreateGuildRecFailed
		appzaplog.Error("CreateCategoryGuildRec err", zap.Error(err), zap.Any("req", r))
		return resp, err
	}
	appzaplog.Debug("CreateCategoryGuildRec")
	return resp, nil
}
