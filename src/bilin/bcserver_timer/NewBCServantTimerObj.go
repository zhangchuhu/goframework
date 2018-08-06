package main

import (
	"bilin/bcserver/bccommon"
	"bilin/bcserver/domain/adapter"
	"bilin/bcserver/domain/collector"
	"bilin/bcserver/domain/entity"
	"bilin/bcserver/domain/service"
	"bilin/common/thriftpool"
	"bilin/protocol"
	"bilin/thrift/gen-go/common"
	"bilin/thrift/gen-go/hotline"
	"bilin/thrift/gen-go/user"
	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/kissgo/httpmetrics"
	"code.yy.com/yytars/goframework/tars/servant"
	"context"
	"fmt"
	"strings"
	"time"
)

const (
	PingTimeOut         = 60
	HandleSuccess       = 0
	ErrorHscanRoomList  = 1
	ErrorGetAllPingTime = 2
	ErrorExitBroRoom    = 3

	ErrorGetUserCount             = 4
	ErrorFreshData                = 5
	ErrorGetHostLeaveTooLongTasks = 6
	ErrorHostUserOfflineTooLong   = 7

	OLD_BC_SERVER = "old_bc_server"
)

// BCServantObj 包含所有直播间相关的操作
type BCServantTimerObj struct {
	bcClient       bilin.BCServantClient
	hotLine        thriftpool.Pool
	userFans       thriftpool.Pool
	mapTimerHandle map[time.Duration]interface{}
}

//定时清空超时用户
func TimerHandleUserTimeOut(obj *BCServantTimerObj, roomList []*entity.Room) {
	const prefix = "TimerHandleUserTimeOut "
	log.Debug(prefix + "begin")

	retCode := HandleSuccess
	defer func(now time.Time) {
		httpmetrics.DefReport(strings.TrimSpace(prefix), int64(retCode), now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	//foreach room
	for _, room := range roomList {
		pingList, err := service.RedisGetAllPingTimeByRoomid(room.Roomid) //一个房间里面就100来人，这里目前不需要优化
		if err != nil {
			retCode = ErrorGetAllPingTime
			log.Error(prefix+"redis.RedisGetAllPingTimeByRoomid", zap.Any("err", err))
			continue
		}

		//foreach user
		now := uint64(time.Now().Unix())
		for userid, lastPingTime := range pingList {
			if now-lastPingTime > PingTimeOut {
				//rpc call
				resp, err := obj.bcClient.ExitBroRoom(context.TODO(), &bilin.ExitBroRoomReq{
					Header: &bilin.Header{
						Roomid: room.Roomid,
						Userid: userid,
					},
				})

				if err != nil || resp.Commonret.Ret != bilin.CommonRetInfo_RETCODE_SUCCEED {
					retCode = ErrorExitBroRoom
					log.Error(prefix+"ExitBroRoom", zap.Any("roomid", room.Roomid), zap.Any("userid", userid), zap.Any("err", err))
					continue
				}

				log.Info(prefix+"ExitBroRoom", zap.Any("roomid", room.Roomid), zap.Any("userid", userid), zap.Any("now", now), zap.Any("lastPingTime", lastPingTime), zap.Any("resp", resp))
			}
		}
	}

}

//定时给java通报直播间用户数
func TimerReportRoomStatis(obj *BCServantTimerObj, roomList []*entity.Room) {
	const prefix = "TimerReportRoomStatis "
	log.Debug(prefix + "begin")

	retCode := HandleSuccess
	defer func(now time.Time) {
		httpmetrics.DefReport(strings.TrimSpace(prefix), int64(retCode), now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	//foreach room
	for _, room := range roomList {
		if room.Status == bilin.BaseRoomInfo_CLOSED { //只处理新的server数据，因为old bc server没有从db中恢复的功能
			if room.From != OLD_BC_SERVER && room.EndTime != 0 && uint64(time.Now().Unix())-room.EndTime >= 1*24*60*60 { //超过1天未使用的房间需要清理
				collector.StorageRoomInfo(room)
				service.RedisRemoveRoom(room.Roomid)

				log.Info(prefix+"Remove room from redis ", zap.Any("room", room))
			}
			continue
		}

		count, err := service.RedisGetUserCount(room.Roomid)
		if err != nil {
			retCode = ErrorGetUserCount
			log.Error(prefix+"redis.RedisGetUserCount", zap.Any("err", err))
		} else {
			//如果房间人数为0，需要设置房间状态为close
			if count == 0 {
				room.Status = bilin.BaseRoomInfo_CLOSED
				room.EndTime = uint64(time.Now().Unix())
				collector.StorageRoomInfo(room)
			}

			//通知java
			var ComRet *common.ComRet
			err = thriftpool.Invoke(service.HotLineService, obj.hotLine, func(client interface{}) (err error) {
				c := client.(*hotline.HotLineServiceClient)
				ComRet, err = c.FreshData(context.TODO(), int32(room.Roomid), count, 0, 0, count, count, 0)
				return
			})
			if err != nil || ComRet.Result_ != "success" {
				retCode = ErrorFreshData
				log.Error(prefix+"FreshData", zap.Any("err", err), zap.Any("ComRet", ComRet))
				continue
			}

			log.Debug(prefix+"FreshData", zap.Any("roomid", room.Roomid), zap.Any("count", count))
		}
	}
}

//定时给所有直播间推送MikeListInfo
func TimerPushAllRoomInfoToRoom(obj *BCServantTimerObj, roomList []*entity.Room) {
	const prefix = "TimerPushAllRoomInfoToRoom "
	log.Debug(prefix + "begin")

	retCode := HandleSuccess
	defer func(now time.Time) {
		httpmetrics.DefReport(strings.TrimSpace(prefix), int64(retCode), now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	//foreach room
	for _, room := range roomList {
		if room.Owner == 0 || room.Status == bilin.BaseRoomInfo_CLOSED {
			continue
		}
		//经与客户端协调
		//后台定时向客户端推送直播间信息的操作，从流量和性能开销考虑，需要优化一下：
		//1、降低频率，建议20s以上
		//2、只推送麦上用户的数据
		adapter.PushMikeListInfoToRoom(room)
		log.Debug(prefix+"PushBaseRoomInfoToRoom  PushMikeListInfoToRoom", zap.Any("room", room))
	}

}

func TimerQueryAttentionMeCount(obj *BCServantTimerObj, roomList []*entity.Room) {
	const prefix = "TimerQueryAttentionMeCount "
	log.Debug(prefix + "begin")

	retCode := HandleSuccess
	defer func(now time.Time) {
		httpmetrics.DefReport(strings.TrimSpace(prefix), int64(retCode), now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	//foreach room
	var userlist []int64
	ownerMap := make(map[uint64]uint64)
	for _, room := range roomList {
		if room.Owner == 0 || room.Status == bilin.BaseRoomInfo_CLOSED {
			continue
		}

		ownerMap[room.Roomid] = room.Owner
		userlist = append(userlist, int64(room.Owner))
	}

	if userlist != nil {
		//从java那边查询粉丝数
		var ComRet *user.QueryAttentionMeCountRet
		err := thriftpool.Invoke(service.UserServervice, obj.userFans, func(client interface{}) (err error) {
			c := client.(*user.UserServiceClient)
			ComRet, err = c.QueryAttentionMeCount(context.TODO(), userlist)
			return
		})
		if err != nil || ComRet.Result_ != "success" {
			retCode = ErrorFreshData
			log.Error(prefix+"QueryAttentionMeCount", zap.Any("err", err), zap.Any("ComRet", ComRet))
			return
		}

		//更新用户粉丝数
		for _, owner := range ownerMap {
			service.RedisSetUserFansCount(owner, uint32(ComRet.AttentionMeCountMap[int64(owner)]))
			log.Debug(prefix+"RedisSetUserFansCount ", zap.Any("owner", owner), zap.Any("FansCount", ComRet.AttentionMeCountMap[int64(owner)]))
		}
	}
}

func TimerCheckMikeErrors(obj *BCServantTimerObj, roomList []*entity.Room) {
	const prefix = "TimerCheckMikeErrors "
	log.Debug(prefix + "begin")

	retCode := HandleSuccess
	defer func(now time.Time) {
		httpmetrics.DefReport(strings.TrimSpace(prefix), int64(retCode), now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	//foreach room
	for _, room := range roomList {
		if room.Owner == 0 || room.Status == bilin.BaseRoomInfo_CLOSED {
			continue
		}

		collector.CheckMikeErrors(room)
	}
}

func TimerHandlerManager(obj *BCServantTimerObj, interval time.Duration) {
	const prefix = "TimerHandlerManager "

	for {
		log.Debug(prefix + "begin")

		retCode := HandleSuccess
		defer func(now time.Time) {
			httpmetrics.DefReport(strings.TrimSpace(prefix), int64(retCode), now, bccommon.SuccessOrFailedFun)
		}(time.Now())

		var HandlerArray []interface{}
		for key, value := range obj.mapTimerHandle {
			if time.Now().Unix()%int64(key/time.Second) == 0 {
				HandlerArray = append(HandlerArray, value)
			}
		}

		var cursor uint64 = 0
		var total_room_count int = 0
		if len(HandlerArray) == 0 {
			goto SLEEP
		}

		log.Info(prefix, zap.Any("HandlerArray", HandlerArray))
		for {
			var err error
			var roomList []*entity.Room
			roomList, cursor, err = service.RedisHscanRoomList(cursor)
			if err != nil {
				retCode = ErrorHscanRoomList
				log.Error(prefix+"redis.RedisGetRoomIdList", zap.Any("err", err), zap.Any("roomlist", roomList))
				break
			}

			log.Debug(prefix, zap.Any("cursor", cursor), zap.Any("roomList_len", len(roomList)))

			total_room_count += len(roomList)
			for _, curFunc := range HandlerArray {
				go curFunc.(func(*BCServantTimerObj, []*entity.Room))(obj, roomList)
			}

			if cursor == 0 { //读完数据了
				log.Info(prefix+"finished loop, waitting next time", zap.Any("cursor", cursor), zap.Any("total_room_count", total_room_count))
				break
			}
		}

	SLEEP:
		time.Sleep(interval)
	}
}

func TimerReportHostLeaveTooLong(obj *BCServantTimerObj, interval time.Duration) {
	const prefix = "TimerReportHostLeaveTooLong "
	for {
		log.Info(prefix + "begin")

		retCode := HandleSuccess
		defer func(now time.Time) {
			httpmetrics.DefReport(strings.TrimSpace(prefix), int64(retCode), now, bccommon.SuccessOrFailedFun)
		}(time.Now())

		tasks, err := service.RedisGetAllHostLeaveTooLongTasks()
		if err != nil {
			retCode = ErrorGetHostLeaveTooLongTasks
			log.Error(prefix+"redis.RedisGetAllHostLeaveTooLongTasks", zap.Any("err", err))
			goto SLEEP
		}

		for _, item := range tasks {
			if item.RoomType == service.OFFICAIL_ROOM {
				room, _ := service.RedisGetRoomInfo(item.RoomId)
				item.HostId = room.Owner
			}
			//先查询主播是否在直播间中
			exist, _ := service.RedisIfUserOnMike(item.RoomId, item.HostId)
			if exist {
				//删除数据
				service.RedisRemoveHostLeaveTooLongTask(item.RoomId)
				continue
			}
			if uint64(time.Now().Unix())-item.LeaveTime >= 300 { //主播不在直播间超过300秒
				//通知java
				var ComRet *common.ComRet
				err = thriftpool.Invoke(service.HotLineService, obj.hotLine, func(client interface{}) (err error) {
					c := client.(*hotline.HotLineServiceClient)
					ComRet, err = c.HostUserOfflineTooLong(context.TODO(), int32(item.RoomId), int64(item.HostId))
					return
				})
				if err != nil || ComRet.Result_ != "success" {
					retCode = ErrorHostUserOfflineTooLong
					log.Error(prefix+"HostUserOfflineTooLong", zap.Any("err", err), zap.Any("ComRet", ComRet))
					continue
				}

				log.Info(prefix+"success process report HostUserOfflineTooLong", zap.Any("item", item))

				//删除数据
				service.RedisRemoveHostLeaveTooLongTask(item.RoomId)
			}

		}

	SLEEP:
		time.Sleep(interval)
	}
}

// NewBCServantObj 被main调用，初始化
func NewBCServantTimerObj() *BCServantTimerObj {
	service.RedisInit()
	service.MysqlInit()
	hotLine, err := thriftpool.NewChannelPool(0, 1000, service.CreateHotLineServiceConn)
	if err != nil {
		log.Panic("can not create thrift connection pool hotLine", zap.Any("err", err))
	}

	userFans, err := thriftpool.NewChannelPool(0, 1000, service.CreateUserServiceConn)
	if err != nil {
		log.Panic("can not create thrift connection pool UserService", zap.Any("err", err))
	}

	comm := servant.NewPbCommunicator()
	objName := fmt.Sprintf("bilin.bcserver2.BCServantObj")
	s := &BCServantTimerObj{
		bcClient:       bilin.NewBCServantClient(objName, comm),
		hotLine:        hotLine,
		userFans:       userFans,
		mapTimerHandle: make(map[time.Duration]interface{}),
	}

	//注册timerhandler
	s.mapTimerHandle[10*time.Second] = TimerHandleUserTimeOut
	s.mapTimerHandle[6*time.Second] = TimerReportRoomStatis
	s.mapTimerHandle[20*time.Second] = TimerPushAllRoomInfoToRoom
	s.mapTimerHandle[3*time.Second] = TimerQueryAttentionMeCount
	s.mapTimerHandle[5*time.Second] = TimerCheckMikeErrors

	go TimerHandlerManager(s, 1*time.Second)

	go TimerReportHostLeaveTooLong(s, 5*time.Second)

	go thriftpool.Ping(service.HotLineService, s.hotLine, func(client interface{}) (err error) {
		c := client.(*hotline.HotLineServiceClient)
		_, err = c.Ping(context.TODO())
		return
	}, 5*time.Second)

	go thriftpool.Ping(service.UserServervice, s.userFans, func(client interface{}) (err error) {
		c := client.(*user.UserServiceClient)
		_, err = c.Ping(context.TODO())
		return
	}, 5*time.Second)

	return s
}
