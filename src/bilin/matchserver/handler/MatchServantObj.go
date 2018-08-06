package handler

import (
	"bilin/common/onlinequery"
	"bilin/common/thriftpool"
	"bilin/protocol"
	u "bilin/protocol/userinfocenter"
	"bilin/thrift/gen-go/callrecord"
	"bilin/thrift/gen-go/common"
	"bilin/thrift/gen-go/hotline"
	"bilin/thrift/gen-go/meeting"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"

	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/tars"
)

const (
	ResultMatchSuccess int32 = 0
	ResultMatchFailed  int32 = 1
)

var _ bilin.MatchServantServer = &MatchServantObj{}

type MatchServantObj struct {
}

func NewMatchServantObj() *MatchServantObj {
	return &MatchServantObj{}
}

func GetUid(ctx context.Context) (uint32, error) {
	if m, ok := tars.FromOutgoingContext(ctx); ok {
		uid, err := strconv.ParseInt(m["uid"], 10, 48)
		if err != nil {
			log.Error("GetUid", zap.Error(err))
			return 0, err
		}
		return uint32(uid), nil
	}
	return 0, fmt.Errorf("no uid found in context")
}

func GetRoomid(ctx context.Context) (int64, error) {
	if m, ok := tars.FromOutgoingContext(ctx); ok {
		rid, err := strconv.ParseInt(m["subscribe-room-push"], 10, 48)
		if err != nil {
			log.Error("GetRoomid", zap.Error(err))
			return -1, err
		}
		return rid, nil
	}
	return -1, fmt.Errorf("no roomid found in context")
}

func (this *MatchServantObj) MatchRandomCall(ctx context.Context, req *bilin.MatchRandomCallReq) (resp *bilin.MatchRandomCallResp, rpcErr error) {
	const prefix = "MatchRandomCall "
	resp = &bilin.MatchRandomCallResp{Result: ResultMatchSuccess, ErrorDesc: ""}

	uid, err := GetUid(ctx)
	if err != nil {
		resp.Result = ResultMatchFailed
		resp.ErrorDesc = err.Error()
		log.Error(prefix+"GetUid fail", zap.Error(err))
		return
	}

	log.Info(prefix+"begin", zap.Any("uid", uid), zap.Any("req", req.String()))

	var sex int
	if req.Sex == Female {
		sex = SexFemale
	} else {
		sex = SexMale
	}
	var playtype int
	if req.MatchType == MatchSex {
		playtype = CallTypeHetero
	} else {
		playtype = CallTypeHomo
	}

	// 获取在线用户数
	stat, _ := GetOnlineStat()
	resp.OnlineUserCount = stat.Online
	resp.MaleCount = stat.Male
	resp.FemaleCount = stat.Female

	// 判断用户是否被禁用或者是外挂
	var spam *meeting.SpamUserLevel
	err = thriftpool.Invoke(MeetingService, spamLevel, func(client interface{}) (err error) {
		c := client.(*meeting.MeetingServiceClient)
		spam, err = c.QueryUserSpamLevel(ctx, int64(uid))
		return
	})
	if err != nil {
		log.Error(prefix+"QueryUserSpamLevel", zap.Error(err), zap.Any("uid", uid))
	} else {
		log.Info(prefix+"QueryUserSpamLevel", zap.Any("uid", uid), zap.Any("spamUserLevel", spam))
	}
	if spam != nil && (spam.Level == 1 || spam.Level == 2 || spam.Level == 100 || spam.Level == 101 || spam.IsCheat) {
		//log.Warn(prefix+"user is disabled", zap.Any("uid", uid), zap.Any("spamLevel", spam.Level), zap.Any("spamCheat", spam.IsCheat))
		WriteSpamRecord(int64(uid), sex, playtype, spam.Level, spam.IsCheat)
		return
	}

	if val, _ := GetOnlineUser(uid); val != "" {
		log.Warn(prefix+"user is already online", zap.Any("uid", uid), zap.Any("user", val))
	}

	AddTalkingHeart(uid)

	isWhite, _ := IsWhite(uid, int(req.Sex))
	userValue, err := newUser(uid, int(req.MatchType), int(req.Sex), isWhite, req.Province)
	if err != nil {
		resp.Result = ResultMatchFailed
		resp.ErrorDesc = err.Error()

		log.Error(prefix+"newUser failed", zap.Uint32("uid", uid), zap.Error(err))
		return
	}

	// 添加到在线用户信息
	AddOnlineUser(uid, userValue)

	// 用户玩过的信息
	AddUserPlay(uid, userValue)

	// 用户行为跟踪
	PlayerBegin(int64(uid), sex, playtype)

	log.Info(prefix+"end", zap.Any("user", userValue))
	return
}

func (this *MatchServantObj) CancleMatchRandom(ctx context.Context, req *bilin.CancleMatchRandomReq) (resp *bilin.CancleMatchRandomResp, rpcErr error) {
	const prefix = "CancleMatchRandom "
	resp = &bilin.CancleMatchRandomResp{Result: ResultMatchSuccess, ErrorDesc: ""}

	uid, err := GetUid(ctx)
	if err != nil {
		resp.Result = ResultMatchFailed
		resp.ErrorDesc = err.Error()

		log.Error(prefix+"GetUid fail", zap.Error(err))
		return
	}

	log.Info(prefix+"begin", zap.Any("uid", uid), zap.Any("req", req.String()))

	val, _ := GetOnlineUser(uid)
	if val == "" {
		log.Warn(prefix+"user is not online", zap.Any("uid", uid))
		return
	}
	var user UserItem
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		log.Error("CancleMatchRandom json.Unmarshal", zap.Error(err), zap.Any("uid", uid))
	}

	// 从匹配队列里删除用户
	isWhite, _ := IsWhite(uid, int(req.Sex))
	_, err = delUser(uid, int(req.MatchType), int(req.Sex), isWhite, req.Province)

	// 删除用户和心跳
	DelTalkingHeart(uid)
	DelOnlineUser(uid)

	// 中途退出
	if req.Matchid != "0" {
		PlayerGiveup2(int64(uid))

		// 未被选中的单播
		match, err := GetMatchIdValue(req.Matchid)
		if err != nil {
			log.Error(prefix+"can not get match id", zap.Error(err), zap.Any("uid", uid), zap.Any("req", req.String()))
			return
		}

		matchResult := &bilin.OptionalMatchingResult{}
		err = json.Unmarshal([]byte(match), matchResult)
		if err != nil {
			log.Error(prefix+"do cancel broadcast json.Unmarshal fail",
				zap.Any("uid", uid),
				zap.Any("matchid", req.Matchid),
				zap.Any("match", match),
				zap.Error(err))
			return
		}

		var uids []int64
		for _, v := range matchResult.Attendees {
			if v.Uid == uid {
				v.Isonline = 0
				continue
			}
			if v.Isonline == 1 && !ExistUserItemRedis(v.Uid) {
				uids = append(uids, int64(v.Uid))
			}
		}
		if len(uids) == 0 {
			DelMatchIdValue(req.Matchid)
			return
		}

		value, err := json.Marshal(matchResult)
		if err != nil {
			log.Error(prefix+"do cancel broadcast json.Marshal fail",
				zap.Any("uid", uid),
				zap.Any("matchid", req.Matchid),
				zap.Any("match", match),
				zap.Any("updated match", string(value)))
			return
		}
		AddMatchIdValue(req.Matchid, string(value))

		log.Info(prefix+"do cancel broadcast",
			zap.Any("uid", uid),
			zap.Any("matchid", req.Matchid),
			zap.Any("match", match),
			zap.Any("updated match", string(value)),
			zap.Any("broadcast uids", uids))
		if user.Sex == Female {
			bc := bilin.MatchingResult{}
			bc.IsSelected = false
			unicast(uids, &bc, bilin.MaxType_MATCH_MSG, bilin.MinType_MATCH_MATCHINGRESULT_MINTYPE)
		} else {
			unicast(uids, matchResult, bilin.MaxType_MATCH_MSG, bilin.MinType_MATCH_OPTIONALMATCHINGRESULT_MINTYPE)
		}
	} else {
		PlayerGiveup1(int64(uid))
	}

	log.Info(prefix+"end", zap.Any("uid", uid))
	return
}

func (this *MatchServantObj) SelectMatchingResult(ctx context.Context, req *bilin.SelectMatchingResultReq) (resp *bilin.SelectMatchingResultResp, rpcErr error) {
	const prefix = "SelectMatchingResult "
	resp = &bilin.SelectMatchingResultResp{Result: ResultMatchSuccess, ErrorDesc: ""}

	uid, err := GetUid(ctx)
	if err != nil {
		resp.Result = ResultMatchFailed
		resp.ErrorDesc = err.Error()

		log.Error(prefix+"GetUid fail", zap.Error(err))
		return
	}

	log.Info(prefix+"begin", zap.Any("uid", uid), zap.Any("req", req.String()))

	// 被选中的单播
	if req.Uid != 0 {
		var selectUids []int64
		bc := bilin.MatchingResult{}
		bc.IsSelected = true
		selectUids = append(selectUids, int64(req.Uid))
		unicast(selectUids, &bc, bilin.MaxType_MATCH_MSG, bilin.MinType_MATCH_MATCHINGRESULT_MINTYPE)
		log.Info(prefix+"selected", zap.Any("uid", uid), zap.Any("selected uid", req.Uid))

		// 选中一次就把男性白名单移除
		DelMaleWhite(req.Uid)

		// 第一次玩随机
		AddFirstPlay(req.Uid)

		PlayerSuccess(int64(uid))
		PlayerSuccess(int64(req.Uid))
	}

	// 未被选中的单播
	match, err := GetMatchIdValue(req.Matchid)
	if err != nil {
		resp.Result = ResultMatchFailed
		resp.ErrorDesc = err.Error()

		log.Error(prefix+"can not get match id", zap.Error(err))
		return
	}

	matchResult := &bilin.OptionalMatchingResult{}
	err = json.Unmarshal([]byte(match), matchResult)
	if err != nil {
		log.Error(prefix+"do cancel broadcast json.Unmarshal fail",
			zap.Any("uid", uid),
			zap.Any("matchid", req.Matchid),
			zap.Any("match", match),
			zap.Error(err))
		return
	}

	var uids []int64
	bc := bilin.MatchingResult{}
	bc.IsSelected = false
	for _, v := range matchResult.Attendees {
		if v.Uid != req.Uid {
			if v.Uid != uid {
				uids = append(uids, int64(v.Uid))
				// 未被选择的，非机器人，加入男性白名单
				if v.Isonline > 0 {
					AddMaleWhite(v.Uid)
					PlayerGiveup2(int64(v.Uid))
				}
			}
		}
	}
	unicast(uids, &bc, bilin.MaxType_MATCH_MSG, bilin.MinType_MATCH_MATCHINGRESULT_MINTYPE)

	log.Info(prefix+"end", zap.Any("uid", uid), zap.Any("bc", bc.String()))
	return
}

func contains(s []int64, e int64) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

//申请通话  无脑转发
func (this *MatchServantObj) ApplyTalking(ctx context.Context, req *bilin.ApplyTalkingRequest) (*bilin.ApplyTalkingRespone, error) {
	resp := &bilin.ApplyTalkingRespone{Result: 0, ErrorDesc: ""}
	log.Info("ApplyTalking", zap.Any("Req", req))

	offline, err := unicast([]int64{int64(req.UnicastUid)},
		&bilin.ApplyTalkingNotify{req.RequestUid, req.Operation, req.Applyid},
		bilin.MaxType_MATCH_MSG, bilin.MinType_MATCH_APPLYTALKING_MINTYPE)
	if (err != nil || contains(offline, int64(req.UnicastUid))) && req.Operation == 0 {
		//resp.Result = 2
		// 未接来电
		var ComRet *common.ComRet
		err = thriftpool.Invoke(CallRecordService, callRecord, func(client interface{}) (err error) {
			c := client.(*callrecord.CallRecordServiceClient)
			ComRet, err = c.AddMissedCall(ctx, int64(req.RequestUid), int64(req.UnicastUid), 1, strconv.FormatUint(req.Applyid, 10))
			return
		})
		if err != nil || ComRet != nil && ComRet.Result_ != "success" {
			log.Error("ApplyTalking AddMissedCall", zap.Any("err", err), zap.Any("result", ComRet))
		}
	}

	return resp, nil
}

func (this *MatchServantObj) ReportTalking(ctx context.Context, req *bilin.ReportTalkingRequest) (*bilin.ReportTalkingResponse, error) {
	resp := &bilin.ReportTalkingResponse{Result: 0, ErrorDesc: ""}
	log.Info("ReportTalking", zap.Any("Req", req))

	offline, err := unicast([]int64{int64(req.UnicastUid)},
		&bilin.ReportTalkingNotify{req.RequestUid, req.Status, req.Reportid},
		bilin.MaxType_MATCH_MSG, bilin.MinType_MATCH_REPORTTALKING_MINTYPE)
	if err != nil || contains(offline, int64(req.UnicastUid)) {
		resp.Result = 2
		resp.ErrorDesc = "对方不在线"
	}

	return resp, nil
}

func (this *MatchServantObj) Talking(ctx context.Context, req *bilin.TalkingRequest) (*bilin.TalkingRespone, error) {
	uid, err := GetUid(ctx)
	if err != nil {
		log.Error("Talking GetUid fail", zap.Error(err))
		return &bilin.TalkingRespone{Result: 1, ErrorDesc: err.Error()}, nil
	}
	roomid, err := GetRoomid(ctx)
	if err != nil {
		log.Error("Talking GetRoomid fail", zap.Error(err))
		return &bilin.TalkingRespone{Result: 1, ErrorDesc: err.Error()}, nil
	}
	log.Info("Talking begin", zap.Any("uid", uid), zap.Any("roomid", roomid), zap.Any("Req", req.String()))

	bc := bilin.TalkingAction{}
	bc.Operation = req.Operation

	var (
		liveDisable bool
		liveId      int32
	)
	// uint32  operation = 1;  //操作请求，0：请求通话；1：取消通话；
	if req.Operation == 0 {
		// 获取房间id
		roomReq := &bilin.GenerateRoomReq{}
		comm := tars.NewCommunicator()
		objName := fmt.Sprintf("bilin.ccserver.CCServantObj")
		client := bilin.NewCCServantClient(objName, comm)
		resp, err := client.GenerateRoom(context.TODO(), roomReq)
		if err != nil {
			log.Error("Talking can not get room id", zap.Any("uid", uid), zap.Error(err))
			return &bilin.TalkingRespone{Result: 1, ErrorDesc: "can not get room id"}, nil
		}
		log.Info("Talking resp msg", zap.Any("uid", uid), zap.Any("resp", resp))

		bc.Cid = uint32(resp.RoomID)

		// 保存正在通话对象
		AddTalkingUser(req.RequestUid, req.UnicastUid, bc.Cid, req.Type)
		liveDisable = false
		liveId = int32(resp.RoomID)
	}

	if req.Operation == 1 {
		// 删除用户和心跳
		DelTalkingHeart(uid)
		DelOnlineUser(uid)
		// 删除正在通话对象
		peer, cid, talktype, begin, end, _ := DelTalkingUser(uid)
		if cid > 0 {
			var calltype int32
			switch talktype {
			case 1:
				calltype = 1
			default:
				calltype = 2
			}
			// 通话记录入库
			var ComRet *common.ComRet
			err = thriftpool.Invoke(CallRecordService, callRecord, func(client interface{}) (err error) {
				c := client.(*callrecord.CallRecordServiceClient)
				ComRet, err = c.AddCallRecordByCCServer(ctx,
					begin,
					strconv.FormatUint(uint64(cid), 10),
					end,
					"",
					int64(peer),
					calltype,
					int64(uid),
					"",
				)
				return
			})
			if err != nil || ComRet != nil && ComRet.Result_ != "success" {
				log.Error("Talking AddCallRecordByCCServer", zap.Any("err", err), zap.Any("result", ComRet),
					zap.Any("begin", begin),
					zap.Any("cid", cid),
					zap.Any("end", end),
					zap.Any("peer", peer),
					zap.Any("uid", uid),
					zap.Any("calltype", calltype))
			}
		}
		liveDisable = true
		if roomid == -1 {
			liveId = int32(cid)
		} else {
			liveId = int32(roomid)
		}
	}

	// 媒体鉴权
	var ComRet *common.ComRet
	err = thriftpool.Invoke(HotLineService, hotLine, func(client interface{}) (err error) {
		c := client.(*hotline.HotLineServiceClient)
		ComRet, err = c.UpdateLiveAuthority(ctx, liveId, []int64{int64(req.RequestUid)}, []hotline.TokenAuthType{hotline.TokenAuthType_tat_up_voice}, liveDisable)
		ComRet, err = c.UpdateLiveAuthority(ctx, liveId, []int64{int64(req.UnicastUid)}, []hotline.TokenAuthType{hotline.TokenAuthType_tat_up_voice}, liveDisable)
		ComRet, err = c.UpdateLiveAuthority(ctx, liveId, []int64{int64(uid)}, []hotline.TokenAuthType{hotline.TokenAuthType_tat_up_voice}, liveDisable)
		return
	})
	if err != nil || ComRet != nil && ComRet.Result_ != "success" {
		log.Error("Talking UpdateLiveAuthority", zap.Any("err", err), zap.Any("result", ComRet))
	}
	// 任务系统
	thriftpool.Invoke(HotLineDataService, hotLineData, func(client interface{}) (err error) {
		var ret int32
		c := client.(*hotline.DataServiceClient)
		tasks := []*hotline.TaskReq{
			{int64(req.RequestUid), "cashRetainTask", "2", int64(liveId)},
			{int64(req.UnicastUid), "cashRetainTask", "2", int64(liveId)},
		}
		if liveDisable {
			ret, err = c.Cancel(ctx, tasks)
		} else {
			ret, err = c.Start(ctx, tasks)
		}
		if err != nil {
			log.Error("Talking Start/Cancel Task", zap.Any("err", err), zap.Any("ret", ret))
		}
		return
	})
	// 生成话单
	var callOp, callType int
	if liveDisable {
		callOp = CallOpEnd
	} else {
		callOp = CallOpStart
	}
	switch req.Type {
	case 1:
		callType = CallTypeDirect
	case 2:
		callType = CallTypeHetero
	case 3:
		callType = CallTypeHomo
	default:
		callType = CallTypeUnknown
	}
	if liveId > 0 {
		WriteCallRecord(callOp, int64(req.RequestUid), int64(req.UnicastUid), int64(liveId), callType)
	}

	var uids []int64

	uids = append(uids, int64(req.UnicastUid))
	uids = append(uids, int64(req.RequestUid))
	unicast(uids, &bc, bilin.MaxType_MATCH_MSG, bilin.MinType_MATCH_TALKACTION_MINTYPE)

	log.Info("Talking end", zap.Any("uid", uid), zap.Any("bc", bc.String()))

	return &bilin.TalkingRespone{Result: 0, ErrorDesc: ""}, nil
}

func (this *MatchServantObj) TalkingHeartbeat(ctx context.Context, req *bilin.TalkingHeartbeatRequest) (*bilin.TalkingHeartbeatRespone, error) {
	uid, err := GetUid(ctx)
	if err != nil {
		log.Error("TalkingHeartbeat GetUid fail", zap.Error(err))
		return &bilin.TalkingHeartbeatRespone{Result: 1, ErrorDesc: err.Error()}, nil
	}
	roomid, err := GetRoomid(ctx)
	if err != nil {
		log.Error("TalkingHeartbeat GetRoomid fail", zap.Error(err))
		return &bilin.TalkingHeartbeatRespone{Result: 1, ErrorDesc: err.Error()}, nil
	}
	if roomid == -1 {
		//log.Warn("TalkingHeartbeat does not subscribe room push", zap.Any("uid", uid))
		// 这种情况对应于用户正在玩随机匹配，但没有在通话
	} else {
		users, _ := onlinequery.GetRoomUser(roomid)
		if len(users) < 2 {
			//bc := bilin.TalkingAction{}
			//bc.Operation = 1
			//bc.CancelReason = 1
			//unicast([]int64{int64(uid)}, &bc, bilin.MaxType_MATCH_MSG, bilin.MinType_MATCH_TALKACTION_MINTYPE)
			log.Info("TalkingHeartbeat cancel talking because of peer left", zap.Any("uid", uid), zap.Any("roomid", roomid))
		}
	}

	AddTalkingHeart(uid)
	return &bilin.TalkingHeartbeatRespone{Result: 0, ErrorDesc: ""}, nil
}

func (this *MatchServantObj) GetComfortWord(ctx context.Context, req *bilin.GetComfortWordRequest) (*bilin.GetComfortWordRespone, error) {
	uid, err := GetUid(ctx)
	if err != nil {
		log.Error("GetComfortWord GetUid fail", zap.Error(err))
		return &bilin.GetComfortWordRespone{Result: 1, ErrorDesc: err.Error(), ComforWord: ""}, nil
	}

	// 是否是第一次玩的时间戳
	timeValue, _ := GetFirstPlay(uid)

	// 失败的次数
	MaleFailvalue, _ := GetMaleWhite(uid)

	log.Info("GetComfortWord begin", zap.Any("uid", uid), zap.Any("timeValue", timeValue), zap.Any("MaleFailvalue", MaleFailvalue))

	// 响应
	resp := &bilin.GetComfortWordRespone{}

	times, _ := strconv.ParseUint(timeValue, 10, 0)
	fails, _ := strconv.ParseUint(MaleFailvalue, 10, 0)

	if times == 0 && fails == 2 {
		resp.ComforWord = "取个好听的名字，换张帅帅的照片可以提高匹配率哦！"
	} else {
		no, _ := GetUserComfortWordNo(uid)
		world, _ := GetComfortWord(no)
		resp.ComforWord = world
	}

	log.Info("GetComfortWord end", zap.Any("resp", resp.String()))
	return resp, nil
}

func (this *MatchServantObj) GetRandomAvatar(ctx context.Context, req *bilin.GetRandomAvatarReq) (*bilin.GetRandomAvatarResp, error) {
	uid, err := GetUid(ctx)
	if err != nil {
		log.Error("GetRandomAvatar GetUid fail", zap.Error(err))
		return &bilin.GetRandomAvatarResp{Result: 1, ErrorDesc: err.Error()}, nil
	}

	log.Info("GetRandomAvatar begin", zap.Any("uid", uid), zap.Any("Req", req.String()))

	playValue, _ := GetUserPlayAll(100)

	var uids []uint32
	count := 0
	for _, v := range playValue {
		var item UserItem
		json.Unmarshal([]byte(v), &item)
		if int32(item.Sex) != req.Sex {
			uids = append(uids, item.Uid)
			count++
		}
	}
	if count == 0 {
		return &bilin.GetRandomAvatarResp{Result: 0, ErrorDesc: "no content"}, nil
	}

	// 发请求查询头像
	userInfoReq := &u.GetUserInfoReq{}

	//利用当前时间的UNIX时间戳初始化rand包
	for i := 0; i < 8; i++ {
		x := rand.Intn(count)
		userInfoReq.Uids = append(userInfoReq.Uids, uint64(uids[x]))
	}

	comm := tars.NewCommunicator()
	objName := fmt.Sprintf("bilin.userinfocenter.UserInfoCenterObj")
	client := u.NewUserInfoCenterObjClient(objName, comm)
	userResp, err := client.GetUserInfo(context.TODO(), userInfoReq)

	resp := &bilin.GetRandomAvatarResp{}
	if err != nil {
		log.Error("GetRandomAvatar GetUserInfo err", zap.Error(err))
		resp.Result = 1
		resp.ErrorDesc = "GetUserInfo failed"
		return resp, nil
	}

	// 返回头像
	for _, v := range userInfoReq.Uids {
		info := getUserInfo(userResp.Users[uint64(v)])
		resp.Avatars = append(resp.Avatars, info.Avatar)
	}

	return resp, nil
}
