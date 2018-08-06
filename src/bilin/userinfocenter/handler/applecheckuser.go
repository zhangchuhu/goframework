package handler

import (
	"bilin/protocol/userinfocenter"
	"bilin/userinfocenter/dao"
	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/kissgo/httpmetrics"
	"context"
	"errors"
	"time"
)

const (
	AppleCheckStatus = 0
)

func (this *UserInfoCenterObj) IsAppleCheckUser(ctx context.Context, req *userinfocenter.IsAppleCheckUserReq) (*userinfocenter.IsAppleCheckUserResp, error) {

	code := int64(0)
	defer func(now time.Time) {
		httpmetrics.DefReport("IsAppleCheckUser", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())

	if req == nil {
		code = ParamInvalid
		log.Error("nill req pointer")
		return nil, errors.New("nill req pointer")
	}

	var (
		resp = &userinfocenter.IsAppleCheckUserResp{
			Uid:            0,
			Applecheckuser: false,
		}
	)

	if req.Uid == 0 {
		log.Warn("IsAppleCheckUser:invalid req uid")
		return resp, nil
	}

	resp.Uid = req.Uid

	if status, err := dao.GetCacheOpenStatus(req.Uid, req.Version, req.Clienttype, req.Ip); err != nil {
		code = GetCacheOpenStatusFailed
		httpmetrics.CounterMetric(GetCacheOpenStatusFailedKey, 1)
		log.Error("GetUserOpenStaus error", zap.Any("req", req), zap.Error(err))
		// 这里不返回 继续rpc查询
	} else {
		if status != dao.NOT_FOUND { //如果找不到 rpc查询
			if status == AppleCheckStatus {
				resp.Applecheckuser = true //默认false
			}
			log.Debug("IsAppleCheckUser success", zap.Any("req", req), zap.Any("resp", *resp))
			httpmetrics.CounterMetric(GetCacheOpenStatusKey, 1)
			return resp, nil
		}
	}

	status, err := dao.GetUserOpenStaus(req.Uid, req.Version, req.Clienttype, req.Ip)
	if err != nil {
		code = GetThriftOpenStatusFailed
		log.Error("GetUserOpenStaus error", zap.Any("req", req), zap.Error(err))
		httpmetrics.CounterMetric(GetThriftOpenStatusFailedKey, 1)
		// 异常直接返回false
		return resp, nil
	}

	httpmetrics.CounterMetric(GetThriftOpenStatusKey, 1)

	if status == AppleCheckStatus {
		resp.Applecheckuser = true
	}

	if err := dao.SetCacheOpenStatus(req.Uid, status, req.Version, req.Clienttype, req.Ip); err != nil {
		code = SetCacheOpenStatusFailed
		log.Error("SetCacheOpenStatus error", zap.Any("req", req), zap.Error(err))
		httpmetrics.CounterMetric(SetCacheOpenStatusFailedKey, 1)
	}

	log.Debug("IsAppleCheckUser success", zap.Any("req", req), zap.Any("resp", *resp))

	return resp, nil
}
