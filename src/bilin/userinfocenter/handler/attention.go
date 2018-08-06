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

// 用户关注，粉丝，魅力值
func (this *UserInfoCenterObj) AttentionInfo(ctx context.Context, r *userinfocenter.AttentionInfoReq) (*userinfocenter.AttentionInfoResp, error) {
	code := int64(0)
	defer func(now time.Time) {
		httpmetrics.DefReport("AttentionInfo", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	resp := &userinfocenter.AttentionInfoResp{}
	attentionme, err := dao.AttentionMeNum(r.Uid)
	if err != nil {
		appzaplog.Error("AttentionInfo AttentionMeNum err", zap.Uint64("uid", r.Uid), zap.Error(err))
		code = AttentionMeFailed
		return resp, err
	}
	fansnum, err := dao.MyAttentionNum(r.Uid)
	if err != nil {
		appzaplog.Error("AttentionInfo MyAttentionNum err", zap.Uint64("uid", r.Uid), zap.Error(err))
		code = AttentionMeFailed
		return resp, err
	}

	userinfo, err := dao.GetUserInfo(r.Uid)
	if err != nil {
		appzaplog.Error("AttentionInfo GetUserInfo err", zap.Uint64("uid", r.Uid), zap.Error(err))
		code = GetDBUserInfoFailed
		return resp, err
	}
	resp.Attentionnum = attentionme
	resp.Fansnum = fansnum
	if userinfo != nil {
		resp.Glamour = userinfo.SHOW_GLAMOUR_VALUE
	}
	return resp, nil
}
