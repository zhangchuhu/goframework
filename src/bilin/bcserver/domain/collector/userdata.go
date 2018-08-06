package collector

import (
	"bilin/bcserver/bccommon"
	"bilin/bcserver/domain/entity"
	"bilin/bcserver/domain/service"
	"bilin/protocol"
	"bilin/thrift/gen-go/hotline"
	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"fmt"
	"time"
)

//专门处理主持人进入事件
func HostEnterRoom(room *entity.Room, host uint64, joinHotLineRet *hotline.JoinHotLineRet) (err error) {
	const prefix = "HostEnterRoom "
	log.Debug(prefix+"begin", zap.Any("roomid", room.Roomid), zap.Any("host", host))

	{
		//set roominfo
		room.Title = joinHotLineRet.GetTitle()
		room.RoomType2 = joinHotLineRet.GetRoomType()
		room.RoomCategoryID = joinHotLineRet.GetRoomCategoryId()
		room.HostBilinID = joinHotLineRet.GetHostBilinId()
		room.RoomPendantLevel = joinHotLineRet.GetRoomPendantLevel()
		room.Status = bilin.BaseRoomInfo_OPEN
		//如果主持人已经在房间里，考虑到某些特殊情况，比如主持人杀进程，重启，然后重进房间
		ifOnMike, _ := service.RedisIfUserOnMike(room.Roomid, host)
		if !ifOnMike && room.RoomType2 != service.OFFICAIL_ROOM {
			room.RoomType = bilin.BaseRoomInfo_ROOMTYPE_THREE //只有主持人真正的进房间才会设置模板为1+3
			room.StartTime = uint64(time.Now().Unix())
		}

		room.LinkStatus = bilin.BaseRoomInfo_OPENLINK
		room.Owner = host
		room.From = "" //陆续清掉oldbcserver里面的缓存信息
		room.Maixuswitch = bilin.BaseRoomInfo_OPENMAIXU
		room.LockStatus = BizQueryRoomLockStatus(room)
		//主持人进房间开播需要写流水,注意先后顺序，因为写流水之后会返回一个流水ID，需要存到redis中
		MysqlStorageLivingRecordInfo(room)
		StorageRoomInfo(room)
	}

	//设置主播和房间id的映射
	service.RedisSetUidToRoomId(host, room.Roomid)

	//java那边有个坑爹的要求，主持人退房间之后还要给定时他发hostleavetoolong才行....
	//如果主持人进入了房间，需要先去掉这个timer
	err = service.RedisRemoveHostLeaveTooLongTask(room.Roomid)
	if err != nil {
		log.Error(prefix+"RedisHostEnterRoom", zap.Any("roomid", room.Roomid), zap.Any("host", host))
		return
	}

	log.Debug(prefix+"end", zap.Any("roomid", room.Roomid), zap.Any("host", host))
	return
}

//java要求，主持人退出房间，不管怎样，都要发一个timer通知hostleavetoolong
func HostLeaveRoom(room *entity.Room, host uint64) (err error) {
	const prefix = "HostLeaveRoom "
	log.Debug(prefix+"begin", zap.Any("roomid", room.Roomid), zap.Any("host", host))

	{
		//清空麦位信息
		InitMikeWheatInfo(room)
		room.LinkStatus = bilin.BaseRoomInfo_CLOSELINK

		//关闭K歌功能
		CloseKaraoke(room)

		//房间信息变化，需要改写redis
		StorageRoomInfo(room)

		//主持人关播时需要更新流水
		MysqlUpdateLivingRecordInfo(room)
	}

	err = service.RedisAddHostLeaveTooLongTask(room, host)
	if err != nil {
		log.Error(prefix+"RedisHostLeaveRoom", zap.Any("roomid", room.Roomid), zap.Any("host", host))
		return
	}

	log.Debug(prefix+"end", zap.Any("roomid", room.Roomid), zap.Any("host", host))
	return
}

func InitUserByJavaResult(RoomId uint64, UserId uint64, joinHotLineRet *hotline.JoinHotLineRet, MapExtend map[uint32]string) (retUser *entity.User, commRet *bilin.CommonRetInfo) {
	const prefix = "InitUserByJavaResult "

	//如果用户已经在房间内，重新赋值一下角色信息
	var redisErr error
	if retUser, redisErr = service.RedisGetUser(RoomId, UserId); redisErr != nil {
		return retUser, bccommon.UserDefinedFailed(bilin.CommonRetInfo_ENTER_BAD_NETWORK, fmt.Sprintf("服务器开小差了，再试试呗~"))
	}
	if retUser != nil {
		retUser.Role = uint32(joinHotLineRet.Status)
		service.RedisAddUser(RoomId, retUser)
		return retUser, bccommon.UserDefinedFailed(bilin.CommonRetInfo_ENTER_ROOM_ALREADY_IN_ROOM, fmt.Sprintf("用户已经在房间"))
	}

	now := time.Now().Unix()
	retUser = &entity.User{
		RoomID:         RoomId,
		UserID:         UserId,
		Role:           uint32(joinHotLineRet.Status),
		Status:         entity.StatusUserJoined,
		BeginJoinTime:  uint64(now),
		NickName:       *joinHotLineRet.Nickname,
		AvatarURL:      *joinHotLineRet.HeaderUrl,
		IsMuted:        0,
		EnterBeginTime: uint64(now),
		LinkBeginTime:  0,
		Sex:            joinHotLineRet.Sex,
		Age:            joinHotLineRet.Age,
		CityName:       *joinHotLineRet.CityName,
		PraiseCount:    0,
		Version:        MapExtend[uint32(bilin.Header_VERSION)],
		Signature:      joinHotLineRet.GetSign(),
	}

	if redisErr = service.RedisAddUser(RoomId, retUser); redisErr != nil {
		log.Warn(prefix+"RedisAddUser error", zap.Any("roomid", RoomId), zap.Any("UserId", UserId))
		return retUser, bccommon.UserDefinedFailed(bilin.CommonRetInfo_ENTER_BAD_NETWORK, fmt.Sprintf("服务器开小差了，再试试呗~"))
	}

	//统计用
	EnterRoomStat(retUser.RoomID, retUser.UserID, retUser.Role)

	log.Debug(prefix+"end", zap.Any("roomid", RoomId), zap.Any("retUser", retUser), zap.Any("joinHotLineRet", *joinHotLineRet))
	return retUser, bccommon.SUCCESSMESSAGE
}

//主持人开播时需要记录开播流水
func MysqlStorageLivingRecordInfo(room *entity.Room) {
	service.MysqlInsertLivingData(room)
}

//主持人关播时需要更新流水
func MysqlUpdateLivingRecordInfo(room *entity.Room) {
	service.MysqlUpdateLivingData(room)
}
