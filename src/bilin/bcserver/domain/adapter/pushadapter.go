package adapter

import (
	"bilin/bcserver/domain/entity"
	"bilin/protocol"
	"github.com/golang/protobuf/proto"

	"bilin/bcserver/domain/collector"
	"bilin/common/pushproxy"
	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
)

type PushService struct {
}

func pushBroadcastToRoom(room *entity.Room, msg proto.Message, msgType bilin.MinType_BC) (err error) {
	const prefix = "pushToRoom"
	var body bilin.BcMessageBody
	body.Type = int32(msgType)
	if body.Data, err = proto.Marshal(msg); err != nil {
		log.Error(prefix+"[-]proto.Marshal failed", zap.Any("err", err))
	}
	pushData, err := proto.Marshal(&body)
	if err != nil {
		log.Error(prefix+"[-]proto.Marshal failed", zap.Any("err", err))
		return
	}

	pushMsg, err := proto.Marshal(&bilin.ServerPush{
		MessageType: int32(bilin.MaxType_BC_MSG),
		PushBuffer:  pushData,
		MessageDesc: "给房间推送广播消息",
		ServiceName: "bc_server",
		MethodName:  "pushBroadcastToRoom",
	})
	if err != nil {
		log.Error(prefix+"[-]proto.Marshal failed", zap.Any("err", err))
		return
	}

	err = pushproxy.PushToRoom(int32(room.Roomid), int32(bilin.MaxType_BC_MSG), string(pushMsg))
	if err != nil {
		log.Error("[-]PushBaseRoomInfoToUser failed push", zap.Any("err", err))
	}

	log.Debug("[+]pushBroadcastToRoom success push", zap.Any("minType", bilin.MinType_BC_name[int32(msgType)]), zap.Any("PushBuffer length", len(pushMsg)))
	return
}

//push通知给用户,包括踢人、禁言、禁麦等，传递的消息类型是mintype
func PushNotifyToUser(roomid uint64, uids []int64, msg proto.Message, msgType bilin.MinType_BC) (err error) {
	const prefix = "PushNotifyToUser "

	var body bilin.BcMessageBody
	body.Type = int32(msgType)
	if body.Data, err = proto.Marshal(msg); err != nil {
		log.Error(prefix+"[-]proto.Marshal failed", zap.Any("err", err))
	}
	pushBody, err := proto.Marshal(&body)
	if err != nil {
		log.Error(prefix+"[-]proto.Marshal failed", zap.Any("err", err))
		return
	}

	multiMsg, err := proto.Marshal(&bilin.MultiPush{
		Msg: &bilin.ServerPush{
			MessageType: int32(bilin.MaxType_BC_MSG),
			PushBuffer:  pushBody,
			MessageDesc: "给房间推送单播消息",
			ServiceName: "bc_server",
			MethodName:  "PushNotifyToUser",
		},
		UserIDs: uids,
	})
	if err != nil {
		log.Error(prefix+"[-]proto.Marshal failed", zap.Any("err", err))
		return
	}

	//appid = tunnel.AppidType_NEW_BCSERVER的时候uid直接填0，uids从protobuf里面取
	pushproxy.PushToUser(int32(roomid), 0, int32(bilin.MaxType_BC_MSG), string(multiMsg))
	if err != nil {
		log.Error("[-]PushNotifyToUser failed push", zap.Any("err", err))
	}

	log.Debug("[+]PushNotifyToUser success push", zap.Any("uids", uids), zap.Any("minType", bilin.MinType_BC_name[int32(msgType)]), zap.Any("PushBuffer length", len(string(multiMsg))))
	return
}

func PushBaseRoomInfoToUser(room *entity.Room, uids []int64) (err error) {
	err = PushNotifyToUser(room.Roomid, uids, collector.BaseRoomInfo(room), bilin.MinType_BC_NotifyBaseRoomInfo)
	return
}

func PushAllRoomInfoToUser(room *entity.Room, uids []int64) (err error) {
	err = PushNotifyToUser(room.Roomid, uids, collector.AllRoomInfo(room), bilin.MinType_BC_NotifyAllRoomInfo)
	if err != nil {
		log.Error("[-]PushAllRoomInfoToUser failed push", zap.Any("err", err))
	}

	log.Debug("[+]PushAllRoomInfoToUser success push", zap.Any("roomid", room.Roomid))
	return
}

func PushBaseRoomInfoToRoom(room *entity.Room) (err error) {
	err = pushBroadcastToRoom(room, collector.BaseRoomInfo(room), bilin.MinType_BC_NotifyBaseRoomInfo)
	if err != nil {
		log.Error("[-]PushBaseRoomInfoToRoom failed push", zap.Any("err", err))
	}

	log.Debug("[+]PushBaseRoomInfoToRoom success push", zap.Any("room", room))
	return
}

func PushAllRoomInfoToRoom(room *entity.Room) (err error) {
	err = pushBroadcastToRoom(room, collector.AllRoomInfo(room), bilin.MinType_BC_NotifyAllRoomInfo)
	if err != nil {
		log.Error("[-]PushAllRoomInfoToRoom failed push", zap.Any("err", err))
	}

	log.Debug("[+]PushAllRoomInfoToRoom success push", zap.Any("roomid", room.Roomid))
	return
}

func PushUserPraiseInfoToRoom(room *entity.Room, praiseCount uint32) (err error) {
	err = pushBroadcastToRoom(room, &bilin.PraiseNotify{Count: praiseCount}, bilin.MinType_BC_NotifyRoomPraise)
	if err != nil {
		log.Error("[-]PushUserPraiseInfoToRoom failed push", zap.Any("err", err))
	}

	log.Debug("[+]PushUserPraiseInfoToRoom success push", zap.Any("roomid", room.Roomid))
	return
}

func PushUserListChangeToRoom(room *entity.Room, enterUsers []*entity.User, exitUids []uint64) (err error) {
	var pushEnterUsers []*bilin.UserInfo
	for _, item := range enterUsers {
		pushUser := collector.LocalUserToSendInfo(item)
		pushEnterUsers = append(pushEnterUsers, pushUser)
	}
	err = pushBroadcastToRoom(room, &bilin.UserListChangeNotify{Enterusers: pushEnterUsers, Exituids: exitUids}, bilin.MinType_BC_NotifyRoomUserListChange)
	if err != nil {
		log.Error("[-]PushUserListChangeToRoom failed push", zap.Any("err", err))
	}

	log.Debug("[+]PushUserListChangeToRoom success push", zap.Any("roomid", room.Roomid))
	return
}

func PushUserListInfoToRoom(room *entity.Room) (err error) {
	err = pushBroadcastToRoom(room, collector.UserList(room), bilin.MinType_BC_NotifyRoomUserListInfo)
	if err != nil {
		log.Error("[-]PushUserListInfoToRoom failed push", zap.Any("err", err))
	}

	log.Debug("[+]PushUserListInfoToRoom success push", zap.Any("roomid", room.Roomid))
	return
}

func PushMikeListInfoToRoom(room *entity.Room) (err error) {
	err = pushBroadcastToRoom(room, collector.MikeList(room), bilin.MinType_BC_NotifyRoomMickListInfo)
	if err != nil {
		log.Error("[-]PushMikeListInfoToRoom failed push", zap.Any("err", err))
	}

	log.Debug("[+]PushMikeListInfoToRoom success push", zap.Any("roomid", room.Roomid))
	return
}

//被禁止公屏发言的用户
func PushBlackListInfoToRoom(room *entity.Room) (err error) {
	err = pushBroadcastToRoom(room, collector.ForbiddenTextUserList(room), bilin.MinType_BC_NotifyRoomForbiddenList)
	if err != nil {
		log.Error("[-]PushBlackListInfoToRoom failed push", zap.Any("err", err))
	}

	log.Debug("[+]PushBlackListInfoToRoom success push", zap.Any("roomid", room.Roomid))
	return
}

//公屏发言
func PushBroIMMsgToRoom(room *entity.Room, data []byte) (err error) {
	const prefix = "PushBroIMMsgToRoom "
	var body bilin.CommonMessageBody
	body.Type = int32(bilin.MinType_COMMON_IM_MSG)
	body.Data = data
	pushData, err := proto.Marshal(&body)
	if err != nil {
		log.Error(prefix+"[-]proto.Marshal failed", zap.Any("err", err))
		return
	}

	pushMsg, err := proto.Marshal(&bilin.ServerPush{
		MessageType: int32(bilin.MaxType_COMMON_MSG),
		PushBuffer:  pushData,
		MessageDesc: "公屏发言",
		ServiceName: "bc_server",
		MethodName:  "PushBroIMMsgToRoom",
	})
	if err != nil {
		log.Error(prefix+"[-]proto.Marshal failed", zap.Any("err", err))
		return
	}

	err = pushproxy.PushToRoom(int32(room.Roomid), int32(bilin.MaxType_COMMON_MSG), string(pushMsg))
	if err != nil {
		log.Error("[-]PushBroIMMsgToRoom failed push", zap.Any("err", err))
	}

	log.Debug("[+]PushBroIMMsgToRoom success push", zap.Any("minType", bilin.MinType_BC_name[int32(bilin.MinType_COMMON_IM_MSG)]), zap.Any("PushBuffer length", len(pushMsg)))
	return
}

func PushRoomClosedNotifyToRoom(room *entity.Room, hostnotifytext string, audiencenotifytext string) (err error) {
	err = pushBroadcastToRoom(room, &bilin.ClosedRoomNotify{Hostnotifytext: hostnotifytext, Audiencenotifytext: audiencenotifytext}, bilin.MinType_BC_NotifyRoomClosed)
	if err != nil {
		log.Error("[-]PushRoomClosedNotifyToRoom failed push", zap.Any("err", err))
	}

	log.Debug("[+]PushRoomClosedNotifyToRoom success push", zap.Any("room", room))
	return
}

//K歌相关push，都是push到房间的广播消息

//即将演唱的歌曲预告 不做此推送，歌曲列表更新时带上预告歌曲
//func PushPreparationSongToRoom(room *entity.Room) (err error) {
//	err = pushBroadcastToRoom(room, collector.GetPrepareSong(room), bilin.MinType_BC_NotifyPreparationSong)
//	if err != nil {
//		log.Error("[-]PushPreparationSongToRoom failed push", zap.Any("err", err))
//	}
//
//	log.Debug("[+]PushPreparationSongToRoom success push", zap.Any("room", room))
//	return
//}

//当前点歌列表
func PushSongsListToRoom(room *entity.Room) (err error) {
	err = pushBroadcastToRoom(room, &bilin.SongsListNotify{Songs: collector.GetSongsList(room)}, bilin.MinType_BC_NotifySongsList)
	if err != nil {
		log.Error("[-]PushSongsListToRoom failed push", zap.Any("err", err))
	}

	log.Debug("[+]PushSongsListToRoom success push", zap.Any("room", room))
	return
}

//嘉宾添加歌曲
func PushAddSongToRoom(room *entity.Room, song *bilin.KaraokeSongInfo) (err error) {
	err = pushBroadcastToRoom(room, &bilin.AddSongNotify{Song: song}, bilin.MinType_BC_NotifyAddSong)
	if err != nil {
		log.Error("[-]PushAddSongToRoom failed push", zap.Any("err", err))
	}

	log.Debug("[+]PushAddSongToRoom success push", zap.Any("room", room))
	return
}

//主持人开始播放歌曲
func PushStartSingToRoom(room *entity.Room, song *bilin.KaraokeSongInfo) (err error) {
	err = pushBroadcastToRoom(room, &bilin.StartSingNotify{Song: song}, bilin.MinType_BC_NotifyStartSing)
	if err != nil {
		log.Error("[-]PushStartSingToRoom failed push", zap.Any("err", err))
	}

	log.Debug("[+]PushStartSingToRoom success push", zap.Any("room", room), zap.Any("song", song))
	return
}

//暂停歌曲
func PushPauseSongToRoom(room *entity.Room, song *bilin.KaraokeSongInfo) (err error) {
	err = pushBroadcastToRoom(room, &bilin.PauseSongNotify{Song: song}, bilin.MinType_BC_NotifyPauseSong)
	if err != nil {
		log.Error("[-]PushPauseSongToRoom failed push", zap.Any("err", err))
	}

	log.Debug("[+]PushPauseSongToRoom success push", zap.Any("room", room), zap.Any("song", song))
	return
}

//结束歌曲
func PushTerminateSongToRoom(room *entity.Room, optUid uint64, song *bilin.KaraokeSongInfo) (err error) {
	err = pushBroadcastToRoom(room, &bilin.TerminateSongNotify{Song: song, Optuserid: optUid}, bilin.MinType_BC_NotifyTerminateSong)
	if err != nil {
		log.Error("[-]PushTerminateSongToRoom failed push", zap.Any("err", err))
	}

	log.Debug("[+]PushTerminateSongToRoom success push", zap.Any("room", room), zap.Any("song", song))
	return
}
