package collector

import (
	"bilin/bcserver/domain/service"
	"bilin/protocol"
	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
)

func GetRoomUsersByPage(roomid uint64, pageNum uint32) (result []*bilin.UserInfo) {
	const prefix = "GetRoomUsersByPage"
	userlist, _ := service.RedisGetDisplayedUsers(roomid, 0)
	if userlist == nil {
		log.Debug(prefix+"end", zap.Any("roomid", roomid), zap.Any("userlist", userlist))
		return nil
	}

	begin_pos := PageUsersCount * int32(pageNum-1)
	end_pos := PageUsersCount * int32(pageNum)
	for idx, item := range userlist {
		if idx >= int(begin_pos) && idx < int(end_pos) { //需要过滤掉主持人 todo
			result = append(result, LocalUserToSendInfo(item))
		}
	}

	log.Debug(prefix+"end", zap.Any("roomid", roomid), zap.Any("result", result))
	return result
}
