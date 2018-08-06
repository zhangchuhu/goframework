package handler

import (
	"bilin/common/onlinepush"
	"bilin/common/onlinequery"
	"bilin/protocol"
	u "bilin/protocol/userinfocenter"

	"context"
	"encoding/json"
	"fmt"
	"time"

	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/tars"
	"github.com/golang/protobuf/proto"
)

type MatchCallBackHanlder struct {
}

var (
	emptyUserInfo = &u.UserInfo{
		Birthday: time.Now().Unix() - 18*365*24*60*60,
	}
)

func unicast(uids []int64, msg proto.Message, maxtype bilin.MaxType, mintype bilin.MinType_MATCH) (offline []int64, err error) {
	const prefix = "unicast "

	var body bilin.CommonMessageBody
	body.Type = int32(mintype)
	if body.Data, err = proto.Marshal(msg); err != nil {
		log.Error(prefix+"[-]proto.Marshal failed", zap.Any("err", err))
		return nil, err
	}
	pushBody, err := proto.Marshal(&body)
	if err != nil {
		log.Error(prefix+"[-]proto.Marshal failed", zap.Any("err", err))
		return nil, err
	}

	multiMsg := bilin.MultiPush{
		Msg: &bilin.ServerPush{
			MessageType: int32(maxtype),
			PushBuffer:  pushBody,
		},
		UserIDs: uids,
	}

	offline, err = onlinepush.PushToUser(multiMsg)
	if err != nil {
		log.Error("[-]PushNotifyToUser failed push", zap.Any("err", err))
		return nil, err
	}

	log.Info("[+]PushNotifyToUser success push",
		zap.Any("uids", uids),
		zap.Any("offline", offline),
		zap.Any("msg", msg.String()),
		zap.Any("maxtype", bilin.MaxType_name[int32(maxtype)]),
		zap.Any("mintype", bilin.MinType_MATCH_name[int32(mintype)]))
	return offline, nil
}

func getUserInfo(u *u.UserInfo) *u.UserInfo {
	if u != nil {
		return u
	} else {
		return emptyUserInfo
	}
}

func GetMatchUserInfo(user UserItem, online int32) *bilin.MatchUserInfo {

	userInfoReq := &u.GetUserInfoReq{}
	userInfoReq.Uids = append(userInfoReq.Uids, uint64(user.Uid))

	comm := tars.NewCommunicator()
	objName := fmt.Sprintf("bilin.userinfocenter.UserInfoCenterObj")
	client := u.NewUserInfoCenterObjClient(objName, comm)
	resp, err := client.GetUserInfo(context.TODO(), userInfoReq)
	if err != nil {
		log.Error("GetMatchUserInfo client.GetUserInfo", zap.Error(err), zap.Any("service", objName))
		resp = nil
	}

	rid, err := onlinequery.GetUserRoom(int64(user.Uid))
	if err != nil {
		log.Error("GetMatchUserInfo onlinequery.GetUserRoom", zap.Error(err), zap.Any("user", user))
	}
	var userOnline int32 = 0
	if rid >= -1 {
		userOnline = 1
	}
	if online != -1 {
		userOnline = online
	}

	var nick string
	var avatar string
	var postion string
	var age int64

	if resp != nil {
		if len(resp.Users) > 0 {
			info := getUserInfo(resp.Users[uint64(user.Uid)])
			nick = info.NickName
			avatar = info.Avatar
			postion = info.City
			age = (time.Now().Unix() - info.Birthday) / (365 * 24 * 60 * 60)
		} else {
			log.Error("GetMatchUserInfo client.GetUserInfo returns zero user", zap.Any("resp", resp), zap.Any("service", objName))
		}
	}

	userinfo := bilin.MatchUserInfo{
		Uid:      user.Uid,
		Sex:      int32(user.Sex),
		Postion:  postion,
		Nick:     nick,
		Avatar:   avatar,
		Isonline: userOnline,
		Age:      int32(age),
	}
	return &userinfo
}

func (m *MatchCallBackHanlder) ReturnMatchOk(user UserItem, matchList UserList) {
	if len(matchList) == 0 {
		log.Error("ReturnMatchOk but matchList is empty", zap.Any("user", user))
		return
	}

	var robotList UserList
	if len(matchList) < 3 {
		robot := GetRobotList()
		n := 0
		for i := len(matchList); i < 3; i++ {
			if len(robot) > n {
				robotList = append(robotList, robot[n])
			} else {
				dummy := UserItem{uint32(n + 1), 0, 0, 0, "haiwai", 1500000000000}
				log.Warn("ReturnMatchOk not enough robot, use a hard code one", zap.Any("robot", dummy))
				robotList = append(robotList, dummy)
			}
			n++
		}
	}

	userinfo := GetMatchUserInfo(user, -1)

	var uids []int64

	bc := bilin.OptionalMatchingResult{}
	bc.Matchid = GenerateMatchid()
	bc.Attendees = append(bc.Attendees, userinfo)
	PlayerSelect(int64(user.Uid))

	for _, user := range matchList {
		userinfo := GetMatchUserInfo(user, -1)
		bc.Attendees = append(bc.Attendees, userinfo)
		PlayerSelect(int64(user.Uid))
	}

	for _, user := range robotList {
		userinfo := GetMatchUserInfo(user, 0)
		bc.Attendees = append(bc.Attendees, userinfo)
	}

	// 排除离线
	for i, x := range bc.Attendees {
		if i == 0 {
			uids = append(uids, int64(x.Uid))
			continue // 第一个是总是女性
		}
		if x.Isonline == 1 {
			uids = append(uids, int64(x.Uid))
		}
	}
	if len(uids) <= 1 || bc.Attendees[0].Isonline == 0 {
		log.Warn("ReturnMatchOk skip",
			zap.Any("user", user),
			zap.Any("matchList", matchList),
			zap.Any("robotList", robotList))
		return
	}

	// 首先匹配到的男性，都是在头部。这里把他换到中间。
	if len(bc.Attendees) >= 3 {
		bc.Attendees[1], bc.Attendees[2] = bc.Attendees[2], bc.Attendees[1]
	}

	val, err := json.Marshal(bc)
	if err != nil {
		log.Error("ReturnMatchOk", zap.Error(err))
		return
	}
	AddMatchIdValue(bc.Matchid, string(val))

	if err := fillUserChatTag(&bc); err != nil {
		log.Error("ReturnMatchOk fillUserChatTag err", zap.Error(err))
	}

	log.Info("ReturnMatchOk broadcast OptionalMatchingResult", zap.Any("OptionalMatchingResult", string(val)))

	unicast(uids, &bc, bilin.MaxType_MATCH_MSG, bilin.MinType_MATCH_OPTIONALMATCHINGRESULT_MINTYPE)
}

func (m *MatchCallBackHanlder) ReturnMatchNoSexOk(matchList UserList) {

	var uids []int64
	bc := bilin.OptionalMatchingResult{}
	bc.Matchid = GenerateMatchid()

	for _, user := range matchList {
		userinfo := GetMatchUserInfo(user, -1)
		bc.Attendees = append(bc.Attendees, userinfo)
		PlayerSelect(int64(user.Uid))
	}

	// 排除离线
	for _, x := range bc.Attendees {
		if x.Isonline == 1 {
			uids = append(uids, int64(x.Uid))
		}
	}
	if len(uids) <= 1 {
		log.Warn("ReturnMatchNoSexOk skip",
			zap.Any("matchList", matchList))
		return
	}

	val, err := json.Marshal(bc)
	if err != nil {
		log.Error("ReturnMatchNoSexOk", zap.Error(err))
		return
	}
	AddMatchIdValue(bc.Matchid, string(val))

	log.Info("ReturnMatchNoSexOk broadcast OptionalMatchingResult", zap.Any("OptionalMatchingResult", string(val)))

	unicast(uids, &bc, bilin.MaxType_MATCH_MSG, bilin.MinType_MATCH_OPTIONALMATCHINGRESULT_MINTYPE)
}
