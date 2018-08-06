package collector

import (
	"bilin/bcserver/domain/entity"
	"bilin/bcserver/domain/service"
	"bilin/protocol"

	"bilin/bcserver/bccommon"
	"bilin/thrift/gen-go/hotline"
	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/kissgo/httpmetrics"
	"sort"
	"time"
)

const (
	DisplayedNum    = 40
	ApplyMikeLimits = 50
	PageUsersCount  = 40
)

//由于老版本有机器人在房间的情况，不能简单根据房间人数为0就删除房间。 先保留接口
func DelRoom(room *entity.Room) (err error) {
	return
}

func AllRoomInfo(room *entity.Room) *bilin.AllRoomInfo {
	return &bilin.AllRoomInfo{
		Baseinfo:      BaseRoomInfo(room),
		Audienceusers: getAudienceUsers(room),
		Mikeinfo:      getMikes(room),
		Forbiddenuids: ForbiddenTextUserList(room),
		Karaokeinfo:   AllRoomKaraokeInfo(room),
		Bizinfo:       getbizinfo(room),
	}
}

func BaseRoomInfo(room *entity.Room) *bilin.BaseRoomInfo {
	return &bilin.BaseRoomInfo{
		Roomid:             room.Roomid,
		Roomstatus:         room.Status,
		Roomtype:           room.RoomType,
		Linkstatus:         room.GetLinkStatus(),
		Title:              room.Title,
		RoomType2:          room.RoomType2,
		RoomCategoryID:     room.RoomCategoryID,
		RoomPendantLevel:   room.RoomPendantLevel,
		HostBilinID:        room.HostBilinID,
		Totalmicknumber:    getMickNumberByRoomType(room.RoomType),
		Host:               getOwner(room),
		Mikewaitingusers:   getApplyMikesNum(room),
		Totalusersnumber:   getTotalUsersNumber(room),
		PageUsersCount:     PageUsersCount,
		Autolink:           room.GetAutoLink(),
		Maixuswitch:        room.Maixuswitch,
		Karaokeswitch:      room.Karaokeswitch,
		Relationlistswitch: room.Relationlistswitch,
		Owneruid:           room.Owner,
	}
}

func MikeList(room *entity.Room) (mickList *bilin.RoomMickListInfo) {
	mickList = new(bilin.RoomMickListInfo)
	mickList.Mikeinfo = getMikes(room)
	mickList.Mikewaitingusers = getApplyMikesNum(room)
	mickList.Roomtype = room.RoomType

	log.Debug("MikeList ", zap.Any("roomid", room.Roomid), zap.Any("mickList", mickList))
	return mickList
}
func UserList(room *entity.Room) (userList *bilin.RoomUserListInfo) {
	userList = new(bilin.RoomUserListInfo)
	userList.Audienceusers = getAudienceUsers(room)
	userList.Totalusersnumber = getTotalUsersNumber(room)

	log.Debug("getAudienceUsers ", zap.Any("roomid", room.Roomid), zap.Any("userList", userList))
	return userList
}

func ForbiddenTextUserList(room *entity.Room) (forbiddenList *bilin.RoomForbiddenList) {
	forbiddenList = new(bilin.RoomForbiddenList)
	result, _ := service.RedisGetForbidenUserList(room.Roomid)
	forbiddenList.Uids = result

	log.Debug("ForbiddenTextUserList ", zap.Any("roomid", room.Roomid), zap.Any("forbiddenList", forbiddenList))
	return
}

func LocalUserToSendInfo(entity_user *entity.User) *bilin.UserInfo {
	return &bilin.UserInfo{
		Userid:    entity_user.UserID,
		Nick:      entity_user.NickName,
		Avatarurl: entity_user.AvatarURL,
		Fanscount: entity_user.FansCount,
		From:      bilin.USERFROM_BROADCAST,
		Mute:      entity_user.IsMuted,
		Sex:       uint32(entity_user.Sex),
		Age:       uint32(entity_user.Age),
		CityName:  entity_user.CityName,
		Signature: entity_user.Signature,
	}
}

func InitMikeWheatInfo(room *entity.Room) (err error) {
	//0号麦是给主持人用的，所以初始化的时候从1开始
	mikeNume := getMickNumberByRoomType(room.RoomType)
	var idx uint32
	for idx = 1; idx <= mikeNume; idx++ {
		service.RedisLockUnlockMikeWheat(room.Roomid, idx, bilin.MikeInfo_EMPTY)
	}

	return
}

//切换模板时，复用麦位信息
func ReuseMikeWheatInfo(room *entity.Room, newType bilin.BaseRoomInfo_ROOMTYPE) (err error) {
	allStatus, _ := service.RedisGetAllMikeWheatStatus(room.Roomid)
	newMikeNum := getMickNumberByRoomType(newType)
	for key, _ := range allStatus {
		if uint32(key) > newMikeNum {
			service.RedisRemoveMikeWheat(room.Roomid, uint32(key))
		}
	}

	//主持人占用0号麦，需要+1
	for idx := len(allStatus); idx < int(newMikeNum)+1; idx++ {
		service.RedisLockUnlockMikeWheat(room.Roomid, uint32(idx), bilin.MikeInfo_EMPTY)
	}

	log.Debug("ReuseMikeWheatInfo ", zap.Any("room", room))
	return
}

func getOwner(room *entity.Room) *bilin.UserInfo {
	if user, _ := service.RedisGetUser(room.Roomid, room.Owner); user != nil {
		return LocalUserToSendInfo(user)
	}

	return nil
}

func getMickNumberByRoomType(roomType bilin.BaseRoomInfo_ROOMTYPE) uint32 {
	switch roomType {
	case bilin.BaseRoomInfo_ROOMTYPE_THREE:
		return 3
	case bilin.BaseRoomInfo_ROOMTYPE_SIX:
		return 6
	case bilin.BaseRoomInfo_ROOMTYPE_RADIO: //一个主持人一个嘉宾
		return 1
	default:
		return 0
	}
}

func IsValidRoomType(roomType bilin.BaseRoomInfo_ROOMTYPE) bool {
	if roomType != bilin.BaseRoomInfo_ROOMTYPE_RADIO && roomType != bilin.BaseRoomInfo_ROOMTYPE_THREE && roomType != bilin.BaseRoomInfo_ROOMTYPE_SIX {
		return false
	}

	return true
}

func CheckMikeNumberUsable(room *entity.Room, mikeIndex uint32) (ret bool) {
	if mikeIndex == 0 || mikeIndex > getMickNumberByRoomType(room.RoomType) {
		log.Warn("CheckMickNumber ", zap.Any("roomid", room.Roomid), zap.Any("mikeIndex", mikeIndex), zap.Any("ret", false))
		return false
	}

	//麦位状态为空才算可用的
	if status, _ := service.RedisGetMikeWheatStatus(room.Roomid, mikeIndex); status != int(bilin.MikeInfo_EMPTY) {
		log.Warn("CheckMickNumber ", zap.Any("roomid", room.Roomid), zap.Any("mikeIndex", mikeIndex), zap.Any("status", status), zap.Any("ret", false))
		return false
	}

	log.Debug("CheckMickNumber ", zap.Any("roomid", room.Roomid), zap.Any("mikeIndex", mikeIndex), zap.Any("ret", true))
	return true
}

func getAudienceUsers(room *entity.Room) (result []*bilin.UserInfo) {
	if userList, _ := service.RedisGetDisplayedUsers(room.Roomid, DisplayedNum); userList != nil {
		for _, item := range userList {
			result = append(result, LocalUserToSendInfo(item))
		}
	}

	log.Debug("CheckMikeNumberUsable ", zap.Any("roomid", room.Roomid), zap.Any("result", result), zap.Any("size", len(result)))

	return
}

//用来检查是否有房间的麦位信息有问题
func CheckMikeErrors(room *entity.Room) {
	const prefix = "CheckMikeErrors "

	mikestatus, _ := service.RedisGetAllMikeWheatStatus(room.Roomid)
	userlist, _ := service.RedisGetOnMikeUserList(room.Roomid)

	mikeNumbers := getMickNumberByRoomType(room.RoomType)
	mikeListLength := len(userlist)

	//先检查房间的麦位数和当前麦位数是否一致  主持人占一个麦位
	if mikeNumbers+1 < uint32(len(mikestatus)) {
		log.Error(prefix+"mikeNumbers  less than len(mikestatus)", zap.Any("roomid", room.Roomid), zap.Any("mikeNumbers", mikeNumbers), zap.Any("mikestatus", mikestatus))
		return
	}

	//检查麦上用户和麦位个数是否相符合
	if mikeNumbers+1 < uint32(mikeListLength) {
		log.Error(prefix+"mikeNumbers less than mikeListLength", zap.Any("roomid", room.Roomid), zap.Any("mikeNumbers", mikeNumbers), zap.Any("mikeListLength", mikeListLength))
		return
	}

	//检查麦上有人，但是麦位为空
	for _, item := range userlist {
		if mikestatus[int(item.MikeIndex)] == int(bilin.MikeInfo_EMPTY) {
			log.Error(prefix+"user on mike, but status == MikeInfo_EMPTY ", zap.Any("user", item))
		}
	}

	//检查同一个麦位是否有多人的情况
	mikeUsers := make(map[int]int)
	for mikeIndex, mikeStatu := range mikestatus {
		for _, item := range userlist {
			if item.MikeIndex == uint32(mikeIndex) {
				mikeUsers[mikeIndex] += 1
			}
		}

		// 麦上没人，麦位为已使用的 情况
		if mikeStatu == int(bilin.MikeInfo_USED) && mikeUsers[mikeIndex] == 0 {
			log.Error(prefix+"user not on mike, but status == MikeInfo_USED ", zap.Any("roomid", room.Roomid), zap.Any("mikeIndex", mikeIndex))
		}

		if mikeUsers[mikeIndex] > 1 {
			log.Error(prefix+"more than 2 people on mike ", zap.Any("roomid", room.Roomid), zap.Any("mikeIndex", mikeIndex), zap.Any("userCount", mikeUsers[mikeIndex]))
		}
	}

}

//在麦上，不在麦上的信息都要返回
func getMikes(room *entity.Room) (result []*bilin.MikeInfo) {
	allStatus, _ := service.RedisGetAllMikeWheatStatus(room.Roomid)
	userList, _ := service.RedisGetOnMikeUserList(room.Roomid)

	if room.Maixuswitch == bilin.BaseRoomInfo_CLOSEMAIXU {
		//没有麦序概念的时候，需要根据上麦时间来排麦序
		sort.Stable(entity.UserSortByOnMikeTimeSlice(userList))
		mikeNum := 1
		for _, item := range userList {
			mike := new(bilin.MikeInfo)
			mike.Mikewheatstatus = bilin.MikeInfo_USED //老版本只要有人就是used
			if item.MikeIndex == 0 {
				mike.Mikeindex = 0 // 主播的麦序永远为0
				//填充粉丝数量
				item.FansCount, _ = service.RedisGetUserFansCount(item.UserID)
			} else {
				mike.Mikeindex = uint32(mikeNum)
				mikeNum += 1 // 其他人的麦序跟着往后面叠加
			}

			mike.Userinfo = LocalUserToSendInfo(item)
			result = append(result, mike)
		}

		//填充剩余的麦位
		for remainMike := mikeNum; remainMike <= int(getMickNumberByRoomType(room.RoomType)); remainMike++ {
			mike := new(bilin.MikeInfo)
			mike.Mikeindex = uint32(remainMike)
			mike.Mikewheatstatus = bilin.MikeInfo_EMPTY
			result = append(result, mike)
		}
	} else {
		//有麦序，按照麦序顺序返回
		for key, value := range allStatus {
			mike := new(bilin.MikeInfo)
			mike.Mikeindex = uint32(key)
			mike.Mikewheatstatus = bilin.MikeInfo_MIKEWHEATSTATUS(value)
			for _, item := range userList {
				if item.MikeIndex == uint32(key) {
					if item.MikeIndex == 0 {
						//填充粉丝数量
						item.FansCount, _ = service.RedisGetUserFansCount(item.UserID)
					}
					mike.Userinfo = LocalUserToSendInfo(item)
				}
			}

			result = append(result, mike)
		}
	}

	log.Debug("getMikes ", zap.Any("roomid", room.Roomid), zap.Any("result", result))
	return
}

func getApplyMikesNum(room *entity.Room) uint32 {
	ret, _ := service.RedisGetApplyMikeUserCount(room.Roomid)

	log.Debug("getApplyMikesNum ", zap.Any("roomid", room.Roomid), zap.Any("ret", ret))
	return uint32(ret)
}

func getTotalUsersNumber(room *entity.Room) uint32 {
	ret, _ := service.RedisGetUserCount(room.Roomid)

	log.Debug("getTotalUsersNumber ", zap.Any("roomid", room.Roomid), zap.Any("ret", ret))
	return uint32(ret)
}

func InitRoomByJavaResult(RoomId uint64, uid uint64, joinHotLineRet *hotline.JoinHotLineRet) (room *entity.Room) {
	const prefix = "InitRoomByJavaResult "

	// 用户进房间, 查看房间是否存在
	if room = GetRoomInfoByRoomId(RoomId); room == nil {
		log.Info(prefix+"room not exist,create it ", zap.Any("RoomId", RoomId))

		//redis  中只需要存房间的基础信息
		room = entity.NewRoom(RoomId)

		//set roominfo
		room.Title = joinHotLineRet.GetTitle()
		room.RoomType2 = joinHotLineRet.GetRoomType()
		room.RoomCategoryID = joinHotLineRet.GetRoomCategoryId()
		room.HostBilinID = joinHotLineRet.GetHostBilinId()
		room.RoomPendantLevel = joinHotLineRet.GetRoomPendantLevel()
		room.LockStatus = BizQueryRoomLockStatus(room)
		room.RoomType = bilin.BaseRoomInfo_ROOMTYPE_THREE
		StorageRoomInfo(room)

		//统计用
		CreateRoomStat(room.Roomid, uid)
	}

	//亲密榜开关
	room.Relationlistswitch = bilin.BaseRoomInfo_OPENRELATIONLIST

	if room.Status == bilin.BaseRoomInfo_CLOSED {
		room.Status = bilin.BaseRoomInfo_OPEN
		StorageRoomInfo(room)
	}

	log.Debug(prefix, zap.Any("room", room))
	return
}

func StorageRoomInfo(room *entity.Room) (err error) {
	const prefix = "StorageRoomInfo "

	//更新redis
	service.RedisAddRoom(room)

	//更新mysql
	service.MysqlStorageRoomInfo(room)

	log.Debug(prefix, zap.Any("room", room))
	return
}

func GetRoomInfoByRoomId(RoomId uint64) (room *entity.Room) {
	const prefix = "GetRoomInfoByRoomId "

	defer func(now time.Time) {
		httpmetrics.DefReport("GetRoomInfoByRoomId", 0, now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	//先从redis取
	if room, _ = service.RedisGetRoomInfo(RoomId); room == nil {
		//再从mysql中获取
		room, _ = service.MysqlGetRoomInfo(RoomId)
		if room != nil {
			service.RedisAddRoom(room)
		}
	}

	log.Debug(prefix, zap.Any("room", room))
	return
}

func getbizinfo(room *entity.Room) (result *bilin.RoomBizInfo) {
	const prefix = "getbizinfo "

	result = &bilin.RoomBizInfo{LockStatus: room.LockStatus}

	log.Info(prefix, zap.Any("room", room))
	return
}
