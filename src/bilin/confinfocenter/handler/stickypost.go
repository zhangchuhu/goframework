package handler

import (
	"bilin/confinfocenter/cache"
	"bilin/confinfocenter/dao"
	"bilin/protocol"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/kissgo/httpmetrics"
	"context"
	"errors"
	"time"
)

// 品类运营置顶区
func (this *ConfInfoServantObj) CategoryStickie(ctx context.Context, r *bilin.CategoryStickieReq) (*bilin.CategoryStickieResp, error) {
	code := int64(0)
	defer func(now time.Time) {
		httpmetrics.DefReport("CategoryStickie", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())

	resp := &bilin.CategoryStickieResp{
		Categoryinfo: make(map[uint64]*bilin.CategoryStickieInfo),
	}
	info, err := dao.GetStickie()
	if err != nil {
		code = GetStickieFailed
		return resp, err
	}

	for _, v := range info {
		resp.Categoryinfo[uint64(v.RoomID)] = &bilin.CategoryStickieInfo{
			Typeid:    v.TypeId,
			Sort:      v.Weight,
			Roomid:    v.RoomID,
			Starttime: v.StartTime.Unix(),
			Endtime:   v.EndTime.Unix(),
			Id:        int64(v.ID),
		}
	}
	return resp, nil
}

// 用户图标icon
func (this *ConfInfoServantObj) BatchUserBabge(ctx context.Context, r *bilin.UserBabgeReq) (*bilin.UserBabgeResp, error) {
	code := int64(0)
	defer func(now time.Time) {
		httpmetrics.DefReport("BatchUserBabge", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	resp := &bilin.UserBabgeResp{}
	info, err := cache.GetUserBadge()
	if err != nil {
		code = GetUserBadgeFailed
		return nil, err
	}

	for _, v := range info {
		resp.Userbabgeinfo = append(resp.Userbabgeinfo, &bilin.UserBabgeInfo{
			Userid: uint64(v.UserID),
			Url:    v.BadgeUrl,
		})
	}
	return resp, nil
}

func (this *ConfInfoServantObj) CreateCategoryStickie(ctx context.Context, r *bilin.CreateCategoryStickieReq) (*bilin.CreateCategoryStickieResp, error) {
	code := int64(0)
	defer func(now time.Time) {
		httpmetrics.DefReport("CreateCategoryStickie", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	resp := &bilin.CreateCategoryStickieResp{}
	sticky := dao.Stickie{
		TypeId:    r.Info.Typeid,
		RoomID:    r.Info.Roomid,
		StartTime: time.Unix(r.Info.Starttime, 0),
		EndTime:   time.Unix(r.Info.Endtime, 0),
	}
	err := sticky.Create()
	if err != nil {
		code = CreateStickieFailed
		return resp, err
	}
	return resp, nil
}
func (this *ConfInfoServantObj) UpdateCategoryStickie(ctx context.Context, r *bilin.UpdateCategoryStickieReq) (*bilin.UpdateCategoryStickieResp, error) {
	code := int64(0)
	defer func(now time.Time) {
		httpmetrics.DefReport("UpdateCategoryStickie", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	resp := &bilin.UpdateCategoryStickieResp{}
	sticky := dao.Stickie{
		TypeId:    r.Info.Typeid,
		RoomID:    r.Info.Roomid,
		StartTime: time.Unix(r.Info.Starttime, 0),
		EndTime:   time.Unix(r.Info.Endtime, 0),
	}
	sticky.ID = uint(r.Info.Id)
	err := sticky.Update()
	if err != nil {
		code = UpdateStickieFailed
		return resp, err
	}
	return resp, nil
}
func (this *ConfInfoServantObj) DelCategoryStickie(ctx context.Context, r *bilin.DelCategoryStickieReq) (*bilin.DelCategoryStickieResp, error) {
	code := int64(0)
	defer func(now time.Time) {
		httpmetrics.DefReport("DelCategoryStickie", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	resp := &bilin.DelCategoryStickieResp{}
	if r.Info.Id == 0 {
		appzaplog.Warn("DelCategoryStickie all", zap.Any("req", r))
		return resp, errors.New("DelCategoryStickie all")
	}
	sticky := dao.Stickie{}
	sticky.ID = uint(r.Info.Id)
	err := sticky.Del()
	if err != nil {
		code = DelStickieFailed
		return resp, err
	}
	return resp, nil
}

// 返回所有可用的置顶信息
func (this *ConfInfoServantObj) AvailableCategoryStickie(ctx context.Context, r *bilin.AvailableCategoryStickieReq) (*bilin.AvailableCategoryStickieResp, error) {
	code := int64(0)
	defer func(now time.Time) {
		httpmetrics.DefReport("AvailableCategoryStickie", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())

	resp := &bilin.AvailableCategoryStickieResp{}
	info, err := dao.GetStickie()
	if err != nil {
		code = GetStickieFailed
		return nil, err
	}

	for _, v := range info {
		resp.Infos = append(resp.Infos, &bilin.CategoryStickieInfo{
			Typeid:    v.TypeId,
			Sort:      v.Weight,
			Roomid:    v.RoomID,
			Starttime: v.StartTime.Unix(),
			Endtime:   v.EndTime.Unix(),
			Id:        int64(v.ID),
		})
	}
	return resp, nil
}
