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
	MAX_REQ_UID_NUM = 20
)

type UserInfoCenterObj struct {
}

func NewUserInfoCenterObj() *UserInfoCenterObj {
	return &UserInfoCenterObj{}
}

func (this *UserInfoCenterObj) GetUserInfo(ctx context.Context, req *userinfocenter.GetUserInfoReq) (*userinfocenter.GetUserInfoResp, error) {

	code := int64(0)
	defer func(now time.Time) {
		httpmetrics.DefReport("GetUserInfo", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())

	if req == nil {
		code = ParamInvalid
		log.Error("nill req pointer")
		return nil, errors.New("nill req pointer")
	}

	var (
		resp = &userinfocenter.GetUserInfoResp{
			Ret:   &userinfocenter.Result{},
			Users: make(map[uint64]*userinfocenter.UserInfo),
		}
	)

	if req.Uids == nil || len(req.Uids) == 0 || len(req.Uids) > MAX_REQ_UID_NUM {
		code = ParamInvalid
		log.Warn("uid is empty or to mush")
		resp.Ret.Code = userinfocenter.Result_PARAM_ERR
		resp.Ret.Desc = "uid is empty or to mush"
		return resp, nil
	}

	for _, uid := range req.Uids {
		//优先读缓存
		pb_user, err := dao.GetCacheUserInfo(uid)
		if err != nil { //如果缓存异常 直接读数据库
			code = GetCacheUserInfoFailed
			log.Error("GetCacheUserInfo fail", zap.Error(err))
			httpmetrics.CounterMetric(GetCacheUserInfoFailedKey, 1)
		}

		if err == nil && pb_user != nil {
			log.Debug("GetCacheUserInfo seccess", zap.Uint64("uid", uid), zap.Any("user", *pb_user))
			resp.Users[uid] = pb_user
			httpmetrics.CounterMetric("GetCacheUserInfo success", 1)
			continue //缓存有的直接返回
		}

		//如果找不到或者有异常 这里继续读数据库
		httpmetrics.CounterMetric(GetDBUserInfoKey, 1)
		user, err := GetUserInfo(uid)
		if err != nil {
			code = GetDBUserInfoFailed
			log.Error("get user all info by uid fail", zap.Error(err))
			resp.Ret.Code = userinfocenter.Result_SYSTEM_ERR
			resp.Ret.Desc = "get user info fail"
			httpmetrics.CounterMetric(GetDBUserInfoFailedKey, 1)
			return resp, err
		}

		if user == nil {
			continue
		}

		resp.Users[uid] = user

		//把数据设置到缓存
		if err := dao.SetCacheUserInfo(uid, user); err != nil {
			code = SetCacheUserInfoFailed
			httpmetrics.CounterMetric(SetCacheUserInfoFailedKey, 1)
			log.Error("SetCacheUserInfo fail", zap.Uint64("uid", uid), zap.Error(err), zap.Any("user", user))
		}
	}

	resp.Ret.Code = userinfocenter.Result_SUCCESS
	resp.Ret.Desc = "success"
	log.Debug("get user info success", zap.Any("uids", req.Uids), zap.Any("resp", *resp))

	return resp, nil
}

func (this *UserInfoCenterObj) GetSingleUserInfo(ctx context.Context, r *userinfocenter.GetSingleUserInfoReq) (*userinfocenter.GetSingleUserInfoResp, error) {
	code := int64(0)
	defer func(now time.Time) {
		httpmetrics.DefReport("GetSingleUserInfo", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())

	if r == nil {
		code = ParamInvalid
		log.Error("nill req pointer")
		return nil, errors.New("nill req pointer")
	}

	var (
		resp = &userinfocenter.GetSingleUserInfoResp{
			Uinfo: &userinfocenter.UserInfo{},
		}
	)

	//优先读缓存
	pb_user, err := dao.GetCacheUserInfo(r.Uid)
	if err != nil { //如果缓存异常 直接读数据库
		code = GetCacheUserInfoFailed
		log.Error("GetCacheUserInfo fail", zap.Error(err))
		httpmetrics.CounterMetric(GetCacheUserInfoFailedKey, 1)
	}

	if err == nil && pb_user != nil {
		log.Debug("GetCacheUserInfo seccess", zap.Uint64("uid", r.Uid), zap.Any("user", *pb_user))
		resp.Uinfo = pb_user
		httpmetrics.CounterMetric("GetCacheUserInfo success", 1)
		return resp, nil //缓存有的直接返回
	}

	//如果找不到或者有异常 这里继续读数据库
	httpmetrics.CounterMetric(GetDBUserInfoKey, 1)
	user, err := GetUserInfo(r.Uid)
	if err != nil {
		code = GetDBUserInfoFailed
		log.Error("get user all info by uid fail", zap.Error(err))
		httpmetrics.CounterMetric(GetDBUserInfoFailedKey, 1)
		return resp, err
	}
	resp.Uinfo = user
	if user == nil {
		return resp, nil
	}
	//把数据设置到缓存
	if err := dao.SetCacheUserInfo(r.Uid, user); err != nil {
		code = SetCacheUserInfoFailed
		httpmetrics.CounterMetric(SetCacheUserInfoFailedKey, 1)
		log.Error("SetCacheUserInfo fail", zap.Uint64("uid", r.Uid), zap.Error(err), zap.Any("user", user))
	}

	log.Debug("get user info success", zap.Uint64("uid", r.Uid), zap.Any("resp", *resp))

	return resp, nil
}

func GetUserInfo(uid uint64) (*userinfocenter.UserInfo, error) {
	//不在缓存, 查询mysql
	user, err := dao.GetUserInfo(uid)
	if err != nil {
		httpmetrics.CounterMetric(GetDBUserInfoTableUserFailedKey, 1)
		log.Error("get user info fail", zap.Error(err))
		return nil, err
	}

	if user == nil {
		log.Warn("user not fount", zap.Uint64("uid", uid))
		return nil, nil
	}

	var user_info userinfocenter.UserInfo

	url, err := dao.GetUserAvatarUrl(user.AvatarId, user.UserId)
	if err != nil {
		httpmetrics.CounterMetric(GetDBUserInfoTableAvatarFailedKey, 1)
		log.Error("GetUserAvatarUrl fail", zap.Error(err))
		return nil, err
	}

	user_info.Avatar = url
	exteninfo, err := dao.GetUserExtension(uid)
	if err != nil {
		httpmetrics.CounterMetric(GetDBUserInfoTableExtensionFailedKey, 1)
		log.Error("GetUserExtension failed", zap.Error(err))
		return nil, err
	}

	if exteninfo == nil {
		user_info.City = "火星"
	} else {
		cc, err := dao.GetCommonCountry(exteninfo.CITY)
		if err != nil {
			httpmetrics.CounterMetric(GetDBUserInfoTableCountryFailedKey, 1)
			log.Error("GetCommonCountry failed", zap.Error(err))
			return nil, err
		}

		if cc == nil {
			user_info.City = "火星"
		} else {
			user_info.City = cc.Name
		}
	}

	user_info.Uid = user.UserId
	user_info.NickName = user.NickName
	user_info.Sex = user.Sex
	user_info.Sign = user.Sign
	user_info.Birthday = user.Birthday.Unix()
	user_info.Showsex = user.ShowSex

	return &user_info, nil
}
