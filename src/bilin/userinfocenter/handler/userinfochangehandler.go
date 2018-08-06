package handler

import (
	"bilin/userinfocenter/dao"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/kissgo/httpmetrics"
	"encoding/json"
	"github.com/nsqio/go-nsq"
	"time"
)

type UserInfoChangeEvent struct {
	Uid         uint64 `json:"uid"`
	NewNickName string `json:"new_nick_name"`
}

type UserInfoEventHandler struct {
}

func (p *UserInfoEventHandler) HandleMessage(message *nsq.Message) error {
	code := int64(0)
	defer func(now time.Time) {
		httpmetrics.DefReport("UserInfoChangeEvent", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())

	var user UserInfoChangeEvent
	if err := json.Unmarshal(message.Body, &user); err != nil {
		code = ParamMarshalFailed
		appzaplog.Error("HandleMessage unmarshal failed", zap.Error(err), zap.String("body", string(message.Body)))
		httpmetrics.CounterMetric(UpdateUserInfoMarshalFaileKey, 1)
		// drop it
		return nil
	}

	err := dao.DelelteCacheUserInfo(user.Uid)
	if err != nil {
		code = DelCacheUserInfoFailed
		httpmetrics.CounterMetric(DelCacheUserInfoFailedKey, 1)
		appzaplog.Error("DeleteCacheUserInfo failed", zap.Error(err), zap.String("body", string(message.Body)))
		return err
	}

	user_info, err := GetUserInfo(user.Uid)
	if err != nil {
		code = GetDBUserInfoFailed
		appzaplog.Error("get user all info by uid fail", zap.Error(err))
		httpmetrics.CounterMetric(GetDBUserInfoFailedKey, 1)
		return nil
	}

	if user_info == nil {
		appzaplog.Warn("user not found", zap.Uint64("uid", user.Uid))
		return nil
	}

	//把数据设置到缓存
	if err = dao.SetCacheUserInfo(user.Uid, user_info); err != nil {
		httpmetrics.CounterMetric(SetCacheUserInfoFailedKey, 1)
		code = SetCacheUserInfoFailed
		appzaplog.Error("SetCacheUserInfo failed", zap.Uint64("uid", user.Uid), zap.Error(err), zap.Any("user", user_info))
	}

	appzaplog.Debug("HandleMessage success", zap.Any("uid", user.Uid))

	return nil
}
