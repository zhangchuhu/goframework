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

// 主播品类推荐
func (this *ConfInfoServantObj) CategoryHostRec(ctx context.Context, r *bilin.CategoryHostRecReq) (*bilin.CategoryHostRecResp, error) {
	code := int64(0)
	defer func(now time.Time) {
		httpmetrics.DefReport("CategoryHostRec", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())

	rec, err := dao.GetHostRec()
	if err != nil {
		code = GetHostRecFailed
		return nil, err
	}
	resp := &bilin.CategoryHostRecResp{}
	for _, v := range rec {
		resp.Cateogryinfos = append(resp.Cateogryinfos, &bilin.CategoryHostRecInfo{
			Hostid: v.HostID,
			Typeid: v.TypeId,
			Id:     int64(v.ID),
		})
	}
	return resp, nil
}

func (this *ConfInfoServantObj) CreateCategoryHostRec(ctx context.Context, r *bilin.CreateCategoryHostRecReq) (*bilin.CreateCategoryHostRecResp, error) {
	code := int64(0)
	defer func(now time.Time) {
		httpmetrics.DefReport("CreateCategoryHostRec", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	resp := &bilin.CreateCategoryHostRecResp{}
	rec := dao.HostRec{
		HostID: r.Info.Hostid,
		TypeId: r.Info.Typeid,
	}
	if err := rec.Create(); err != nil {
		code = CreateHostRecFailed
		return resp, err
	}
	return resp, nil
}
func (this *ConfInfoServantObj) UpdateCategoryHostRec(ctx context.Context, r *bilin.UpdateCategoryHostRecReq) (*bilin.UpdateCategoryHostRecResp, error) {
	code := int64(0)
	defer func(now time.Time) {
		httpmetrics.DefReport("UpdateCategoryHostRec", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	resp := &bilin.UpdateCategoryHostRecResp{}
	rec := dao.HostRec{
		HostID: r.Info.Hostid,
		TypeId: r.Info.Typeid,
	}
	rec.ID = uint(r.Info.Id)
	if err := rec.Update(); err != nil {
		code = UpdateHostRecFailed
		return resp, err
	}
	return resp, nil
}
func (this *ConfInfoServantObj) DelCategoryHostRec(ctx context.Context, r *bilin.DelCategoryHostRecReq) (*bilin.DelCategoryHostRecResp, error) {
	code := int64(0)
	defer func(now time.Time) {
		httpmetrics.DefReport("DelCategoryHostRec", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	resp := &bilin.DelCategoryHostRecResp{}

	if r.Info.Id == 0 {
		appzaplog.Warn("DelCategoryHostRec all", zap.Any("req", r))
		return resp, errors.New("DelCategoryHostRec all")
	}

	rec := dao.HostRec{}
	rec.ID = uint(r.Info.Id)
	if err := rec.Del(); err != nil {
		code = DelHostRecFailed
		return resp, err
	}
	return resp, nil
}
