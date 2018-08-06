package handler

import (
	"bilin/bcserver/domain/entity"
	"bilin/common/thriftpool"
	"bilin/protocol"
	"bilin/thrift/gen-go/hotline"
	"context"
	"fmt"
	"strings"
	"time"

	"bilin/bcserver/bccommon"
	"bilin/bcserver/config"
	"bilin/bcserver/domain/adapter"
	"bilin/bcserver/domain/collector"
	"bilin/bcserver/domain/service"
	"bilin/thrift/gen-go/bilin_msg_filter"
	"bilin/thrift/gen-go/common"
	"bilin/thrift/gen-go/findfriendsbroadcast"
	"bilin/thrift/gen-go/officialhotline"
	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/kissgo/httpmetrics"
	"code.yy.com/yytars/goframework/tars"
	"sort"
)

const (
	GBaseRoomInfo_AUTOLINK = 2 // 2018-05-11 comment Discarded in the future
	PINGROOM               = 0
	ENTERROOM              = 1
)

var _ bilin.BCServantServer = &BCServantObj{}

// BCServantObj 包含所有直播间相关的操作
type BCServantObj struct {
	hotLine             thriftpool.Pool
	hotLineDataService  thriftpool.Pool
	msgFilter           thriftpool.Pool
	officailHotLine     thriftpool.Pool
	findFriendBroadcast thriftpool.Pool
	rooms               map[uint64]*entity.Room
	invisibleUids       []uint64
	relationlistclient  bilin.RelationListServantClient
}

// NewBCServantObj 被main调用，初始化
func NewBCServantObj() *BCServantObj {
	service.RedisInit()
	service.MysqlInit()
	service.KafkaProducerInit()
	collector.RoomStatisticsInit("bcserver")
	log.Info("NewBCServantObj begin", zap.Any("appconfig", config.GetAppConfig()))

	hotLine, err := thriftpool.NewChannelPool(0, 1000, service.CreateHotLineServiceConn)
	if err != nil {
		log.Panic("can not create thrift connection pool hotLine", zap.Any("err", err))
	}
	hotLineDataService, err := thriftpool.NewChannelPool(0, 1000, service.CreateHotLineDataServiceConn)
	if err != nil {
		log.Panic("can not create thrift connection pool hotLineDataService", zap.Any("err", err))
	}
	msgFilter, err := thriftpool.NewChannelPool(0, 1000, service.CreateMsgFilterServiceConn)
	if err != nil {
		log.Panic("can not create thrift connection pool hotLine", zap.Any("err", err))
	}
	officailHotLine, err := thriftpool.NewChannelPool(0, 1000, service.CreateOfficailHotLineServiceConn)
	if err != nil {
		log.Panic("can not create thrift connection pool officialHotline", zap.Any("err", err))
	}

	findFriendBroadcast, err := thriftpool.NewChannelPool(0, 1000, service.CreateFindFriendsBroadcastServiceConn)
	if err != nil {
		log.Panic("can not create thrift connection pool FindFriendBroadcast", zap.Any("err", err))
	}

	comm := tars.NewCommunicator()
	s := &BCServantObj{
		hotLine:             hotLine,
		hotLineDataService:  hotLineDataService,
		msgFilter:           msgFilter,
		officailHotLine:     officailHotLine,
		findFriendBroadcast: findFriendBroadcast,
		invisibleUids:       config.GetAppConfig().InvisibleUids,
		relationlistclient:  bilin.NewRelationListServantClient("bilin.relationlist.RelationListPbObj", comm),
	}
	go thriftpool.Ping(service.HotLineService, s.hotLine, func(client interface{}) (err error) {
		c := client.(*hotline.HotLineServiceClient)
		_, err = c.Ping(context.TODO())
		return
	}, 10*time.Minute)
	go thriftpool.Ping(service.OfficailService, s.officailHotLine, func(client interface{}) (err error) {
		c := client.(*officialhotline.OfficialHotlineServiceClient)
		_, err = c.Ping(context.TODO())
		return
	}, 10*time.Minute)
	return s
}

func (this *BCServantObj) CommonCheckAuth(roomid uint64, userid uint64) (room *entity.Room, user *entity.User, err error) {
	const prefix = "CommonCheckAuth "

	defer func(now time.Time) {
		httpmetrics.DefReport("CommonCheckAuth", 0, now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	if room = collector.GetRoomInfoByRoomId(roomid); room == nil {
		return nil, nil, fmt.Errorf("room not exist, roomid:%d", roomid)
	}

	if user, _ = service.RedisGetUser(roomid, userid); user == nil {
		return nil, nil, fmt.Errorf("user :%d not in room:%d", userid, roomid)
	}

	return room, user, nil
}

func transThriftErrToCommRet(status int32) (commRet *bilin.CommonRetInfo) {
	switch status {
	case -1:
		return bccommon.UserDefinedFailed(bilin.CommonRetInfo_ENTER_USER_NO_RIGHT, fmt.Sprintf("房主拒绝你进来哦~"))
	case -2:
		return bccommon.UserDefinedFailed(bilin.CommonRetInfo_ENTER_ROOM_NOT_START, fmt.Sprintf("直播异常，主持人未开播"))
	case -3:
		return bccommon.UserDefinedFailed(bilin.CommonRetInfo_ENTER_ROOM_CLOSED, fmt.Sprintf("直播间已被关闭"))
	case -4:
		return bccommon.UserDefinedFailed(bilin.CommonRetInfo_ENTER_ROOM_FORBIDDEN, fmt.Sprintf("直播间涉嫌违规"))
	default:
		return bccommon.UserDefinedFailed(bilin.CommonRetInfo_ENTER_BAD_NETWORK, fmt.Sprintf("服务器开小差了，再试试呗~"))
	}
}

func (this *BCServantObj) privateEnterBroRoom(ctx context.Context, RoomId uint64, UserId uint64, Mapextend map[uint32]string, EnterType int32, RoomPwd string) (room *entity.Room, user *entity.User, commRet *bilin.CommonRetInfo) {
	const prefix = "privateEnterBroRoom "
	log.Info(prefix+"begin", zap.Any("RoomId", RoomId), zap.Any("UserId", UserId))

	//查询观众信息
	var joinHotLineRet *hotline.JoinHotLineRet
	thriftErr := thriftpool.Invoke(service.HotLineService, this.hotLine, func(client interface{}) (err error) {
		c := client.(*hotline.HotLineServiceClient)
		joinHotLineRet, err = c.JoinHotLine(ctx, int32(RoomId), int64(UserId))
		return
	})
	if thriftErr != nil {
		log.Warn(prefix+"JoinHotLine", zap.Any("err", thriftErr), zap.Any("joinHotLineRet", joinHotLineRet))
		return nil, nil, bccommon.UserDefinedFailed(bilin.CommonRetInfo_ENTER_BAD_NETWORK, fmt.Sprintf("服务器开小差了，再试试呗~"))
	}
	if joinHotLineRet.Result_ != "success" {
		log.Warn(prefix+"JoinHotLine", zap.Any("err", thriftErr), zap.Any("joinHotLineRet", joinHotLineRet))
		return nil, nil, bccommon.UserDefinedFailed(bilin.CommonRetInfo_ENTER_ROOM_FAILED, fmt.Sprintf("服务器开小差了，再试试呗~"))
	}
	if joinHotLineRet.Status < 0 {
		log.Warn(prefix+"JoinHotLine", zap.Any("thriftErr", thriftErr), zap.Any("joinHotLineRet", joinHotLineRet))
		return nil, nil, transThriftErrToCommRet(joinHotLineRet.Status)
	}
	log.Info(prefix+"JoinHotLine", zap.Any("RoomId", RoomId), zap.Any("UserId", UserId), zap.Any("ret", *joinHotLineRet))

	room = collector.InitRoomByJavaResult(RoomId, UserId, joinHotLineRet)

	//if room.Status == bilin.BaseRoomInfo_FORBIDDEN {
	//	log.Error(prefix+"BaseRoomInfo_FORBIDDEN ", zap.Any("err", thriftErr), zap.Any("room", room))
	//	return nil, nil, transThriftErrToCommRet(-3)
	//}

	if joinHotLineRet.Status == entity.ROLE_HOST {
		collector.HostEnterRoom(room, UserId, joinHotLineRet)
	} else {
		if ENTERROOM == EnterType { //观众进房间，需要鉴权
			if pass, authErr := collector.BizAuth(room, UserId, RoomPwd); authErr != nil || !pass {
				log.Error(prefix+"Incorrect password", zap.Any("password", RoomPwd))
				return nil, nil, bccommon.UserDefinedFailed(bilin.CommonRetInfo_ENTER_ROOM_PWDERR, fmt.Sprintf("房间已上锁，请输入正确密码"))
			}
		}

	}

	user, commRet = collector.InitUserByJavaResult(RoomId, UserId, joinHotLineRet, Mapextend)
	return room, user, commRet
}

func (this *BCServantObj) notifyUserEnter(ctx context.Context, user *entity.User, room *entity.Room) (err error) {
	const prefix = "notifyUserEnter "

	//直播间活动任务相关操作
	go thriftpool.Invoke(service.HotLineDataService, this.hotLineDataService, func(client interface{}) (err error) {
		c := client.(*hotline.DataServiceClient)
		intRes, err := c.Start(ctx, []*hotline.TaskReq{service.NewNotifyTask(user.UserID, room.Roomid, service.UserEnterRoomTask)})
		if err != nil {
			log.Warn(prefix+"UserEnterRoomTask", zap.Any("RoomId", room.Roomid), zap.Any("UserId", user.UserID), zap.Any("err", err), zap.Any("intRes", intRes))
			return
		}

		log.Warn(prefix+"UserEnterRoomTask begin", zap.Any("intRes", intRes))

		if user.Role == entity.ROLE_HOST {
			intHostRes, err2 := c.Start(ctx, []*hotline.TaskReq{service.NewNotifyTask(user.UserID, room.Roomid, service.HostStartLivingTask)})
			if err2 != nil {
				log.Warn(prefix+"HostStartLivingTask", zap.Any("RoomId", room.Roomid), zap.Any("UserId", user.UserID), zap.Any("err", err2), zap.Any("intHostRes", intHostRes))
			}

			log.Info(prefix+"HostStartLivingTask begin", zap.Any("intHostRes", intHostRes))
		}

		return
	})

	if room.RoomType2 == service.OFFICAIL_ROOM && user.Role == entity.ROLE_HOST {
		//查看当前主持人麦位是否有人
		hasUser := false
		userlist, _ := service.RedisGetOnMikeUserList(room.Roomid)
		for _, mikeUser := range userlist {
			if mikeUser.MikeIndex == 0 && mikeUser.UserID != 0 {
				hasUser = true
				log.Info(prefix, zap.Any("host mike hasUser", mikeUser.UserID))
				break
			}
		}

		if hasUser { //当前麦位已经有人，等待java通知上麦
			return
		} else {
			//官频没人，直接上麦，并通知java
			var ComRet *officialhotline.OfficialHotlineRet
			err = thriftpool.Invoke(service.HotLineService, this.officailHotLine, func(client interface{}) (err error) {
				c := client.(*officialhotline.OfficialHotlineServiceClient)
				ComRet, err = c.OnOfficialMike(ctx, int32(room.Roomid), int64(user.UserID))
				return
			})
			if err != nil || ComRet.Result_ != "success" {
				log.Warn(prefix+"OnOfficialMike", zap.Any("err", err), zap.Any("ComRet", ComRet))
			} else {
				log.Info(prefix+"OnOfficialMike", zap.Any("err", err), zap.Any("ComRet", ComRet))
			}
		}

	}

	//主持人上0号麦
	if room.Owner == user.UserID && user.Role == entity.ROLE_HOST {
		if room.RoomType2 != service.OFFICAIL_ROOM {
			//对房间做一些初始化工作
			service.RedisClearMikeWheat(room.Roomid)
			collector.InitMikeWheatInfo(room)
			if mikelist, _ := service.RedisGetOnMikeUserList(room.Roomid); mikelist != nil {
				for _, mikeuser := range mikelist {
					this.RemoveUserFromMike(ctx, room, mikeuser)
				}

				//通知mikelist变化
				adapter.PushMikeListInfoToRoom(room)
			}
		}

		user.MikeIndex = 0
		this.AddUserToMike(ctx, room, user)

		//通知mikelist变化
		adapter.PushMikeListInfoToRoom(room)

		//经客户端要求，需要发送一个baseRoomInfo下去
		adapter.PushBaseRoomInfoToRoom(room)
	}

	//通知直播间内用户列表变化增量，为了兼容坑爹的老版本
	adapter.PushUserListChangeToRoom(room, []*entity.User{user}, nil)

	//通知直播间其他人，列表变化
	adapter.PushUserListInfoToRoom(room)

	return
}

func (this *BCServantObj) getUserPrivilegeInfo(ctx context.Context, room *entity.Room, owner uint64, guest_uid uint64) (result *bilin.UserPrivilegeInfoInRoom) {
	const prefix = "getUserPrivilegeInfo "
	log.Info(prefix+"begin", zap.Any("roomid", room.Roomid), zap.Any("owner", owner), zap.Any("guest_uid", guest_uid))

	result = &bilin.UserPrivilegeInfoInRoom{}
	//用户特权信息
	result.Headgear, _ = service.RedisGetUserHeadgear(guest_uid)

	if room.Relationlistswitch == bilin.BaseRoomInfo_CLOSERELATIONLIST || owner == guest_uid {
		return
	}

	//查询用户与主播的亲密度，勋章等信息
	resp, rpcErr := this.relationlistclient.GetUserRelationMedal(ctx, &bilin.GetUserRelationMedalReq{
		Header: &bilin.Header{Userid: guest_uid, Roomid: room.Roomid},
		Owner:  owner,
	})
	if rpcErr != nil {
		log.Error(prefix+"GetUserRelationMedal ", zap.Any("rpcErr", rpcErr), zap.Any("resp", resp))
		return
	}

	result.Medaltext = resp.Medalname
	result.Medalurl = resp.MedalUrl
	log.Info(prefix+"end", zap.Any("roomid", room.Roomid), zap.Any("owner", owner), zap.Any("guest_uid", guest_uid), zap.Any("PrivilegeInfo", result))
	return
}

// EnterBroRoom 是用户进入直播间的操作
func (this *BCServantObj) EnterBroRoom(ctx context.Context, req *bilin.EnterBroRoomReq) (resp *bilin.EnterBroRoomResp, err error) {
	const prefix = "EnterBroRoom "
	resp = &bilin.EnterBroRoomResp{Commonret: bccommon.SUCCESSMESSAGE}
	roomid := req.Header.Roomid
	userid := req.Header.Userid

	var user *entity.User
	var room *entity.Room
	var commRet *bilin.CommonRetInfo

	if bccommon.Contains(this.invisibleUids, userid) {
		if room = collector.GetRoomInfoByRoomId(roomid); room != nil && room.Owner != userid { //如果是主持人,直接开播
			goto RETURN
		}

	}

	defer func(now time.Time) {
		httpmetrics.DefReport("EnterBroRoom", int64(resp.Commonret.Ret), now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	if userid == 0 || roomid == 0 {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_ILLEGAL_MESSAGE, bccommon.COMMERRORILLEGALDESC)
		log.Warn("[-]"+prefix+"failed", zap.Any("resp", resp), zap.Any("err", err))
		return resp, nil
	}

	log.Info(prefix+"begin", zap.Any("req", req))

	if room, user, commRet = this.privateEnterBroRoom(ctx, roomid, userid, req.Header.Mapextend, ENTERROOM, req.Roompwd); commRet != bccommon.SUCCESSMESSAGE {
		if commRet.Ret == bilin.CommonRetInfo_ENTER_ROOM_ALREADY_IN_ROOM {
			goto RETURN
		}
		resp.Commonret = commRet
		log.Warn("[-]EnterBroRoom failed", zap.Any("resp", resp))
		return
	}

	this.notifyUserEnter(ctx, user, room)

RETURN:
	//返回当前直播间信息
	resp.Allroominfo = collector.AllRoomInfo(room)

	//用户特权信息  跟老版本(5.1)兼容
	headgear, _ := service.RedisGetUserHeadgear(req.Header.Userid)
	resp.Allroominfo.PrivilegeInfo = &bilin.UserPrivilegeInfoInRoom{Headgear: headgear}

	//新版本（>=5.2）用户特权信息都放在resp下，也不知道当初为什么会放roominfo里的，估计脑子进水了。。。。
	resp.Privilegeinfo = this.getUserPrivilegeInfo(ctx, room, room.Owner, userid)

	log.Info("[+]EnterBroRoom success", zap.Any("req", req), zap.Any("resp", resp))

	return
}

// 客户端连接媒体结果通知
func (this *BCServantObj) ConnMediaResult(context.Context, *bilin.ConnMediaResultReq) (*bilin.ConnMediaResultResp, error) {
	return nil, nil
}

//
// 直播间PING请求
func (this *BCServantObj) PingBroRoom(ctx context.Context, req *bilin.PingBroRoomReq) (resp *bilin.PingBroRoomResp, err error) {
	const prefix = "PingBroRoom "
	resp = &bilin.PingBroRoomResp{Commonret: bccommon.SUCCESSMESSAGE}
	roomid := req.Header.Roomid
	userid := req.Header.Userid

	if bccommon.Contains(this.invisibleUids, userid) {
		return
	}

	defer func(now time.Time) {
		httpmetrics.DefReport("PingBroRoom", int64(resp.Commonret.Ret), now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	if userid == 0 || roomid == 0 {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_ILLEGAL_MESSAGE, bccommon.COMMERRORILLEGALDESC)
		log.Error("[-]"+prefix+"failed", zap.Any("resp", resp), zap.Any("userid", userid), zap.Any("roomid", roomid))
		return resp, nil
	}

	if exist, e := service.RedisSetPingTime(roomid, userid); e != nil || exist {
		return
	}

	// key not exist  重进频道
	var user *entity.User
	var room *entity.Room
	var commRet *bilin.CommonRetInfo
	if room, user, commRet = this.privateEnterBroRoom(ctx, roomid, userid, req.Header.Mapextend, PINGROOM, ""); commRet != bccommon.SUCCESSMESSAGE {
		resp.Commonret = commRet
		log.Error(prefix+"[-]privateEnterBroRoom failed", zap.Any("resp", resp))
		return
	}

	this.notifyUserEnter(ctx, user, room)

	//返回当前直播间信息
	adapter.PushAllRoomInfoToUser(room, []int64{int64(userid)})
	return
}

func (this *BCServantObj) RemoveUserFromMike(ctx context.Context, room *entity.Room, user *entity.User) (err error) {
	const prefix = "RemoveUserFromMike "
	if err = service.RedisRemoveUserFromMike(room.Roomid, user); err != nil {
		log.Warn("RemoveUserFromMike ", zap.Any("roomid", room.Roomid), zap.Any("user", user))
	}

	//notify user on unmike
	adapter.PushNotifyToUser(room.Roomid, []int64{int64(user.UserID)}, &bilin.MikeNotify{Optuserid: room.Owner, Opt: uint32(bilin.MikeOperationReq_UNMIKE)}, bilin.MinType_BC_NotifyRoomMickOpt)

	var ComRet *common.ComRet
	err = thriftpool.Invoke(service.HotLineService, this.hotLine, func(client interface{}) (err error) {
		c := client.(*hotline.HotLineServiceClient)
		ComRet, err = c.UpdateLiveAuthority(ctx, int32(room.Roomid), []int64{int64(user.UserID)}, []hotline.TokenAuthType{hotline.TokenAuthType_tat_up_voice}, true)
		return
	})
	if err != nil || ComRet.Result_ != "success" {
		log.Error(prefix+"UpdateLiveAuthority", zap.Any("err", err), zap.Any("result", ComRet.Result_), zap.Any("errorMsg", ComRet.ErrorMsg))
	}

	// bug fix ,主播下播需要通知亲密榜统计停止
	resp, rpcErr := this.relationlistclient.RSUserMikeOption(ctx, &bilin.RSUserMikeOptionReq{
		Header: &bilin.Header{Userid: user.UserID, Roomid: room.Roomid},
		Owner:  room.Owner,
		Opt:    bilin.RSUserMikeOptionReq_UNMIKE,
	})
	if rpcErr != nil {
		log.Error("RSUserMikeOption ", zap.Any("rpcErr", rpcErr), zap.Any("resp", resp))
	}

	log.Info("RemoveUserFromMike ", zap.Any("roomid", room.Roomid), zap.Any("userid", user.UserID), zap.Any("ComRet", ComRet), zap.Any("err", err))
	return nil
}

func (this *BCServantObj) AddUserToMike(ctx context.Context, room *entity.Room, user *entity.User) (err error) {
	const prefix = "AddUserToMike "

	user.OnMikeTime = uint64(time.Now().Unix())
	if err = service.RedisAddUserToMike(room.Roomid, user); err != nil {
		log.Warn("AddUserToMike ", zap.Any("roomid", room.Roomid), zap.Any("user", user))
	}

	//notify user on mike
	adapter.PushNotifyToUser(room.Roomid, []int64{int64(user.UserID)}, &bilin.MikeNotify{Optuserid: room.Owner, Opt: uint32(bilin.MikeOperationReq_ONMIKE)}, bilin.MinType_BC_NotifyRoomMickOpt)

	var ComRet *common.ComRet
	err = thriftpool.Invoke(service.HotLineService, this.hotLine, func(client interface{}) (err error) {
		c := client.(*hotline.HotLineServiceClient)
		ComRet, err = c.UpdateLiveAuthority(ctx, int32(room.Roomid), []int64{int64(user.UserID)}, []hotline.TokenAuthType{hotline.TokenAuthType_tat_up_voice}, false)
		return
	})
	if err != nil || ComRet.Result_ != "success" {
		log.Error(prefix+"UpdateLiveAuthority", zap.Any("err", err), zap.Any("result", ComRet.Result_), zap.Any("errorMsg", ComRet.ErrorMsg))
	}

	//亲密榜更新 嘉宾上麦时才需要更新
	resp, rpcErr := this.relationlistclient.RSUserMikeOption(ctx, &bilin.RSUserMikeOptionReq{
		Header: &bilin.Header{Userid: user.UserID, Roomid: room.Roomid},
		Owner:  room.Owner,
		Opt:    bilin.RSUserMikeOptionReq_ONMIKE,
	})
	if rpcErr != nil {
		log.Error("RSUserMikeOption ", zap.Any("rpcErr", rpcErr), zap.Any("resp", resp))
	}

	log.Info("AddUserToMike ", zap.Any("roomid", room.Roomid), zap.Any("userid", user.UserID), zap.Any("ComRet", ComRet), zap.Any("err", err))
	return err
}

//用户自动上麦，补前面那个下麦的人的坑位
func (this *BCServantObj) supplementUserOnMike(ctx context.Context, room *entity.Room, mikeIdx uint32) (success bool) {
	if room.RoomType == bilin.BaseRoomInfo_ROOMTYPE_RADIO {
		return false
	}

	if room.GetLinkStatus() == bilin.BaseRoomInfo_CLOSELINK || room.GetAutoLink() == bilin.BaseRoomInfo_CLOSEAUTOTOMIKE {
		return false
	}

	var user *entity.User
	if user, _ = service.RedisGetOneApplyMikeUser(room.Roomid); user == nil { //这里已经从applylist中pop出来了
		return false
	}

	user.MikeIndex = mikeIdx
	this.AddUserToMike(ctx, room, user)

	return true
}

//自动连线状态下，选一个位置让用户上麦
func (this *BCServantObj) autoOnMikeOperation(ctx context.Context, room *entity.Room, user *entity.User) (ret bool) {
	//检查是否有空的麦位
	mikeMap, _ := service.RedisGetAllMikeWheatStatus(room.Roomid)
	var identifier []int
	for idx, value := range mikeMap {
		if value == int(bilin.MikeInfo_EMPTY) && idx != 0 {
			identifier = append(identifier, idx)
		}
	}

	if len(identifier) > 0 {
		sort.Ints(identifier)
		user.MikeIndex = uint32(identifier[0])
		this.AddUserToMike(ctx, room, user)
		//从排麦列表中删除等待的用户
		service.RedisRemoveUserFromApplyMikeList(room.Roomid, user.UserID)
		log.Info("autoOnMikeOperation ", zap.Any("roomid", room.Roomid), zap.Any("userid", user.UserID), zap.Any("mikeIndex", user.MikeIndex))
		return true
	}
	return false
}

func (this *BCServantObj) fillAllEmptyMikeWheat(ctx context.Context, room *entity.Room) (result bool) {
	//检查排麦列表，填充麦上用户
	result = false
	var userlist []*entity.User
	userlist, _ = service.RedisGetApplyMikeUserList(room.Roomid)
	for _, item := range userlist {
		//用户上麦
		if ret := this.autoOnMikeOperation(ctx, room, item); !ret {
			break
		}
		result = true //有人上麦成功就设为true
	}

	log.Info("fillAllEmptyMikeWheat ", zap.Any("roomid", room.Roomid), zap.Any("result", result))
	return
}

// AudienceLinkOperation 观众请求麦位、取消麦位
func (this *BCServantObj) AudienceLinkOperation(ctx context.Context, req *bilin.AudienceLinkOperationReq) (resp *bilin.AudienceLinkOperationResp, err error) {
	const prefix = "AudienceLinkOperation "
	resp = &bilin.AudienceLinkOperationResp{Commonret: bccommon.SUCCESSMESSAGE}
	roomId := req.Header.Roomid
	userId := req.Header.Userid
	log.Debug(prefix+"begin", zap.Any("req", req))

	defer func(now time.Time) {
		httpmetrics.DefReport("AudienceLinkOperation", int64(resp.Commonret.Ret), now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	var user *entity.User
	var room *entity.Room
	if room, user, err = this.CommonCheckAuth(roomId, userId); err != nil {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_ILLEGAL_MESSAGE, bccommon.COMMERRORILLEGALDESC)
		log.Error("[-]"+prefix+"failed", zap.Any("resp", resp), zap.Any("err", err))
		return resp, nil
	}

	if req.Linkop == bilin.AudienceLinkOperationReq_LINK {
		switch room.GetLinkStatus() {
		case bilin.BaseRoomInfo_CLOSELINK:
			resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_AUDIENCE_LINK_OPERATION_FAILED, fmt.Sprintf("房间已关闭连线"))
			log.Error("[-]AudienceLinkOperation failed", zap.Any("resp", resp), zap.Any("Linkstatus", room.GetLinkStatus()))
			return resp, nil
		case GBaseRoomInfo_AUTOLINK:
			fallthrough
		case bilin.BaseRoomInfo_OPENLINK:
			//如果麦序开关已打开，1+6模板的直播间内，老版本（5.0以前）用户无法申请上麦，但是可以下麦
			if room.Maixuswitch == bilin.BaseRoomInfo_OPENMAIXU && room.RoomType == bilin.BaseRoomInfo_ROOMTYPE_SIX && len(user.Version) == 0 {
				resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_AUDIENCE_LINK_OPERATION_FAILED, fmt.Sprintf("请先升级到最新版本哦~"))
				log.Error("[-]AudienceLinkOperation failed", zap.Any("resp", resp), zap.Any("version not support", user.Version))
				return resp, nil
			}

			//电台模板不支持自动上麦
			if room.RoomType == bilin.BaseRoomInfo_ROOMTYPE_RADIO {
				break
			}

			//如果用户已经在麦上，直接返回false
			if retUser, _ := service.RedisGetUserOnMike(room.Roomid, user.UserID); retUser != nil {
				resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_AUDIENCE_LINK_OPERATION_FAILED, fmt.Sprintf("您已经在麦上了哦~"))
				log.Warn("autoOnMikeOperation user already on mike.", zap.Any("roomid", room.Roomid), zap.Any("userid", user.UserID), zap.Any("mikeIndex", user.MikeIndex))
				return resp, nil
			}

			if room.GetAutoLink() == bilin.BaseRoomInfo_OPENAUTOTOMIKE { //直接上麦
				if ret := this.autoOnMikeOperation(ctx, room, user); ret {
					adapter.PushMikeListInfoToRoom(room)
					return resp, nil
				}
			}
		default:
			resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_AUDIENCE_LINK_OPERATION_STATUS_ERR, fmt.Sprintf("连线状态错误"))
			log.Error("[-]AudienceLinkOperation failed, unknown linkstatus", zap.Any("roomId", roomId), zap.Any("Linkstatus", room.GetLinkStatus()))
			return resp, nil
		}

		// 如果已经在排麦列表里， 直接返回错误
		if exist := service.RedisIfUserOnApplyMikeList(roomId, userId); exist {
			resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_AUDIENCE_LINK_OPERATION_FAILED, fmt.Sprintf("用户已经在申请连线列表中"))
			log.Error("[-]AudienceLinkOperation failed", zap.Any("resp", resp))
			return resp, nil
		}

		if cout, _ := service.RedisGetApplyMikeUserCount(roomId); cout >= collector.ApplyMikeLimits {
			resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_AUDIENCE_LINK_OPERATION_FULL_MEMBER, fmt.Sprintf("当前排队人数较多，先等等呗~"))
			log.Error("[-]AudienceLinkOperation failed", zap.Any("resp", resp))
			return resp, nil
		}

		service.RedisAddUserToApplyMikeList(roomId, user)
	} else {
		//取消申请，或者下麦
		service.RedisRemoveUserFromApplyMikeList(roomId, user.UserID)

		//如果在麦序上， 需要下麦
		if mikeUser, _ := service.RedisGetUserOnMike(roomId, user.UserID); mikeUser != nil {
			this.RemoveUserFromMike(ctx, room, mikeUser)
			this.supplementUserOnMike(ctx, room, mikeUser.MikeIndex)

			adapter.PushMikeListInfoToRoom(room)
		}

	}

	//通知直播间其他人，申请排麦人数变化
	adapter.PushBaseRoomInfoToRoom(room)

	log.Debug(prefix+"end", zap.Any("resp", resp))
	return
}

//
// 用户离开直播间通知
func (this *BCServantObj) ExitBroRoom(ctx context.Context, req *bilin.ExitBroRoomReq) (resp *bilin.ExitBroRoomResp, err error) {
	const prefix = "ExitBroRoom "
	resp = &bilin.ExitBroRoomResp{Commonret: bccommon.SUCCESSMESSAGE}
	RoomId := req.Header.Roomid
	UserId := req.Header.Userid
	log.Debug(prefix+"begin", zap.Any("req", req))

	defer func(now time.Time) {
		httpmetrics.DefReport("ExitBroRoom", int64(resp.Commonret.Ret), now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	var room *entity.Room
	var user *entity.User
	if room, user, err = this.CommonCheckAuth(RoomId, UserId); err != nil {

		//从redis中删除用户
		service.RedisRemoveUser(RoomId, UserId)
		log.Warn("[-]"+prefix+"failed", zap.Any("resp", resp), zap.Any("RoomId", RoomId), zap.Any("UserId", UserId), zap.Any("err", err))

		return resp, nil
	}

	//从redis中删除用户
	service.RedisRemoveUser(RoomId, UserId)

	//直播间活动任务相关操作
	go thriftpool.Invoke(service.HotLineDataService, this.hotLineDataService, func(client interface{}) (err error) {
		c := client.(*hotline.DataServiceClient)
		intRes, err := c.Cancel(ctx, []*hotline.TaskReq{service.NewNotifyTask(user.UserID, room.Roomid, service.UserEnterRoomTask)})
		if err != nil {
			log.Error(prefix+"UserExitBroRoomTask", zap.Any("RoomId", RoomId), zap.Any("UserId", UserId), zap.Any("err", err), zap.Any("intRes", intRes))
		}

		log.Info(prefix+"UserExitBroRoomTask end", zap.Any("RoomId", RoomId), zap.Any("UserId", UserId), zap.Any("intRes", intRes))

		if user.Role == entity.ROLE_HOST {
			intHostRes, err2 := c.Cancel(ctx, []*hotline.TaskReq{service.NewNotifyTask(user.UserID, room.Roomid, service.HostStartLivingTask)})
			if err2 != nil {
				log.Error(prefix+"HostExitLivingTask", zap.Any("RoomId", RoomId), zap.Any("UserId", UserId), zap.Any("err", err2), zap.Any("intHostRes", intHostRes))
			}

			log.Info(prefix+"HostExitLivingTask end", zap.Any("RoomId", RoomId), zap.Any("UserId", UserId), zap.Any("intHostRes", intHostRes))
		}

		return
	})

	//统计用
	join_type := 2 //听众

	//如果是主持人 1,设置房间为非连线状态 ; 2,清空排麦列表, 3, 通知所有在麦序上的用户下麦, 4 初始化麦位状态
	if user.Role == entity.ROLE_HOST {
		service.RedisClearApplyMikeList(RoomId)

		if mikelist, _ := service.RedisGetOnMikeUserList(RoomId); mikelist != nil {
			for _, mikeuser := range mikelist {
				this.RemoveUserFromMike(ctx, room, mikeuser)
			}
		}

		//通知java，主持人退出房间
		var ComRet *common.ComRet
		err = thriftpool.Invoke(service.HotLineService, this.hotLine, func(client interface{}) (err error) {
			c := client.(*hotline.HotLineServiceClient)
			ComRet, err = c.HostUserLeave(ctx, int32(RoomId), int64(UserId))
			return
		})
		if err != nil || ComRet.Result_ != "success" {
			log.Error(prefix+"HostUserLeave", zap.Any("err", err), zap.Any("ComRet", ComRet))
		}

		collector.HostLeaveRoom(room, user.UserID)
		adapter.PushAllRoomInfoToRoom(room)

		//统计用
		join_type = 1
		collector.ExitRoomStat(RoomId, UserId, user.Role, int64(user.BeginJoinTime), time.Now().Unix(), join_type)

		log.Debug(prefix+"end", zap.Any("resp", resp))
		return
	}

	// 如果用户在麦序上
	if mikeUser, _ := service.RedisGetUserOnMike(RoomId, UserId); mikeUser != nil {
		this.RemoveUserFromMike(ctx, room, mikeUser)

		//如果是自动连线状态，需要从apply列表中选一个用户上麦
		this.supplementUserOnMike(ctx, room, mikeUser.MikeIndex)

		adapter.PushMikeListInfoToRoom(room)

		join_type = 1 //表示是嘉宾
	}

	// 如果正在排麦
	if userlist, e := service.RedisGetApplyMikeUserList(RoomId); e == nil {
		for _, item := range userlist {
			if item.UserID == UserId {
				service.RedisRemoveUserFromApplyMikeList(RoomId, UserId)

				//更新直播间内排麦用户数量
				adapter.PushBaseRoomInfoToRoom(room)
			}
		}
	}

	//通知直播间内用户列表变化增量，为了兼容坑爹的老版本
	adapter.PushUserListChangeToRoom(room, nil, []uint64{UserId})

	//广播用户列表变化
	adapter.PushUserListInfoToRoom(room)

	//统计用
	collector.ExitRoomStat(RoomId, UserId, user.Role, int64(user.BeginJoinTime), time.Now().Unix(), join_type)

	log.Debug(prefix+"end", zap.Any("RoomId", RoomId), zap.Any("UserId", UserId), zap.Any("resp", resp))
	return
}

//
// 直播间踢人
func (this *BCServantObj) KickUser(ctx context.Context, req *bilin.KickUserReq) (resp *bilin.KickUserResp, err error) {
	const prefix = "KickUser "
	resp = &bilin.KickUserResp{Commonret: bccommon.SUCCESSMESSAGE}
	roomId := req.Header.Roomid
	userId := req.Header.Userid
	log.Info(prefix+"begin", zap.Any("req", req))

	defer func(now time.Time) {
		httpmetrics.DefReport("KickUser", int64(resp.Commonret.Ret), now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	//查询用户和被踢人是否在频道内
	var optUser *entity.User
	var beKickedUser *entity.User
	if optUser, _ = service.RedisGetUser(roomId, userId); optUser == nil {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_ILLEGAL_MESSAGE, bccommon.COMMERRORILLEGALDESC)
		log.Error(prefix+"RedisGetUser not exist", zap.Any("roomId", roomId), zap.Any("userId", userId))
		return
	}

	if beKickedUser, _ = service.RedisGetUser(roomId, req.Kickeduserid); beKickedUser == nil {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_KICK_USER_NOT_IN_ROOM, fmt.Sprintf("用户不在直播间"))
		log.Error(prefix+"RedisGetUser not exist", zap.Any("roomId", roomId), zap.Any("userId", userId))
		return
	}

	// 检查权限
	if optUser.Role >= beKickedUser.Role {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_KICK_USER_NO_RIGHT, fmt.Sprintf("你没有该权限哦~"))
		log.Error(prefix+"permission denied", zap.Any("room", roomId), zap.Any("userId", userId), zap.Any("kickUserId", req.Kickeduserid))
		return
	}

	// 通知用户
	adapter.PushNotifyToUser(roomId, []int64{int64(req.Kickeduserid)}, &bilin.KickNotify{Optuserid: userId}, bilin.MinType_BC_NotifyRoomKickUser)

	//通知java
	var ComRet *common.ComRet
	err = thriftpool.Invoke(service.HotLineService, this.hotLine, func(client interface{}) (err error) {
		c := client.(*hotline.HotLineServiceClient)
		ComRet, err = c.KickUser(ctx, int32(roomId), int64(userId), int64(req.Kickeduserid), entity.ROLE_AUDIENCE)
		return
	})
	if err != nil || ComRet.Result_ != "success" {
		log.Error(prefix+"KickUser", zap.Any("err", err), zap.Any("result", ComRet.Result_), zap.Any("errorMsg", ComRet.ErrorMsg))
	}

	log.Info(prefix+"end", zap.Any("resp", resp), zap.Any("ret", ComRet))
	return
}

//
// 直播间禁麦和开麦  主持人抱听众上下麦
func (this *BCServantObj) MikeOperation(ctx context.Context, req *bilin.MikeOperationReq) (resp *bilin.MikeOperationResp, err error) {
	const prefix = "MikeOperation "
	resp = &bilin.MikeOperationResp{Commonret: bccommon.SUCCESSMESSAGE}
	roomId := req.Header.Roomid
	userId := req.Header.Userid
	log.Debug(prefix+"begin", zap.Any("req", req))

	defer func(now time.Time) {
		httpmetrics.DefReport("MikeOperation", int64(resp.Commonret.Ret), now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	var room *entity.Room
	var hoster *entity.User
	if room, hoster, err = this.CommonCheckAuth(roomId, userId); err != nil {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_ILLEGAL_MESSAGE, bccommon.COMMERRORILLEGALDESC)
		log.Error("[-]"+prefix+"failed", zap.Any("resp", resp), zap.Any("err", err))
		return
	}

	//检查用户权限
	if hoster.Role != entity.ROLE_HOST {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_MIKE_STAGE_NO_RIGHT, fmt.Sprintf("你没有该权限哦~"))
		log.Error(prefix, zap.Any("User", userId), zap.Any("Role", hoster.Role))
		return
	}

	//如果是上下麦操作，需要查询用户是否在频道内
	var mikeUser *entity.User
	if req.Opt == bilin.MikeOperationReq_UNMIKE || req.Opt == bilin.MikeOperationReq_ONMIKE {
		if mikeUser, _ = service.RedisGetUser(roomId, req.Userid); mikeUser == nil {
			resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_MIKE_OPREATION_USER_NOT_IN_ROOM, fmt.Sprintf("用户不在房间里"))
			log.Error(prefix+"MikeOperation failed", zap.Any("roomId", roomId), zap.Any("userId", userId))
			return
		}
	}

	switch req.Opt {
	case bilin.MikeOperationReq_UNMIKE:
		var retUser *entity.User
		if retUser, _ = service.RedisGetUserOnMike(roomId, mikeUser.UserID); retUser == nil {
			resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_MIKE_OPREATION_FAILED, fmt.Sprintf("用户不在麦上"))
			return resp, nil
		}

		//下麦
		this.RemoveUserFromMike(ctx, room, retUser)
		this.supplementUserOnMike(ctx, room, retUser.MikeIndex)
	case bilin.MikeOperationReq_ONMIKE:
		//检查用户是否在麦上
		if retUser, _ := service.RedisGetUserOnMike(roomId, mikeUser.UserID); retUser != nil {
			resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_MIKE_OPREATION_FAILED, fmt.Sprintf("您已经在麦上了哦~"))
			return resp, nil
		}

		//如果指定了麦位，检查麦位合理性，然后上麦
		if collector.CheckMikeNumberUsable(room, req.Mikeidx) {
			mikeUser.MikeIndex = req.Mikeidx
			this.AddUserToMike(ctx, room, mikeUser)
			//从排麦列表中删除等待的用户
			service.RedisRemoveUserFromApplyMikeList(room.Roomid, mikeUser.UserID)
		} else {
			if bilin.BaseRoomInfo_OPENMAIXU == room.Maixuswitch && req.Mikeidx != 0 {
				//麦位不可用，直接返回错误
				resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_MIKE_WHEAT_IS_USED, fmt.Sprintf("当前麦位已经有人了~"))
				log.Error(prefix, zap.Any("User", userId), zap.Any("resp", resp))
				return
			}

			//自动上麦
			if ret := this.autoOnMikeOperation(ctx, room, mikeUser); !ret {
				resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_MIKE_STAGE_IS_FULL, fmt.Sprintf("当前麦位已经坐满人了~"))
				log.Error(prefix, zap.Any("User", userId), zap.Any("resp", resp))
				return
			}

		}
	case bilin.MikeOperationReq_LOCKMIKE:
		if collector.CheckMikeNumberUsable(room, req.Mikeidx) {
			service.RedisLockUnlockMikeWheat(room.Roomid, req.Mikeidx, bilin.MikeInfo_LOCK)
		}
	case bilin.MikeOperationReq_UNLOCKMIKE:
		status, _ := service.RedisGetMikeWheatStatus(room.Roomid, req.Mikeidx)
		if status != int(bilin.MikeInfo_LOCK) {
			resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_MIKE_OPREATION_FAILED, fmt.Sprintf("当前麦位并未锁定~"))
			log.Error(prefix, zap.Any("User", userId), zap.Any("resp", resp))
			return
		}

		if !this.supplementUserOnMike(ctx, room, req.Mikeidx) {
			service.RedisLockUnlockMikeWheat(room.Roomid, req.Mikeidx, bilin.MikeInfo_EMPTY)
		}
	default:
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_ILLEGAL_MESSAGE, bccommon.COMMERRORILLEGALDESC)
		log.Error(prefix, zap.Any("User", userId), zap.Any("resp", resp))
		return
	}

	//通知直播间用户
	adapter.PushMikeListInfoToRoom(room)

	log.Debug(prefix+"end", zap.Any("resp", resp))
	return
}

//
// 更改直播间连线状态
func (this *BCServantObj) ChangeBroRoomLinkStatus(ctx context.Context, req *bilin.ChangeBroRoomLinkStatusReq) (resp *bilin.ChangeBroRoomLinkStatusResp, err error) {
	const prefix = "ChangeBroRoomLinkStatus "
	resp = &bilin.ChangeBroRoomLinkStatusResp{Commonret: bccommon.SUCCESSMESSAGE}
	roomId := req.Header.Roomid
	userId := req.Header.Userid
	log.Debug(prefix+"begin", zap.Any("req", req))

	defer func(now time.Time) {
		httpmetrics.DefReport("ChangeBroRoomLinkStatus", int64(resp.Commonret.Ret), now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	var room *entity.Room
	var user *entity.User
	if room, user, err = this.CommonCheckAuth(roomId, userId); err != nil {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_ILLEGAL_MESSAGE, bccommon.COMMERRORILLEGALDESC)
		log.Error("[-]"+prefix+"failed", zap.Any("resp", resp), zap.Any("err", err))
		return resp, nil
	}

	//检查用户权限
	if user.Role != entity.ROLE_HOST {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_CHANGE_ROOM_LINK_STATUS_NO_RIGHT, fmt.Sprintf("你没有该权限哦~"))
		log.Error(prefix, zap.Any("User", userId), zap.Any("Role", user.Role))
		return
	}

	//不需要做任何操作
	if room.GetLinkStatus() == req.Linkstatus {
		log.Warn(prefix, zap.Any("User", userId), zap.Any("status not change", req.Linkstatus))
		return
	}

	room.SetLinkstatus(req.Linkstatus)

	// 如果是关闭连线，需要清空排序列表
	if room.GetLinkStatus() == bilin.BaseRoomInfo_CLOSELINK {
		service.RedisClearApplyMikeList(roomId)
	} else if room.GetLinkStatus() == bilin.BaseRoomInfo_OPENLINK &&
		room.GetAutoLink() == bilin.BaseRoomInfo_OPENAUTOTOMIKE &&
		room.RoomType != bilin.BaseRoomInfo_ROOMTYPE_RADIO {
		if mikeChange := this.fillAllEmptyMikeWheat(ctx, room); mikeChange {
			adapter.PushMikeListInfoToRoom(room)
		}
	}

	//更新房间信息
	collector.StorageRoomInfo(room)

	adapter.PushBaseRoomInfoToRoom(room)

	log.Debug(prefix+"end", zap.Any("resp", resp))
	return
}

// 是否开启自动连麦
func (this *BCServantObj) ChangeBroRoomAutoToMikeStatus(ctx context.Context, req *bilin.ChangeBroRoomAutoToMikeStatusReq) (resp *bilin.ChangeBroRoomAutoToMikeStatusResp, err error) {
	const prefix = "ChangeBroRoomAutoToMikeStatus "
	resp = &bilin.ChangeBroRoomAutoToMikeStatusResp{Commonret: bccommon.SUCCESSMESSAGE}
	roomId := req.Header.Roomid
	userId := req.Header.Userid
	log.Debug(prefix+"begin", zap.Any("req", req))

	defer func(now time.Time) {
		httpmetrics.DefReport("BroRoomPraise", int64(resp.Commonret.Ret), now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	var room *entity.Room
	var user *entity.User
	if room, user, err = this.CommonCheckAuth(roomId, userId); err != nil {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_ILLEGAL_MESSAGE, bccommon.COMMERRORILLEGALDESC)
		log.Error("[-]"+prefix+"failed", zap.Any("resp", resp), zap.Any("err", err))
		return resp, nil
	}

	//电台模板不支持该操作
	if room.RoomType == bilin.BaseRoomInfo_ROOMTYPE_RADIO {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_CHANGE_ROOM_AUTO_TO_MIKE_FAILED, fmt.Sprintf("电台模板不支持该操作"))
		log.Error("[-]ChangeBroRoomAutoToMikeStatus failed", zap.Any("resp", resp))
		return resp, nil
	}

	//关闭连线状态，无法开启自动连麦
	if room.GetLinkStatus() == bilin.BaseRoomInfo_CLOSELINK {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_CHANGE_ROOM_AUTO_TO_MIKE_FAILED, fmt.Sprintf("服务器开小差了，再试试呗~"))
		log.Error(prefix, zap.Any("User", userId), zap.Any("Role", user.Role))
		return
	}

	//检查用户权限
	if user.Role != entity.ROLE_HOST {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_CHANGE_ROOM_AUTO_TO_MIKE_NO_RIGHT, fmt.Sprintf("服务器开小差了，再试试呗~"))
		log.Error(prefix, zap.Any("User", userId), zap.Any("Role", user.Role))
		return
	}

	//不需要做任何操作
	if room.GetAutoLink() == req.Autolink {
		log.Warn(prefix, zap.Any("User", userId), zap.Any("status not change", req.Autolink))
		return
	}
	room.SetAutoLink(req.Autolink)

	//电台模板不支持自动上麦
	if room.GetAutoLink() == bilin.BaseRoomInfo_OPENAUTOTOMIKE {
		if mikeChange := this.fillAllEmptyMikeWheat(ctx, room); mikeChange {
			adapter.PushMikeListInfoToRoom(room)
		}

	}

	//更新房间信息
	collector.StorageRoomInfo(room)

	adapter.PushBaseRoomInfoToRoom(room)

	log.Debug(prefix+"end", zap.Any("resp", resp))
	return
}

//
// 右下角点击,客户端聚合请求
func (this *BCServantObj) BroRoomPraise(ctx context.Context, req *bilin.BroRoomPraiseReq) (resp *bilin.BroRoomPraiseResp, err error) {
	const prefix = "BroRoomPraise "
	resp = &bilin.BroRoomPraiseResp{Commonret: bccommon.SUCCESSMESSAGE}
	roomId := req.Header.Roomid
	userId := req.Header.Userid
	log.Debug(prefix+"begin", zap.Any("req", req))

	defer func(now time.Time) {
		httpmetrics.DefReport("BroRoomPraise", int64(resp.Commonret.Ret), now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	var room *entity.Room
	if room, _, err = this.CommonCheckAuth(roomId, userId); err != nil {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_ILLEGAL_MESSAGE, bccommon.COMMERRORILLEGALDESC)
		log.Error("[-]"+prefix+"failed", zap.Any("resp", resp), zap.Any("err", err))
		return
	}
	if req.PraiseCount <= 0 {
		log.Warn(prefix+"end, PraiseCount less than 0, return!", zap.Any("resp", resp))
		return
	}
	//user.PraiseCount = req.PraiseCount
	//service.RedisAddUser(roomId, user)

	adapter.PushUserPraiseInfoToRoom(room, req.PraiseCount)

	log.Debug(prefix+"end", zap.Any("resp", resp))
	return
}

// 主播获取申请连线用户
func (this *BCServantObj) GetBroRoomPreparedAudience(ctx context.Context, req *bilin.GetBroRoomPreparedAudienceReq) (resp *bilin.GetBroRoomPreparedAudienceResp, err error) {
	const prefix = "GetBroRoomPreparedAudience "
	resp = &bilin.GetBroRoomPreparedAudienceResp{Commonret: bccommon.SUCCESSMESSAGE}
	roomId := req.Header.Roomid
	userId := req.Header.Userid
	log.Info(prefix+"begin", zap.Any("req", req))

	defer func(now time.Time) {
		httpmetrics.DefReport("GetBroRoomPreparedAudience", int64(resp.Commonret.Ret), now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	if _, _, err = this.CommonCheckAuth(roomId, userId); err != nil {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_ILLEGAL_MESSAGE, bccommon.COMMERRORILLEGALDESC)
		log.Error("[-]"+prefix+"failed", zap.Any("resp", resp), zap.Any("err", err))
		return
	}

	//任何人都可以获取申请连线用户列表
	//if user.Role != entity.ROLE_HOST {
	//	resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_GET_ROOM_LINK_AUDIENCE_NO_RIGHT, fmt.Sprintf("用户没有权限"))
	//	log.Error("[-]GetBroRoomPreparedAudience failed", zap.Any("resp", resp))
	//	return
	//}

	var userlist []*entity.User
	userlist, err = service.RedisGetApplyMikeUserList(roomId)
	if err != nil {
		log.Error(prefix, zap.Any("User", userId), zap.Any("err", err))
	}

	for _, item := range userlist {
		pushUser := collector.LocalUserToSendInfo(item)
		resp.Preparedusers = append(resp.Preparedusers, pushUser)
	}

	log.Debug(prefix+"end", zap.Any("resp", resp))
	return
}

// 主持人设置台上嘉宾静音  流程： 主持人设置-》服务器通知给该用户——》用户上报结果--》服务器修改状态--》发送广播
func (this *BCServantObj) MuteUser(ctx context.Context, req *bilin.MuteUserReq) (resp *bilin.MuteUserResp, err error) {
	const prefix = "MuteUser "
	resp = &bilin.MuteUserResp{Commonret: bccommon.SUCCESSMESSAGE}
	roomId := uint64(req.Header.Roomid)
	userId := uint64(req.Header.Userid)
	muteuserid := uint64(req.Muteuserid)
	log.Info(prefix+"begin", zap.Any("req", req))

	defer func(now time.Time) {
		httpmetrics.DefReport("MuteUser", int64(resp.Commonret.Ret), now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	var user *entity.User
	if _, user, err = this.CommonCheckAuth(roomId, userId); err != nil {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_ILLEGAL_MESSAGE, bccommon.COMMERRORILLEGALDESC)
		log.Error("[-]"+prefix+"failed", zap.Any("resp", resp), zap.Any("err", err))
		return resp, nil
	}

	if user.Role != entity.ROLE_HOST {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_MUTE_USER_NO_RIGHT, fmt.Sprintf("你没有该权限哦~"))
		log.Error("[-]MuteUser failed", zap.Any("resp", resp))
		return
	}

	var muteuser *entity.User
	if muteuser, _ = service.RedisGetUserOnMike(roomId, muteuserid); muteuser == nil {
		//用户不在麦序上
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_MUTE_USER_NOT_ON_MIKE, fmt.Sprintf("用户不在麦序上"))
		log.Error("[-]MuteUser failed", zap.Any("resp", resp))
		return
	}
	//对比麦序状态
	if muteuser.IsMuted == uint32(req.Opt) {
		log.Info("[-]MuteUser failed, IsMuted not change ", zap.Any("resp", resp))
		return
	}

	//通知用户设置静音相关操作
	adapter.PushNotifyToUser(roomId, []int64{int64(muteuserid)}, &bilin.MuteNotify{Optuserid: userId, Opt: uint32(req.Opt)}, bilin.MinType_BC_NotifyRoomAudienceMute)

	log.Info(prefix+"end", zap.Any("resp", resp))
	return
}

// 被静音的用户给服务器报结果
func (this *BCServantObj) MuteResult(ctx context.Context, req *bilin.MuteResultReq) (resp *bilin.MuteResultResp, err error) {
	const prefix = "MuteResult "
	resp = &bilin.MuteResultResp{Commonret: bccommon.SUCCESSMESSAGE}
	roomId := req.Header.Roomid
	userId := req.Header.Userid
	log.Info(prefix+"begin", zap.Any("req", req))

	defer func(now time.Time) {
		httpmetrics.DefReport("MuteResult", int64(resp.Commonret.Ret), now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	var room *entity.Room
	if room, _, err = this.CommonCheckAuth(roomId, userId); err != nil {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_ILLEGAL_MESSAGE, bccommon.COMMERRORILLEGALDESC)
		log.Error("[-]"+prefix+"failed", zap.Any("resp", resp), zap.Any("err", err))
		return resp, nil
	}

	var muteuser *entity.User
	if muteuser, _ = service.RedisGetUserOnMike(roomId, userId); muteuser == nil {
		//用户不在麦序上
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_MUTE_RESULT_NOT_ON_MIKE, fmt.Sprintf("用户不在麦上"))
		log.Error("[-]MuteUser failed", zap.Any("resp", resp))
		return
	}

	//设置mute属性, 没有抱麦操作，不需要发单播通知
	muteuser.IsMuted = uint32(req.Opt)
	if e := service.RedisAddUserToMike(roomId, muteuser); e != nil {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_MUTE_RESULT_FAILED, fmt.Sprintf("服务器开小差了，再试试呗~	"))
		log.Error("[-]MuteUser failed", zap.Any("resp", resp))
		return
	}

	//广播给频道用户，麦序信息变化
	adapter.PushMikeListInfoToRoom(room)

	log.Info(prefix+"end", zap.Any("resp", resp))
	return
}

// 主持人设置禁止公屏发言
func (this *BCServantObj) ForbiddenUser(ctx context.Context, req *bilin.ForbiddenUserReq) (resp *bilin.ForbiddenUserResp, err error) {
	const prefix = "ForbiddenUser "
	resp = &bilin.ForbiddenUserResp{Commonret: bccommon.SUCCESSMESSAGE}
	roomId := req.Header.Roomid
	userId := req.Header.Userid
	log.Info(prefix+"begin", zap.Any("req", req))

	defer func(now time.Time) {
		httpmetrics.DefReport("ForbiddenUser", int64(resp.Commonret.Ret), now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	var room *entity.Room
	var user *entity.User
	if room, user, err = this.CommonCheckAuth(roomId, userId); err != nil {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_ILLEGAL_MESSAGE, bccommon.COMMERRORILLEGALDESC)
		log.Error("[-]ForbiddenUser failed", zap.Any("resp", resp), zap.Any("err", err))
		return resp, nil
	}

	if user.Role != entity.ROLE_HOST {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_FORBIDDEN_USER_NO_RIGHT, fmt.Sprintf("你没有该权限哦~"))
		log.Error("[-]ForbiddenUser failed", zap.Any("resp", resp))
		return resp, nil
	}

	//exist, _ := service.RedisGetForbidenStatus(roomId, req.Forbiddenuserid)
	//if exist == req.Opt {
	//	log.Info("[-]ForbiddenUser user status not change", zap.Any("resp", resp))
	//	return
	//}

	service.RedisSetForbidenStatus(roomId, req.Forbiddenuserid, req.Opt)

	//通知用户被禁止公屏发言
	var optInt uint32
	if req.Opt {
		optInt = 1
	} else {
		optInt = 0
	}
	adapter.PushNotifyToUser(roomId, []int64{int64(req.Forbiddenuserid)}, &bilin.ForbiddenNotify{Optuserid: userId, Opt: optInt}, bilin.MinType_BC_NotifyUserBeForbidden)

	//通知频道内所有人  黑名单
	adapter.PushBlackListInfoToRoom(room)

	log.Info(prefix+"end", zap.Any("resp", resp))
	return
}

// 发送公屏消息
func (this *BCServantObj) SendRoomMessage(ctx context.Context, req *bilin.SendRoomMessageReq) (resp *bilin.SendRoomMessageResp, err error) {
	const prefix = "SendRoomMessage "
	resp = &bilin.SendRoomMessageResp{Commonret: bccommon.SUCCESSMESSAGE}
	roomId := req.Header.Roomid
	userId := req.Header.Userid
	log.Debug(prefix+"begin", zap.Any("req", req))

	defer func(now time.Time) {
		httpmetrics.DefReport(strings.TrimSpace(prefix), int64(resp.Commonret.Ret), now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	var room *entity.Room
	if room, _, err = this.CommonCheckAuth(roomId, userId); err != nil {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_ILLEGAL_MESSAGE, bccommon.COMMERRORILLEGALDESC)
		log.Error("[-]SendRoomMessage failed", zap.Any("resp", resp), zap.Any("err", err))
		return resp, nil
	}

	// 检查敏感词
	var checkRet int32
	err = thriftpool.Invoke(service.MsgFilterService, this.msgFilter, func(client interface{}) (err error) {
		c := client.(*bilin_msg_filter.MsgFilterClient)
		checkRet, err = c.CheckMsg(ctx, string(req.Data))
		return
	})
	if err != nil {
		log.Error(prefix+"CheckMsg", zap.Any("checkRet", checkRet), zap.Any("err", err))
		return
	}
	if checkRet == 2 { //返回2就是服务器判定不通过。
		log.Warn(prefix+"CheckMsg filter msg warnning", zap.Any("checkRet", checkRet), zap.Any("msg", req.Data))
		return
	}

	// 发送 min_type = 1008 类型的广播
	adapter.PushBroIMMsgToRoom(room, req.Data)

	log.Info(prefix+"end", zap.Any("resp", resp))
	return
}

// 客户端主动拉取房间全量信息
func (this *BCServantObj) GetAllRoomInfo(ctx context.Context, req *bilin.GetAllRoomInfoReq) (resp *bilin.GetAllRoomInfoResp, err error) {
	const prefix = "GetAllRoomInfo "
	resp = &bilin.GetAllRoomInfoResp{Commonret: bccommon.SUCCESSMESSAGE}
	roomid := req.Header.Roomid
	userid := req.Header.Userid
	log.Info(prefix+"begin", zap.Any("req", req))

	defer func(now time.Time) {
		httpmetrics.DefReport(strings.TrimSpace(prefix), int64(resp.Commonret.Ret), now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	if userid == 0 || roomid == 0 {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_ILLEGAL_MESSAGE, bccommon.COMMERRORILLEGALDESC)
		log.Error("[-]"+prefix+"failed", zap.Any("resp", resp))
		return resp, nil
	}

	var room *entity.Room
	if room = collector.GetRoomInfoByRoomId(roomid); room == nil {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_ENTER_ROOM_CLOSED, fmt.Sprintf("直播间已被关闭"))
		log.Error("[-]"+prefix+"failed", zap.Any("resp", resp))
		return resp, nil
	}

	resp.Allroominfo = collector.AllRoomInfo(room)
	log.Debug("[+]"+prefix+"success", zap.Any("req", req), zap.Any("resp", resp))
	return
}

func (this *BCServantObj) changeRoomType(ctx context.Context, room *entity.Room, newType bilin.BaseRoomInfo_ROOMTYPE) (err error) {
	const prefix = "changeRoomType "

	//清空排麦列表
	service.RedisClearApplyMikeList(room.Roomid)

	//如果是娱乐模板和电台模板切换，需要清空麦序，清空排麦列表
	if room.RoomType == bilin.BaseRoomInfo_ROOMTYPE_RADIO || newType == bilin.BaseRoomInfo_ROOMTYPE_RADIO {

		//更改麦位个数,只保留主持人麦位，其他麦位需要下麦
		mikelist, _ := service.RedisGetOnMikeUserList(room.Roomid)
		for _, item := range mikelist {
			if item.MikeIndex > 0 {
				this.RemoveUserFromMike(ctx, room, item)
			}
		}

		//电台模板和娱乐模板互切需要关闭K歌功能
		collector.CloseKaraoke(room)

		if newType == bilin.BaseRoomInfo_ROOMTYPE_RADIO {
			//娱乐模板切电台模板需要初始化麦位状态，因为之前可能会有一些锁定的麦位
			collector.InitMikeWheatInfo(room)

			//娱乐模板切电台模板需要停止寻友广播
			var ComRet *findfriendsbroadcast.OfflineFindFriendsBroadcastRet
			errThrift := thriftpool.Invoke(service.FindFriendsBroadcastService, this.findFriendBroadcast, func(client interface{}) (err error) {
				c := client.(*findfriendsbroadcast.FindFriendsBroadcastServiceClient)
				ComRet, err = c.OfflineFindFriendsBroadcastByUserIdList(ctx, []int64{int64(room.Owner)})
				return
			})
			if errThrift != nil || ComRet.Result_ != "success" {
				log.Error(prefix+"OfflineFindFriendsBroadcastByUserIdList", zap.Any("hostid", room.Owner), zap.Any("errThrift", errThrift), zap.Any("ComRet", ComRet))
			} else {
				log.Info(prefix+"OfflineFindFriendsBroadcastByUserIdList success", zap.Any("hostid", room.Owner), zap.Any("ComRet", ComRet))
			}
		}

		//电台模板和娱乐模板互切，都需要把连线开关打开
		room.LinkStatus = bilin.BaseRoomInfo_OPENLINK
	}

	if newType == bilin.BaseRoomInfo_ROOMTYPE_THREE {
		//4,5,6号麦位下麦
		mikelist, _ := service.RedisGetOnMikeUserList(room.Roomid)
		for _, item := range mikelist {
			if item.MikeIndex > 3 {
				this.RemoveUserFromMike(ctx, room, item)
			}
		}
	}

	//复用麦位信息,针对模板重新设置麦位
	collector.ReuseMikeWheatInfo(room, newType)

	//经产品沟通，切模板统一不做上麦操作，需要清麦序就可以了
	//if newType == bilin.BaseRoomInfo_ROOMTYPE_SIX {
	//	if room.GetLinkStatus() == bilin.BaseRoomInfo_OPENLINK && room.GetAutoLink() == bilin.BaseRoomInfo_OPENAUTOTOMIKE {
	//		if mikeChange := this.fillAllEmptyMikeWheat(ctx, room); mikeChange {
	//			room.RoomType = newType // 因为mikelistinfo里面需要带上新的模板信息，所以这里加上这个赋值
	//			adapter.PushMikeListInfoToRoom(room)
	//		}
	//	}
	//}

	log.Info(prefix, zap.Any("room", room))
	return err
}

// 主持人切换模板   BaseRoomInfo.ROOMTYPE
func (this *BCServantObj) ChangeBroRoomType(ctx context.Context, req *bilin.ChangeBroRoomTypeReq) (resp *bilin.ChangeBroRoomTypeResp, err error) {
	const prefix = "ChangeBroRoomType "
	resp = &bilin.ChangeBroRoomTypeResp{Commonret: bccommon.SUCCESSMESSAGE}
	roomid := req.Header.Roomid
	userid := req.Header.Userid
	log.Info(prefix+"begin", zap.Any("req", req))

	defer func(now time.Time) {
		httpmetrics.DefReport(strings.TrimSpace(prefix), int64(resp.Commonret.Ret), now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	var room *entity.Room
	var user *entity.User
	if room, user, err = this.CommonCheckAuth(roomid, userid); err != nil {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_ILLEGAL_MESSAGE, bccommon.COMMERRORILLEGALDESC)
		log.Error("[-]ChangeBroRoomType failed", zap.Any("resp", resp))
		return resp, nil
	}

	//检查用户权限
	if user.Role != entity.ROLE_HOST {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_CHANGE_ROOM_TYPE_NO_RIGHT, fmt.Sprintf("你没有该权限哦~"))
		log.Error(prefix, zap.Any("User", userid), zap.Any("Role", user.Role))
		return
	}

	if valid := collector.IsValidRoomType(req.Roomtype); !valid {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_CHANGE_ROOM_TYPE_INVALID, fmt.Sprintf("直播模板不可用"))
		log.Error("[-]ChangeBroRoomType failed", zap.Any("resp", resp))
		return resp, nil
	}

	//切1+6模板时需要开关为开才行
	if req.Roomtype == bilin.BaseRoomInfo_ROOMTYPE_SIX && room.Maixuswitch == bilin.BaseRoomInfo_CLOSEMAIXU {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_CHANGE_ROOM_TYPE_INVALID, fmt.Sprintf("直播模板不可用"))
		log.Error("[-]ChangeBroRoomType failed", zap.Any("resp", resp))
		return resp, nil
	}

	if room.RoomType == req.Roomtype {
		log.Info("[-]ChangeBroRoomType room.RoomType not change", zap.Any("resp", resp))
		return
	}

	log.Info(prefix, zap.Any("from roomType", room.RoomType), zap.Any("to roomType", req.Roomtype))

	//先下麦
	this.changeRoomType(ctx, room, req.Roomtype)

	room.RoomType = req.Roomtype
	collector.StorageRoomInfo(room)

	//通知所有人模板切换
	adapter.PushBaseRoomInfoToRoom(room)
	adapter.PushMikeListInfoToRoom(room)
	log.Info("[+]"+prefix+"success", zap.Any("resp", resp))
	return
}

// 客户端请求分页列表信息
func (this *BCServantObj) GetBroRoomUsersByPage(ctx context.Context, req *bilin.GetBroRoomUsersByPageReq) (resp *bilin.GetBroRoomUsersByPageResp, err error) {
	const prefix = "GetBroRoomUsersByPage "
	resp = &bilin.GetBroRoomUsersByPageResp{Commonret: bccommon.SUCCESSMESSAGE}
	roomid := req.Header.Roomid
	log.Debug(prefix+"begin", zap.Any("req", req))

	defer func(now time.Time) {
		httpmetrics.DefReport(strings.TrimSpace(prefix), int64(resp.Commonret.Ret), now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	resp.Audienceusers = collector.GetRoomUsersByPage(roomid, req.Pagenumber)
	log.Debug("[+]"+prefix+"success", zap.Any("resp", resp))
	return
}

// 运营后台通知禁止某个直播间
func (this *BCServantObj) ForbiddenRoom(ctx context.Context, req *bilin.ForbiddenRoomReq) (resp *bilin.ForbiddenRoomResp, err error) {
	const prefix = "ForbiddenRoom "
	resp = &bilin.ForbiddenRoomResp{Commonret: bccommon.SUCCESSMESSAGE}
	roomid := req.Header.Roomid
	log.Debug(prefix+"begin", zap.Any("req", req))

	defer func(now time.Time) {
		httpmetrics.DefReport(strings.TrimSpace(prefix), int64(resp.Commonret.Ret), now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	var room *entity.Room
	if room = collector.GetRoomInfoByRoomId(roomid); room == nil {
		return resp, nil
	}

	room.Status = bilin.BaseRoomInfo_FORBIDDEN
	collector.StorageRoomInfo(room)

	adapter.PushRoomClosedNotifyToRoom(room, req.Hostnotifytext, req.Audiencenotifytext)

	//临时处理6.4防敏感
	{
		// 通知用户
		adapter.PushNotifyToUser(roomid, []int64{int64(req.Header.Userid)}, &bilin.KickNotify{Optuserid: 0}, bilin.MinType_BC_NotifyRoomKickUser)

		//通知java
		var ComRet *common.ComRet
		err = thriftpool.Invoke(service.HotLineService, this.hotLine, func(client interface{}) (err error) {
			c := client.(*hotline.HotLineServiceClient)
			ComRet, err = c.KickUser(ctx, int32(roomid), int64(0), int64(req.Header.Userid), entity.ROLE_HOST)
			return
		})
		if err != nil || ComRet.Result_ != "success" {
			log.Error(prefix+"KickUser", zap.Any("err", err), zap.Any("result", ComRet.Result_), zap.Any("errorMsg", ComRet.ErrorMsg))
		}
	}

	log.Debug("[+]"+prefix+"success", zap.Any("resp", resp))
	return
}

func (this *BCServantObj) GetRoomPrivilegeInfo(ctx context.Context, req *bilin.RoomPrivilegeInfoReq) (resp *bilin.RoomPrivilegeInfoResp, err error) {
	const prefix = "GetRoomPrivilegeInfo "
	roomid := req.Header.Roomid
	userid := req.Header.Userid
	resp = &bilin.RoomPrivilegeInfoResp{Commonret: bccommon.SUCCESSMESSAGE}
	log.Debug(prefix+"begin", zap.Any("req", req))

	defer func(now time.Time) {
		httpmetrics.DefReport(strings.TrimSpace(prefix), int64(resp.Commonret.Ret), now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	var room *entity.Room
	if room, _, err = this.CommonCheckAuth(roomid, userid); err != nil {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_ILLEGAL_MESSAGE, bccommon.COMMERRORILLEGALDESC)
		log.Error("[-]ChangeBroRoomType failed", zap.Any("resp", resp))
		return resp, nil
	}

	//用户特权信息
	resp.Privilegeinfo = this.getUserPrivilegeInfo(ctx, room, room.Owner, userid)

	log.Debug("[+]"+prefix+"success", zap.Any("resp", resp))
	return
}

// 主持人清空排麦列表 只有关闭自动连麦时才能清空
func (this *BCServantObj) ClearRoomPreparedAudience(ctx context.Context, req *bilin.ClearRoomPreparedAudienceReq) (resp *bilin.ClearRoomPreparedAudienceResp, err error) {
	const prefix = "ClearRoomPreparedAudience "
	roomid := req.Header.Roomid
	userid := req.Header.Userid
	resp = &bilin.ClearRoomPreparedAudienceResp{Commonret: bccommon.SUCCESSMESSAGE}
	log.Debug(prefix+"begin", zap.Any("req", req))

	defer func(now time.Time) {
		httpmetrics.DefReport(strings.TrimSpace(prefix), int64(resp.Commonret.Ret), now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	var room *entity.Room
	var user *entity.User
	if room, user, err = this.CommonCheckAuth(roomid, userid); err != nil {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_ILLEGAL_MESSAGE, bccommon.COMMERRORILLEGALDESC)
		log.Error("[-]ChangeBroRoomType failed", zap.Any("resp", resp))
		return resp, nil
	}

	//检查用户权限
	if user.Role != entity.ROLE_HOST {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_CHANGE_ROOM_TYPE_NO_RIGHT, fmt.Sprintf("你没有该权限哦~"))
		log.Error(prefix, zap.Any("User", userid), zap.Any("Role", user.Role))
		return
	}

	var multicastUids []int64
	if userlist, e := service.RedisGetApplyMikeUserList(roomid); e == nil {
		for _, item := range userlist {
			multicastUids = append(multicastUids, int64(item.UserID))
		}
	}

	//发送push消息,多播
	adapter.PushNotifyToUser(room.Roomid, multicastUids, &bilin.ClearRoomPreparedAudienceNotify{}, bilin.MinType_BC_NotifyRoomClearPreparedAudience)

	//清空排麦列表
	service.RedisClearApplyMikeList(room.Roomid)

	//发广播
	adapter.PushBaseRoomInfoToRoom(room)

	log.Debug("[+]"+prefix+"success", zap.Any("resp", resp))
	return
}
