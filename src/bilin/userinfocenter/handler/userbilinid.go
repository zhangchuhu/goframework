package handler

import (
	"bilin/protocol/userinfocenter"
	"bilin/userinfocenter/dao"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/kissgo/httpmetrics"
	"context"
	"time"
)

// 根据比邻uid查询比邻号
func (this *UserInfoCenterObj) BatchUserBiLinId(ctx context.Context, r *userinfocenter.BatchUserBiLinIdReq) (*userinfocenter.BatchUserBiLinIdResp, error) {
	code := int64(0)
	defer func(now time.Time) {
		httpmetrics.DefReport("BatchUserBiLinId", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	resp := &userinfocenter.BatchUserBiLinIdResp{
		Uid2Bilinid: make(map[uint64]uint64),
	}
	bilinidmap, err := dao.BatchUserBiID(r.Uid)
	if err != nil {
		appzaplog.Error("BatchUserBiLinId BatchUserBiID err", zap.Error(err), zap.Uint64s("uids", r.Uid))
		code = BatchUserBLNumFailed
		return resp, err
	}
	for _, v := range bilinidmap {
		resp.Uid2Bilinid[v.USER_ID] = v.BILIN_ID
	}
	return resp, nil
}

// 根据比邻号查询uid
func (this *UserInfoCenterObj) BatchUserIdByBiLinId(ctx context.Context, r *userinfocenter.BatchUserIdByBiLinIdReq) (*userinfocenter.BatchUserIdByBiLinIdResp, error) {
	code := int64(0)
	defer func(now time.Time) {
		httpmetrics.DefReport("BatchUserIdByBiLinId", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	resp := &userinfocenter.BatchUserIdByBiLinIdResp{
		Bilinid2Uid: make(map[uint64]uint64),
	}
	bilinidmap, err := dao.BatchUserID(r.Bilinid)
	if err != nil {
		appzaplog.Error("BatchUserIdByBiLinId BatchUserID err", zap.Error(err), zap.Uint64s("uids", r.Bilinid))
		code = BatchUserBLNumFailed
		return resp, err
	}
	for _, v := range bilinidmap {
		resp.Bilinid2Uid[v.BILIN_ID] = v.USER_ID
	}
	return resp, nil
}
