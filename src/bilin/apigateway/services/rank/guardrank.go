package rank

import (
	"bilin/apigateway/config"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

const (
//GUARD_RANK_URL = "http://bilin.yy.com/guardRankingList.do"
)

const DIAL_TIMEOUT = time.Duration(5 * time.Second)
const HEAD_TIMEOUT = time.Duration(5 * time.Second)

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, DIAL_TIMEOUT)
}

type GuardRankResp struct {
	IsEncrypt string        `json:"isEncrypt"`
	Data      GuardRankData `json:"data"`
}

type GuardRankData struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Body []GuardRankUser `json:"body"`
}

type GuardRankUser struct {
	NickName       string `json:"nickName"`
	UID            uint64 `json:"uid"`
	Header         string `json:"header"`
	Sign           string `json:"sign"`
	AttentionCount int64  `json:"attentionCount"`
	Status         int64  `json:"status"`
	GuardCount     int64  `json:"guardCount"`
	LiveID         int64  `json:"liveId"`
}

func GetGuardRankList() (error, []GuardRankUser) {

	tr := &http.Transport{
		//使用带超时的连接函数
		Dial: dialTimeout,
		//建立连接后读超时
		ResponseHeaderTimeout: HEAD_TIMEOUT,
	}

	client := http.Client{
		Transport: tr,
		//总超时，包含连接读写
		Timeout: DIAL_TIMEOUT,
	}

	req, err := http.NewRequest("GET", config.GetAppConfig().GuardRankUrl, nil)
	if err != nil {
		appzaplog.Error("NewRequest err", zap.Error(err))
		return err, nil
	}

	resp, err := client.Do(req)
	if err != nil {
		appzaplog.Error("client.Do err", zap.Error(err))
		return err, nil
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		appzaplog.Error("ReadAll err", zap.Error(err))
		return err, nil
	}

	//Log.Info("resp status: %s,resp body: %s", resp.Status, string(respBody))

	if resp.StatusCode != http.StatusOK {
		errMsg := "resp status:" + resp.Status + ", resp body" + string(respBody)
		appzaplog.Error("status error", zap.String("http", errMsg))
		return errors.New(errMsg), nil
	}

	var response GuardRankResp

	err = json.Unmarshal(respBody, &response)
	if err != nil {
		errMsg := "Unmarshal fail:" + err.Error()
		appzaplog.Error("json error", zap.String("json", errMsg))
		return errors.New(errMsg), nil
	}

	if response.Data.Code != 0 {
		return errors.New(response.Data.Msg), nil
	}

	appzaplog.Debug("GetGuardRankList resp info", zap.Any("data", response.Data))
	return nil, response.Data.Body
}

func GetGuardRankInfo() (*RankInfo, error) {

	rank_info := GetDefaultGuardRankInfo()
	err, guard_rank_list := GetGuardRankList()
	if err != nil {
		appzaplog.Error("GetGuardRankList err", zap.Error(err))
		return nil, err
	}

	for i := 0; i < len(guard_rank_list) && i < TOP_NUM; i++ {
		user := RankUser{}
		user.UserID = guard_rank_list[i].UID
		user.NickName = guard_rank_list[i].NickName
		user.Avatar = guard_rank_list[i].Header
		rank_info.Users = append(rank_info.Users, &user)
	}

	appzaplog.Debug("GetGuardRankInfo resp ", zap.Any("rankinfo", rank_info))
	return rank_info, nil
}

func GetDefaultGuardRankInfo() *RankInfo {
	rank_info := &RankInfo{}
	//rank_info.TargetURL = "http://" + config.GetAppConfig().RankTargetHost + GUARD_RANK_TARGET_URL
	rank_info.TargetURL = config.GetAppConfig().GuardRankTargetURL
	rank_info.Title = "今日守护榜"
	rank_info.Icon = GUARD_RANK_ICON_URL
	rank_info.FirstBadge = FIRST_GUARD_BADGE_URL
	rank_info.SecondBadge = SECOND_GUARD_BADGE_URL
	rank_info.ThirdBadge = THIRD_GUARD_BADGE_URL
	return rank_info
}
