package rank

import (
	"bilin/apigateway/config"
	"bilin/clientcenter"
	"bilin/common/cacheprocessor"
	"bilin/protocol/userinfocenter"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"context"
	"errors"
	"strings"
	"sync/atomic"
	"time"
)

const (
	TOP_NUM = 3

	ANCHOR_RANK_TARGET_URL     = "/bilin18/longlist/index.html#/index?id=0"
	CONTRIBUTE_RANK_TARGET_URL = "/bilin18/longlist/index.html#/index?id=1"
	GUARD_RANK_TARGET_URL      = "/bilin18/longlist/index.html#/index?id=2"

	FIRST_ANCHOR_BADGE_URL  = "https://vipweb.bs2cdn.yy.com/vipinter_c3f522e296c540f5be53a820b58b6ac2.png"
	SECOND_ANCHOR_BADGE_URL = "https://vipweb.bs2cdn.yy.com/vipinter_ee666ea19dff436e9a0852f96ffc09cc.png"
	THIRD_ANCHOR_BADGE_URL  = "https://vipweb.bs2cdn.yy.com/vipinter_421efb4155d64f7ea24e74ec4675cc0f.png"

	FIRST_GUARD_BADGE_URL  = "https://vipweb.bs2cdn.yy.com/vipinter_a28da2e556af47f580ad8d0a6b5be80b.png"
	SECOND_GUARD_BADGE_URL = "https://vipweb.bs2cdn.yy.com/vipinter_1a97578f7e0340cc92cb49a8ece6df1d.png"
	THIRD_GUARD_BADGE_URL  = "https://vipweb.bs2cdn.yy.com/vipinter_beccf61ee73d4f57ad176204ecbf78e2.png"

	FIRST_CONTRIBUTE_BADGE_URL  = "https://vipweb.bs2cdn.yy.com/vipinter_c2da8a8e97db4e3f96dcce6aec4bb039.png"
	SECOND_CONTRIBUTE_BADGE_URL = "https://vipweb.bs2cdn.yy.com/vipinter_87b41092cc3d41c4bea5271f858506b5.png"
	THIRD_CONTRIBUTE_BADGE_URL  = "https://vipweb.bs2cdn.yy.com/vipinter_db219faf00f14ec18d196da05a625fc7.png"

	//已经不需要点击效果
	ANCHOR_RANK_ICON_URL = "https://vipweb.bs2cdn.yy.com/vipinter_19c13173b15b4493a396d4c38355c459.png"
	//ANCHOR_RANK_ONCLICK_ICON_URL     = "https://vipweb.bs2cdn.yy.com/vipinter_661fc7845bd74c53951d9b31841e81c2.png"
	GUARD_RANK_ICON_URL = "https://vipweb.bs2cdn.yy.com/vipinter_08ceefab0a124dfd88a616df2c22cca9.png"
	//GUARD_RANK_ONCLICK_ICON_URL      = "https://vipweb.bs2cdn.yy.com/vipinter_492241f196a945b99742fe37279bf3d1.png"
	CONTRIBUTE_RANK_ICON_URL = "https://vipweb.bs2cdn.yy.com/vipinter_582a6692e1b643cdbaea4e7e3bba9a05.png"
	//CONTRIBUTE_RANK_ONCLICK_ICON_URL = "https://vipweb.bs2cdn.yy.com/vipinter_61fad78a2ddd42e1b5b1afe2696017df.png"
)

type ChanlBillRank struct {
	UID   uint64 `json:"uid"`
	Value int64  `json:"value"`
	Rank  int64  `json:"rank"`
}

type RankUser struct {
	UserID   uint64 `json:"user_id"`
	NickName string `json:"nick_name"`
	Avatar   string `json:"avatar"`
}

type RankInfo struct {
	Users       []*RankUser `json:"users"`
	TargetURL   string      `json:"target_url"`
	Icon        string      `json:"icon"`
	Title       string      `json:"title"`
	FirstBadge  string      `json:"first_badge"`
	SecondBadge string      `json:"second_badge"`
	ThirdBadge  string      `json:"third_badge"`
}

type TadayRank struct {
	ContributeRank *RankInfo `json:"contribute_rank"` //土豪榜
	AnchorRank     *RankInfo `json:"anchor_rank"`     //主播榜
	GuardRank      *RankInfo `json:"guard_rank"`      //守护榜
}

type TadayRankBody struct {
	IsShow    bool      `json:"is_show"`
	TadayRank TadayRank `json:"today_rank"`
}

var ContributeRank atomic.Value
var AnchorRank atomic.Value
var GuardRank atomic.Value
var ChannelBillRank atomic.Value

func InitRank() error {
	if conf := config.GetAppConfig(); conf != nil {
		InitThriftConnentPool(strings.Join(conf.RankThriftAddrs, ";"))
		guard_rank_info := GetDefaultGuardRankInfo()
		GuardRank.Store(guard_rank_info)
		contribute_rank_info := GetDefaultContributeRankInfo()
		ContributeRank.Store(contribute_rank_info)
		anchor_rank_info := GetDefaultAnchorRankInfo()
		AnchorRank.Store(anchor_rank_info)
		bill_rank_map := make(map[uint64]ChanlBillRank)
		ChannelBillRank.Store(&bill_rank_map)

		//不成功不更新 但不影响启动
		cacheprocessor.CacheProcessor("LoadContributeRankInfo", 20*time.Second, LoadContributeRankInfo)
		cacheprocessor.CacheProcessor("LoadAnchorRankInfo", 20*time.Second, LoadAnchorRankInfo)
		cacheprocessor.CacheProcessor("LoadGuardRankInfo", 20*time.Second, LoadGuardRankInfo)
		cacheprocessor.CacheProcessor("LoadChannelBillRankInfo", 20*time.Second, LoadChannelBillRankInfo)
	} else {
		return errors.New("appconfig not init")
	}

	return nil
}

func GetContributeRank() *RankInfo {
	return ContributeRank.Load().(*RankInfo)
}

func GetAnchorRank() *RankInfo {
	return AnchorRank.Load().(*RankInfo)
}

func GetGuardRank() *RankInfo {
	return GuardRank.Load().(*RankInfo)
}

func GetChannelBillRank() *map[uint64]ChanlBillRank {
	return ChannelBillRank.Load().(*map[uint64]ChanlBillRank)
}

func GetRankUserInfo(uids []uint64) ([]*RankUser, error) {
	var list []*RankUser
	if len(uids) == 0 {
		appzaplog.Debug("GetRankUserInfo empty")
		return list, nil
	}

	resp, err := clientcenter.UserInfoClient().GetUserInfo(context.TODO(), &userinfocenter.GetUserInfoReq{uids})
	if err != nil {
		appzaplog.Error("GetUserInfo err", zap.Error(err))
		return nil, err
	}

	if resp.Ret == nil || resp.Ret.Code != userinfocenter.Result_SUCCESS {
		appzaplog.Error("GetUserInfo err")
		return nil, errors.New("GetUserInfo err")
	}

	if resp.Users == nil {
		//appzaplog.Warn("no users", zap.Any("uids", uids))
		//return nil, errors.New("no user")
		appzaplog.Error("no users", zap.Any("uids", uids)) //这种情况不应该存在,如果有,还是确认一下
		return list, nil
	}

	for _, uid := range uids {
		if v, ok := resp.Users[uid]; ok {
			u := RankUser{}
			u.UserID = v.Uid
			u.NickName = v.NickName
			u.Avatar = v.Avatar
			list = append(list, &u)
		}
	}
	appzaplog.Info("GetRankUserInfo", zap.Any("list", list))
	return list, nil
}

func LoadContributeRankInfo() error {
	contribute_info, err := GetContributeRankInfo()
	if err != nil {
		appzaplog.Error("GetContributeRankInfo fail, not reload", zap.Error(err))
		return err
	}
	appzaplog.Debug("GetContributeRankInfo success, reload data", zap.Any("contribute_info", contribute_info))
	ContributeRank.Store(contribute_info)
	return nil
}

func LoadAnchorRankInfo() error {
	anchor_info, err := GetAnchorRankInfo()
	if err != nil {
		appzaplog.Error("GetAnchorRankInfo fail, not reload", zap.Error(err))
		return err
	}
	appzaplog.Debug("GetRankUserInfo success, reload data", zap.Any("anchor_info", anchor_info))
	AnchorRank.Store(anchor_info)
	return nil
}

func LoadGuardRankInfo() error {
	guard_info, err := GetGuardRankInfo()
	if err != nil {
		appzaplog.Error("GetGuardRankInfo fail, not reload", zap.Error(err))
		return err
	}
	appzaplog.Debug("GetGuardRankInfo success, reload data", zap.Any("guard_info", guard_info))
	GuardRank.Store(guard_info)
	return nil
}

func LoadChannelBillRankInfo() error {
	channel_bill_info, err := GetChannelBillRankInfo()
	if err != nil {
		appzaplog.Error("GetChannelBillRankInfo fail, not reload", zap.Error(err))
		return err
	}
	appzaplog.Debug("GetChannelBillRankInfo success, reload data", zap.Any("channel_bill_info", channel_bill_info))
	ChannelBillRank.Store(channel_bill_info)
	return nil
}
