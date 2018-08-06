package collector

import (
	"bilin/bcserver/bccommon"
	"bilin/bcserver/domain/entity"
	"bilin/bcserver/domain/service"
	"bilin/protocol"
	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/tars"
	"context"
	"encoding/json"
)

type KafkaMessage struct {
	Roomid     uint64 `json:"roomid"`
	LockStatus uint32 `json:"lock_status"`
	Owner      uint64 `json:"owner"`
}

var bizRoomCenterClient bilin.BizRoomCenterServantClient

func init() {
	comm := tars.NewCommunicator()
	bizRoomCenterClient = bilin.NewBizRoomCenterServantClient("bilin.bizroomcenter.BizRoomCenterPbObj", comm)
}

func BizAuth(room *entity.Room, UserId uint64, roompwd string) (result bool, err error) {
	const prefix = "BizAuth "

	result = true
	//主播不需要鉴权
	if room.LockStatus == entity.UNLOCKROOM || room.Owner == UserId {
		return
	}

	resp, err := bizRoomCenterClient.GetBizRoomInfo(context.TODO(), &bilin.GetBizRoomInfoReq{Roomid: room.Roomid})
	if err != nil || resp.Commonret.Ret != bccommon.SUCCESSMESSAGE.Ret {
		result = false
		log.Error(prefix+"failed", zap.Any("roomid", room.Roomid), zap.Any("UserId", UserId))
		return
	}

	if resp.Bizroominfo == nil || bccommon.IsNil(resp.Bizroominfo) {
		return
	}

	if resp.Bizroominfo.Roompwd != roompwd {
		result = false
		log.Error(prefix+"failed", zap.Any("roomid", room.Roomid), zap.Any("UserId", UserId), zap.Any("resp", resp), zap.Any("roompwd", roompwd))
		return
	}

	log.Info(prefix+"success", zap.Any("roomid", room.Roomid), zap.Any("UserId", UserId))
	return
}

func BizQueryRoomLockStatus(room *entity.Room) (status uint32) {
	const prefix = "BizQueryRoomLockStatus "

	status = entity.UNLOCKROOM
	resp, err := bizRoomCenterClient.GetBizRoomLockStatus(context.TODO(), &bilin.GetBizRoomLockStatusReq{Roomid: room.Roomid})
	if err != nil || resp.Commonret.Ret != bccommon.SUCCESSMESSAGE.Ret {
		log.Error(prefix+"failed", zap.Any("roomid", room.Roomid), zap.Any("err", err), zap.Any("resp", resp))
		return
	}

	status = resp.Lockstatus

	log.Info(prefix, zap.Any("roomid", room.Roomid), zap.Any("status", status))
	return
}

func BizSetRoomPassword(room *entity.Room, roompwd string) (err error) {
	const prefix = "BizSetRoomPassword "

	//调用tars接口存密码
	resp, err := bizRoomCenterClient.SetRoomPassword(context.TODO(), &bilin.SetRoomPasswordReq{Roomid: room.Roomid, Password: roompwd})
	if err != nil || resp.Commonret.Ret != bccommon.SUCCESSMESSAGE.Ret {
		log.Error(prefix+"failed", zap.Any("roomid", room.Roomid))
		return
	}

	room.LockStatus = entity.LOCKROOM
	//更新redis
	service.RedisAddRoom(room)

	//notify kafka
	jsonBytes, err := json.Marshal(&KafkaMessage{Roomid: room.Roomid, LockStatus: room.LockStatus, Owner: room.Owner})
	if err != nil {
		log.Error(prefix+"json.Marshal(KafkaMessage)", zap.Any("err", err))
		return
	}
	service.KafkaProduceMessage(string(jsonBytes))

	log.Info(prefix, zap.Any("roomid", room.Roomid))
	return
}

func BizRemoveRoomPassword(room *entity.Room) (err error) {
	const prefix = "BizRemoveRoomPassword "

	//调用tars接口删除密码
	resp, err := bizRoomCenterClient.RemoveRoomPassword(context.TODO(), &bilin.RemoveRoomPasswordReq{Roomid: room.Roomid})
	if err != nil || resp.Commonret.Ret != bccommon.SUCCESSMESSAGE.Ret {
		log.Error(prefix+"failed", zap.Any("roomid", room.Roomid))
		return
	}

	room.LockStatus = entity.UNLOCKROOM
	//更新redis
	service.RedisAddRoom(room)

	//notify kafka
	jsonBytes, err := json.Marshal(&KafkaMessage{Roomid: room.Roomid, LockStatus: room.LockStatus, Owner: room.Owner})
	if err != nil {
		log.Error(prefix+"json.Marshal(KafkaMessage)", zap.Any("err", err))
		return
	}
	service.KafkaProduceMessage(string(jsonBytes))

	log.Info(prefix, zap.Any("roomid", room.Roomid))
	return
}
