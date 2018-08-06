package cache

import (
	"bilin/protocol"
	"bilin/roominfocenter/dao"
	"sync"

	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
)

var (
	lock          sync.Mutex
	roominfocache map[uint64]*bilin.RoomInfo
)

func RefreshRoomCache() error {
	//rooms, err := dao.SyncLivingRoomInfos()
	rooms, err := dao.SyncLivingRoomInfosByScan()
	if err != nil {
		appzaplog.Error("[+]LivingRoomsInfo SyncLivingRoomInfos failed", zap.Error(err))
		return err
	}

	Livingrooms := make(map[uint64]*bilin.RoomInfo)
	// take usernumbers
	for roomid, info := range rooms {
		if info.Status != bilin.BaseRoomInfo_OPEN {
			continue
		}
		info := &bilin.RoomInfo{
			Roomid:         roomid,
			Starttime:      info.StartTime,
			RoomcategoryID: info.RoomCategoryID,
			Owner:          info.Owner,
			Title:          info.Title,
			RoomType2:      info.RoomType2,
			OwnerBilinID:   uint64(info.HostBilinID),
			LockStatus:     info.LockStatus,
		}
		if exist, _ := dao.HostInRoom(info.Owner, info.Roomid); !exist {
			appzaplog.Debug("[+]LivingRoomsInfo host absent room", zap.Uint64("owner", info.Owner), zap.Uint64("roomid", info.Roomid))
			continue
		}
		if number, err := dao.SyncUserCount(roomid); err != nil {
			appzaplog.Error("[-]LivingRoomsInfo SyncUserCount failed", zap.Uint64("roomid", roomid), zap.Uint64("owner", info.Owner), zap.Error(err))
			return err
		} else {
			appzaplog.Debug("[+]LivingRoomsInfo SyncUserCount ok", zap.Uint64("roomid", roomid), zap.Uint64("owner", info.Owner), zap.Int64("count", number))
			info.Usernumber = uint64(number)
		}
		Livingrooms[roomid] = info
	}
	appzaplog.Debug("refresh", zap.Int("roomnum", len(Livingrooms)))
	lock.Lock()
	roominfocache = Livingrooms
	lock.Unlock()
	return nil
}

func GetRoomCache() (info map[uint64]*bilin.RoomInfo) {
	lock.Lock()
	info = roominfocache
	lock.Unlock()
	return
}
