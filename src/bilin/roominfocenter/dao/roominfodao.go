package dao

import (
	"bilin/bcserver/domain/entity"
	"bilin/protocol"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"math"
	"strconv"
	"time"
)

// from bcserver domain entity
// 做了一些裁剪，注意对应的域和bcserver保持一致
type Room struct {
	Roomid     uint64                        `json:"roomid"`
	Owner      uint64                        `json:"owner"`
	Status     bilin.BaseRoomInfo_ROOMSTATUS `json:"status"`
	RoomType   bilin.BaseRoomInfo_ROOMTYPE   `json:"roomtype"`
	Linkstatus bilin.BaseRoomInfo_LINKSTATUS `json:"linkstatus"`
	Title      string                        `json:"title"`
	//TotalMickNumber  uint32
	//MikeWaitingUsers uint32
	RoomType2        int32  `json:"roomType2"`
	RoomCategoryID   int32  `json:"roomCategoryID"`
	RoomPendantLevel int32  `json:"roomPendantLevel"`
	HostBilinID      int64  `json:"hostBilinID"`
	StartTime        uint64 `json:"starttime"`
	LockStatus       int32  `json:"lock_status"` // 0 未锁，1 锁定
}

// SyncLivingRoomInfosByScan return only bilin.BaseRoomInfo_OPEN room
func SyncLivingRoomInfosByScan() (roomList map[uint64]*entity.Room, err error) {
	if RedisClient == nil {
		return nil, redisNotInitErr
	}
	const prefix = "SyncLivingRoomInfos "
	var (
		roomid uint64 = math.MaxUint64
		keys   []string
		cursor uint64 = 0
	)
	ctx, cancle := context.WithDeadline(context.Background(), time.Now().Add(2*time.Second))
	defer cancle()
	roomList = make(map[uint64]*entity.Room)
	for {
		keys, cursor, err = RedisClient.HScan("bc_roomlist", cursor, "", 1000).Result()
		if err != nil && err != redis.Nil {
			appzaplog.Error(prefix+"redis.HGetAll", zap.Error(err))
			return
		}

		for index, item := range keys {
			if index%2 != 0 {
				if roomid != math.MaxUint64 {
					room := &entity.Room{}
					if marshalErr := json.Unmarshal([]byte(item), room); marshalErr != nil {
						appzaplog.Warn(prefix+"Unmarshal failed", zap.Any("room", room), zap.Any("item", item))
						continue
					}
					// check room status here
					if room.Status == bilin.BaseRoomInfo_OPEN {
						roomList[roomid] = room
					}
				}
			} else {
				roomid, err = strconv.ParseUint(item, 10, 64)
				if err != nil {
					roomid = math.MaxUint64 // illegal roomid
					continue
				}
			}

		}
		if cursor == 0 {
			break
		}
		if deadline, ok := ctx.Deadline(); ok {
			if time.Now().After(deadline) {
				appzaplog.Error("hscan take too much time", zap.Int("alreadyscanlen", len(roomList)))
				err = errors.New("hscan take more than 2 second")
				break
			}
		}
	}

	return
}

func SyncLivingRoomInfos() (roomList map[uint64]*entity.Room, err error) {
	if RedisClient == nil {
		return nil, redisNotInitErr
	}
	const prefix = "SyncLivingRoomInfos "
	var (
		redisVal map[string]string
		roomid   uint64
		//room     = &entity.Room{}
	)
	roomList = make(map[uint64]*entity.Room)
	redisVal, err = RedisClient.HGetAll("bc_roomlist").Result()
	if err != nil && err != redis.Nil {
		appzaplog.Error(prefix+"redis.HGetAll", zap.Error(err))
		return
	}

	for roomidStr, value := range redisVal {
		room := &entity.Room{}
		if roomid, err = strconv.ParseUint(roomidStr, 10, 64); err != nil {
			appzaplog.Error("RedisGetRoomIdList ParseUint failed", zap.Error(err))
			continue
		}
		if len(value) == 0 {
			appzaplog.Info(prefix+"room not find in redis", zap.Uint64("roomid", roomid))
			continue
		}
		if err = json.Unmarshal([]byte(value), room); err != nil {
			appzaplog.Warn(prefix+"Unmarshal failed", zap.Uint64("roomid", roomid))
			continue
		}
		roomList[roomid] = room
	}
	//appzaplog.Debug(prefix, zap.Any("roomList", roomList))
	return
}

func SyncUserCount(roomid uint64) (total int64, err error) {
	if RedisClient == nil {
		return 0, redisNotInitErr
	}
	const prefix = "SyncUserCount "

	if total, err = RedisClient.HLen(fmt.Sprintf("users_%d", roomid)).Result(); err != nil {
		return 0, fmt.Errorf("查询redis失败roomid: %d", roomid)
	}

	//appzaplog.Debug(prefix, zap.Uint64("roomid", roomid), zap.Int64("total", total))
	return total, nil
}

func HostInRoom(hostid, roomid uint64) (bool, error) {
	if RedisClient == nil {
		// fake hostonline
		return true, nil
	}
	if exist, err := RedisClient.HExists(fmt.Sprintf("stage_%d", roomid), fmt.Sprintf("%d", hostid)).Result(); err != nil {
		return true, fmt.Errorf("查询redis失败roomid: %d", roomid)
	} else {
		return exist, nil
	}
}
