package service

import (
	"bilin/bcserver/domain/entity"
	"encoding/json"
	"github.com/go-redis/redis"
	"time"

	"bilin/bcserver/bccommon"
	"bilin/bcserver/config"
	"bilin/protocol"
	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	httpmetrics "code.yy.com/yytars/goframework/kissgo/httpmetrics"
	"fmt"
	"sort"
	"strconv"
)

var (
	RedisClient *redis.Client
)

func RedisInit() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:         config.GetAppConfig().RedisAddr,
		Password:     "", // no password set
		DB:           0,  // use default DB
		DialTimeout:  1 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		PoolSize:     10,
		PoolTimeout:  5 * time.Second,
	})

}

//禁用，这个api请勿调用，仅仅用于调试
func RedisRemoveRoom(roomid uint64) (err error) {
	const prefix = "RedisRemoveRoom "

	log.Info(prefix, zap.Any("roomid", roomid))
	//清空用户列表
	if err = RedisClient.Del(fmt.Sprintf("users_%d", roomid)).Err(); err != nil {
		log.Error(prefix+"redis.Del", zap.Any("err", err))
	}

	log.Info(prefix+"redis.Del users", zap.Any("roomid", roomid), zap.Any("key", fmt.Sprintf("users_%d", roomid)))

	//清空用户心跳信息
	if err = RedisClient.Del(fmt.Sprintf("users_ping_time_%d", roomid)).Err(); err != nil {
		log.Error(prefix+"redis.Del", zap.Any("err", err))
	}

	log.Info(prefix+"redis.Del heartbeats", zap.Any("roomid", roomid), zap.Any("key", fmt.Sprintf("users_ping_time_%d", roomid)))

	//清空麦序上的用户
	redisClearMikeList(roomid)

	//清空排麦用户
	RedisClearApplyMikeList(roomid)

	//清空被禁止公屏发言的用户
	RedisClearForbidenUserList(roomid)

	//删除房间
	if err = RedisClient.HDel(fmt.Sprintf("bc_roomlist"), fmt.Sprintf("%d", roomid)).Err(); err != nil {
		log.Error(prefix+"redis.HDel", zap.Any("err", err))
	}

	log.Info(prefix+"redis.HDel bc_roomlist", zap.Any("key", fmt.Sprintf("%d", roomid)))
	return
}

//频道信息存redis
func RedisIfRoomExist(roomid uint64) (ret bool, err error) {
	const prefix = "RedisIfRoomExist "

	if ret, err = RedisClient.HExists("bc_roomlist", fmt.Sprintf("%d", roomid)).Result(); err != nil {
		log.Error(prefix+"redis.HExists", zap.Any("err", err))
		return
	}

	log.Debug(prefix, zap.Any("roomid", roomid))
	return
}

func RedisAddRoom(room *entity.Room) (err error) {
	const prefix = "RedisAddRoom "

	jsonBytes, err := json.Marshal(room)
	if err != nil {
		log.Error(prefix+"json.Marshal(room)", zap.Any("err", err))
		return
	}
	err = RedisClient.HSet("bc_roomlist", fmt.Sprintf("%d", room.Roomid), string(jsonBytes)).Err()
	if err != nil {
		log.Error(prefix+"redis.HSet", zap.Any("err", err))
		return
	}

	log.Debug(prefix, zap.Any("room", room))
	return
}

func RedisHscanRoomList(position uint64) (roomlist []*entity.Room, cursor uint64, err error) {
	const prefix = "RedisHsacnRoomList "

	var keys []string
	keys, cursor, err = RedisClient.HScan("bc_roomlist", position, "", 100).Result()
	if err != nil && err != redis.Nil {
		log.Error(prefix+"redis.HScan", zap.Any("err", err))
		return
	}

	for index, item := range keys {
		if index%2 != 0 {
			room := &entity.Room{}
			if marshalErr := json.Unmarshal([]byte(item), room); marshalErr != nil {
				log.Warn(prefix+"Unmarshal failed", zap.Any("room", room), zap.Any("item", item))
				continue
			}

			roomlist = append(roomlist, room)
		}

	}

	log.Debug(prefix, zap.Any("cursor", cursor))
	return
}

func RedisGetRoomIdList() (roomList []uint64, err error) {
	const prefix = "RedisGetRoomList "

	var redisVal []string
	redisVal, err = RedisClient.HKeys("bc_roomlist").Result()
	if err != nil && err != redis.Nil {
		log.Error(prefix+"redis.HSet", zap.Any("err", err))
		return
	}

	for _, value := range redisVal {
		roomid, _ := strconv.Atoi(value)
		roomList = append(roomList, uint64(roomid))
	}

	log.Debug(prefix, zap.Any("roomList", roomList))
	return
}

func RedisGetRoomInfo(roomid uint64) (room *entity.Room, err error) {
	const prefix = "RedisGetRoomInfo "
	room = &entity.Room{}

	defer func(now time.Time) {
		httpmetrics.DefReport("GetRoomInfoByRoomId", 0, now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	var redisVal string
	redisVal, err = RedisClient.HGet("bc_roomlist", fmt.Sprintf("%d", roomid)).Result()
	if err != nil && err != redis.Nil {
		log.Error(prefix+"redis.HGet", zap.Any("roomid", roomid), zap.Any("err", err))
		return nil, err
	}

	if len(redisVal) == 0 {
		log.Info(prefix+"room not find in redis", zap.Any("roomid", roomid))
		return nil, nil
	}

	if err = json.Unmarshal([]byte(redisVal), room); err != nil {
		log.Warn(prefix+"Unmarshal failed", zap.Any("roomid", roomid))
		return nil, err
	}

	log.Debug(prefix, zap.Any("roomid", roomid))
	return room, nil
}

// 频道内用户基本进出事件
func RedisAddUser(roomid uint64, user *entity.User) (err error) {
	const prefix = "RedisAddUser "
	//用户进入频道成功，先写入缓存redis
	now := uint64(time.Now().Unix())
	userBytes, err := json.Marshal(user)
	if err != nil {
		log.Error(prefix+"json.Marshal(user)", zap.Any("err", err))
		return
	}
	err = RedisClient.HSet(fmt.Sprintf("users_%d", roomid), fmt.Sprintf("%d", user.UserID), string(userBytes)).Err()
	if err != nil {
		log.Error(prefix+"redis.HSet", zap.Any("err", err))
		return
	}

	err = RedisClient.HSet(fmt.Sprintf("users_ping_time_%d", roomid), fmt.Sprintf("%d", user.UserID), fmt.Sprintf("%d", now)).Err()
	if err != nil {
		log.Error(prefix+"redis.HSet", zap.Any("err", err))
		return
	}

	log.Debug(prefix, zap.Any("roomid", roomid), zap.Any("user", user))
	return nil
}

func RedisGetUser(roomid uint64, userid uint64) (user *entity.User, err error) {
	const prefix = "RedisGetUser "
	user = &entity.User{}

	defer func(now time.Time) {
		httpmetrics.DefReport("RedisGetUser", 0, now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	val, err := RedisClient.HGet(fmt.Sprintf("users_%d", roomid), fmt.Sprintf("%d", userid)).Result()
	if err != nil && err != redis.Nil {
		log.Error(prefix+"redis.HGet", zap.Any("err", err))
		return
	}

	if len(val) == 0 {
		log.Warn(prefix+"user not find in redis", zap.Any("roomid", roomid), zap.Any("userid", userid))
		return nil, nil
	}

	if err = json.Unmarshal([]byte(val), user); err != nil {
		log.Warn(prefix+"json.Unmarshal error", zap.Any("roomid", roomid), zap.Any("userid", userid))
		return nil, err
	}

	log.Debug(prefix, zap.Any("roomid", roomid), zap.Any("user", user))
	return user, nil
}

func RedisGetUserCount(roomid uint64) (total int64, err error) {
	const prefix = "RedisGetUserCount "

	if total, err = RedisClient.HLen(fmt.Sprintf("users_%d", roomid)).Result(); err != nil {
		return 0, fmt.Errorf("查询redis失败roomid: %d", roomid)
	}

	log.Debug(prefix, zap.Any("roomid", roomid), zap.Any("roomid", total))
	return total, nil
}

func RedisRemoveUser(roomid uint64, userid uint64) (err error) {
	const prefix = "RedisRemoveUser "

	if _, err = RedisClient.HDel(fmt.Sprintf("users_%d", roomid), fmt.Sprintf("%d", userid)).Result(); err != nil {
		log.Error(prefix+"redis.HDel", zap.Any("err", err))
		return
	}

	if _, err = RedisClient.HDel(fmt.Sprintf("users_ping_time_%d", roomid), fmt.Sprintf("%d", userid)).Result(); err != nil {
		log.Error(prefix+"redis.HDel", zap.Any("err", err))
		return
	}

	log.Debug(prefix, zap.Any("roomid", roomid), zap.Any("userid", userid))
	return nil
}

func RedisGetRoomUserList(roomid uint64) (userList []uint64, err error) {
	const prefix = "RedisGetRoomUserList "

	var redisVal map[string]string
	if redisVal, err = RedisClient.HGetAll(fmt.Sprintf("users_%d", roomid)).Result(); err != nil {
		log.Error(prefix+"redis.HDel", zap.Any("err", err))
		return
	}

	for _, value := range redisVal {
		user := &entity.User{}
		if err = json.Unmarshal([]byte(value), user); err != nil {
			log.Warn(prefix, zap.Any("value", value), zap.Any("err", err))
			return nil, err
		}

		if user.Status == entity.StatusUserJoined {
			userList = append(userList, user.UserID)
		}
	}

	log.Debug(prefix, zap.Any("roomid", roomid))
	return userList, nil
}

func RedisGetDisplayedUsers(roomid uint64, num uint32) (userList []*entity.User, err error) {
	const prefix = "RedisGetDisplayedUsers "

	var redisVal map[string]string
	if redisVal, err = RedisClient.HGetAll(fmt.Sprintf("users_%d", roomid)).Result(); err != nil && err != redis.Nil {
		log.Error(prefix+"redis.HDel", zap.Any("err", err))
		return
	}

	for _, value := range redisVal {
		user := &entity.User{}
		if err = json.Unmarshal([]byte(value), user); err != nil {
			log.Warn(prefix, zap.Any("value", value), zap.Any("err", err))
			return nil, err
		}

		userList = append(userList, user)
	}

	sort.Stable(entity.UserSortByBeginJoinTimeSlice(userList))
	log.Debug(prefix, zap.Any("roomid", roomid))

	//等于0表示取全量的数据
	if num == 0 || len(userList) <= int(num) {
		return userList, nil
	}

	return userList[:num], nil
}

// 频道内用户心跳
func RedisSetPingTime(roomid uint64, userid uint64) (exist bool, err error) {
	const prefix = "RedisSetPingTime "

	key := fmt.Sprintf("users_ping_time_%d", roomid)
	field := fmt.Sprintf("%d", userid)

	exist, err = RedisClient.HExists(key, field).Result()
	if err != nil {
		log.Error(prefix+"redis.HExists", zap.Any("err", err))
		return
	}

	if exist {
		err = RedisClient.HSet(key, field, fmt.Sprintf("%d", time.Now().Unix())).Err()
		if err != nil {
			log.Error(prefix+"redis.HSet", zap.Any("err", err))
			return
		}

	}

	log.Debug(prefix, zap.Any("roomid", roomid), zap.Any("userid", userid), zap.Any("exist", exist))
	return exist, nil
}

// 获取频道内所有用户心跳信息
func RedisGetAllPingTimeByRoomid(roomid uint64) (result map[uint64]uint64, err error) {
	const prefix = "RedisGetAllPingTimeByRoomid "

	var redisVal map[string]string
	if redisVal, err = RedisClient.HGetAll(fmt.Sprintf("users_ping_time_%d", roomid)).Result(); err != nil && err != redis.Nil {
		log.Error(prefix+"redis.HGetAll", zap.Any("err", err))
		return
	}

	result = make(map[uint64]uint64)
	for key, value := range redisVal {
		i, e := strconv.Atoi(key)
		if e != nil {
			continue
		}
		j, e := strconv.Atoi(value)
		if e != nil {
			continue
		}
		result[uint64(i)] = uint64(j)
	}

	log.Debug(prefix, zap.Any("roomid", roomid), zap.Any("result", result))
	return
}

// 频道内禁止公屏聊天用户
func RedisGetForbidenStatus(roomid uint64, userid uint64) (ret bool, err error) {
	const prefix = "RedisGetForbidenStatus "

	if ret, err = RedisClient.SIsMember(fmt.Sprintf("users_forbiden_text_%d", roomid), fmt.Sprintf("%d", userid)).Result(); err != nil && err != redis.Nil {
		log.Error(prefix+"redis.SIsMember", zap.Any("err", err))
		return false, err
	}

	log.Debug(prefix, zap.Any("roomid", roomid), zap.Any("userid", userid), zap.Any("ret", ret))
	return ret, nil
}

func RedisSetForbidenStatus(roomid uint64, userid uint64, opt bool) (err error) {
	const prefix = "RedisSetForbidenStatus "

	if opt {
		if _, err = RedisClient.SAdd(fmt.Sprintf("users_forbiden_text_%d", roomid), fmt.Sprintf("%d", userid)).Result(); err != nil {
			log.Error(prefix+"redis.SAdd", zap.Any("err", err))
			return
		}
	} else {
		if _, err = RedisClient.SRem(fmt.Sprintf("users_forbiden_text_%d", roomid), fmt.Sprintf("%d", userid)).Result(); err != nil {
			log.Error(prefix+"redis.SRem", zap.Any("err", err))
			return
		}
	}

	log.Debug(prefix, zap.Any("roomid", roomid), zap.Any("userid", userid), zap.Any("opt", opt))
	return nil
}

func RedisGetForbidenUserList(roomid uint64) (userlist []uint64, err error) {
	const prefix = "RedisGetForbidenUserList "

	var strVal []string
	if strVal, err = RedisClient.SMembers(fmt.Sprintf("users_forbiden_text_%d", roomid)).Result(); err != nil && err != redis.Nil {
		log.Error(prefix+"redis.SMembers", zap.Any("err", err))
		return
	}

	for _, v := range strVal {
		item, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			log.Error(prefix+"ParseUint", zap.Any("err", err))
			return nil, err
		}
		userlist = append(userlist, item)
	}

	log.Debug(prefix, zap.Any("roomid", roomid), zap.Any("userlist", userlist))
	return
}

func RedisClearForbidenUserList(roomid uint64) (err error) {
	const prefix = "RedisClearForbidenUserList "

	if _, err = RedisClient.Del(fmt.Sprintf("users_forbiden_text_%d", roomid)).Result(); err != nil {
		log.Error(prefix+"redis.Del", zap.Any("err", err))
		return
	}

	log.Info(prefix, zap.Any("roomid", roomid))
	return
}

// 频道内上麦用户信息，由于需要和老的版本兼容，所以需要三个hash 来存储麦上用户信息
// 分别是   用户id --> 用户信息
//		   麦位id --> 麦位信息 (空/有人/锁住)
func RedisAddUserToMike(roomid uint64, user *entity.User) (err error) {
	const prefix = "RedisAddUserToMike "

	userBytes, err := json.Marshal(user)
	if err != nil {
		log.Error(prefix+"json.Marshal(user)", zap.Any("err", err))
		return
	}
	if err = RedisClient.HSet(fmt.Sprintf("stage_%d", roomid), fmt.Sprintf("%d", user.UserID), string(userBytes)).Err(); err != nil {
		log.Error(prefix+"redis.HSet", zap.Any("err", err))
		return
	}

	//更新麦位信息
	if err = RedisClient.HSet(fmt.Sprintf("mikeidx_%d", roomid), fmt.Sprintf("%d", user.MikeIndex), fmt.Sprintf("%d", bilin.MikeInfo_USED)).Err(); err != nil {
		log.Error(prefix+"redis.HSet", zap.Any("err", err))
		return
	}

	log.Debug(prefix, zap.Any("roomid", roomid), zap.Any("user", user))
	return
}

func RedisRemoveUserFromMike(roomid uint64, user *entity.User) (err error) {
	const prefix = "RedisRemoveUserFromMike "

	if err = RedisClient.HDel(fmt.Sprintf("stage_%d", roomid), fmt.Sprintf("%d", user.UserID)).Err(); err != nil {
		log.Error(prefix+"redis.HDel", zap.Any("err", err))
		return
	}

	//更新麦位信息
	if err = RedisClient.HSet(fmt.Sprintf("mikeidx_%d", roomid), fmt.Sprintf("%d", user.MikeIndex), fmt.Sprintf("%d", bilin.MikeInfo_EMPTY)).Err(); err != nil {
		log.Error(prefix+"redis.HSet", zap.Any("err", err))
		return
	}

	log.Debug(prefix, zap.Any("roomid", roomid), zap.Any("userid", user.UserID), zap.Any("MikeIndex", user.MikeIndex))
	return
}

func RedisIfUserOnMike(roomid uint64, userid uint64) (ret bool, err error) {
	const prefix = "RedisIfUserOnMike "

	if ret, err = RedisClient.HExists(fmt.Sprintf("stage_%d", roomid), fmt.Sprintf("%d", userid)).Result(); err != nil {
		log.Error(prefix+"redis.HExists", zap.Any("err", err))
		return
	}

	log.Debug(prefix, zap.Any("roomid", roomid), zap.Any("userid", userid), zap.Any("ret", ret))
	return
}

func RedisGetUserOnMike(roomid uint64, userid uint64) (user *entity.User, err error) {
	const prefix = "RedisGetUserOnMike "
	user = &entity.User{}

	val, err := RedisClient.HGet(fmt.Sprintf("stage_%d", roomid), fmt.Sprintf("%d", userid)).Result()
	if err != nil && err != redis.Nil {
		log.Error(prefix+"redis.HGet", zap.Any("err", err))
		return
	}

	if len(val) == 0 {
		log.Warn(prefix+"user not find in redis", zap.Any("roomid", roomid), zap.Any("userid", userid))
		return nil, nil
	}

	if err = json.Unmarshal([]byte(val), user); err != nil {
		log.Warn(prefix, zap.Any("User", userid), zap.Any("not in room", roomid))
		return nil, err
	}

	log.Debug(prefix, zap.Any("roomid", roomid), zap.Any("user", user))
	return user, nil
}

func RedisGetOnMikeUserList(roomid uint64) (userlist []*entity.User, err error) {
	const prefix = "RedisGetOnMikeUserList "

	var redisVal map[string]string
	if redisVal, err = RedisClient.HGetAll(fmt.Sprintf("stage_%d", roomid)).Result(); err != nil && err != redis.Nil {
		log.Error(prefix+"redis.HGetAll", zap.Any("err", err))
		return
	}

	for _, value := range redisVal {
		user := &entity.User{}
		if err = json.Unmarshal([]byte(value), user); err != nil {
			log.Warn(prefix+"json.Unmarshal", zap.Any("value", value), zap.Any("err", err))
			return nil, err
		}

		userlist = append(userlist, user)
	}

	log.Debug(prefix, zap.Any("roomid", roomid), zap.Any("userlist", userlist))
	return
}

func RedisGetOnMikeUserCount(roomid uint64) (num int64, err error) {
	const prefix = "RedisGetOnMikeUserCount "

	if num, err = RedisClient.HLen(fmt.Sprintf("stage_%d", roomid)).Result(); err != nil {
		log.Error(prefix+"redis.HLen", zap.Any("err", err))
		return
	}

	log.Debug(prefix, zap.Any("roomid", roomid), zap.Any("num", num))
	return
}

// 【注意】暂时不提供该接口  因为服务器清麦的时候需要通知给用户
func redisClearMikeList(roomid uint64) (err error) {
	const prefix = "redisClearMikeList "

	if err = RedisClient.Del(fmt.Sprintf("stage_%d", roomid)).Err(); err != nil {
		log.Error(prefix+"redis.Del", zap.Any("err", err))
		return
	}

	if err = RedisClient.Del(fmt.Sprintf("mikeidx_%d", roomid)).Err(); err != nil {
		log.Error(prefix+"redis.Del", zap.Any("err", err))
		return
	}

	log.Info(prefix, zap.Any("roomid", roomid))
	return
}

//锁定/解锁麦位  opt 2: 锁定   0： 解锁,麦位上有人是不能锁定和解锁的
func RedisLockUnlockMikeWheat(roomid uint64, mikeindex uint32, opt bilin.MikeInfo_MIKEWHEATSTATUS) (err error) {
	const prefix = "RedisLockUnlockMikeWheat "

	if opt != bilin.MikeInfo_EMPTY && opt != bilin.MikeInfo_LOCK {
		log.Error(prefix+"unRecognize opt", zap.Any("roomid", roomid), zap.Any("mikeindex", mikeindex), zap.Any("opt", opt))
		return fmt.Errorf("unRecognize opt")
	}

	//更新麦位信息
	if err = RedisClient.HSet(fmt.Sprintf("mikeidx_%d", roomid), fmt.Sprintf("%d", mikeindex), fmt.Sprintf("%d", opt)).Err(); err != nil {
		log.Error(prefix+"redis.HSet", zap.Any("err", err))
		return
	}

	log.Info(prefix, zap.Any("roomid", roomid), zap.Any("mikeindex", mikeindex), zap.Any("opt", opt))
	return
}

//删除某个麦位
func RedisRemoveMikeWheat(roomid uint64, mikeindex uint32) (err error) {
	const prefix = "RedisRemoveMikeWheat "

	if err = RedisClient.HDel(fmt.Sprintf("mikeidx_%d", roomid), fmt.Sprintf("%d", mikeindex)).Err(); err != nil {
		log.Error(prefix+"redis.HDel", zap.Any("err", err))
	}

	log.Info(prefix, zap.Any("roomid", roomid), zap.Any("mikeindex", mikeindex))
	return
}

//删除某个房间的麦位
func RedisClearMikeWheat(roomid uint64) (err error) {
	const prefix = "RedisClearMikeWheat "

	if err = RedisClient.Del(fmt.Sprintf("mikeidx_%d", roomid)).Err(); err != nil {
		log.Error(prefix+"redis.Del", zap.Any("err", err))
		return
	}

	log.Info(prefix, zap.Any("roomid", roomid))
	return
}

func RedisGetMikeWheatStatus(roomid uint64, mikeindex uint32) (status int, err error) {
	const prefix = "RedisGetMikeWheatStatus "

	val, err := RedisClient.HGet(fmt.Sprintf("mikeidx_%d", roomid), fmt.Sprintf("%d", mikeindex)).Result()
	if err != nil && err != redis.Nil {
		log.Error(prefix+"redis.HGet", zap.Any("err", err))
		return 0, err
	}

	log.Debug(prefix, zap.Any("roomid", roomid), zap.Any("mikeindex", mikeindex), zap.Any("val", val))
	return strconv.Atoi(val)
}

//获取麦位状态
func RedisGetAllMikeWheatStatus(roomid uint64) (mikemap map[int]int, err error) {
	const prefix = "RedisGetAllMikeWheatStatus "
	var redisVal map[string]string
	if redisVal, err = RedisClient.HGetAll(fmt.Sprintf("mikeidx_%d", roomid)).Result(); err != nil && err != redis.Nil {
		log.Error(prefix+"redis.HGetAll", zap.Any("err", err))
		return
	}

	mikemap = make(map[int]int)
	for key, value := range redisVal {
		i, e := strconv.Atoi(key)
		if e != nil {
			continue
		}
		j, e := strconv.Atoi(value)
		if e != nil {
			continue
		}
		mikemap[i] = j
	}

	log.Debug(prefix, zap.Any("roomid", roomid), zap.Any("mikemap", mikemap))
	return
}

// 排麦用户操作
func RedisAddUserToApplyMikeList(roomid uint64, user *entity.User) (err error) {
	const prefix = "RedisAddUserToApplyMikeList "

	userBytes, err := json.Marshal(user)
	if err != nil {
		log.Error(prefix+"json.Marshal(user)", zap.Any("err", err))
		return
	}

	if err = RedisClient.RPush(fmt.Sprintf("apply_taling_users_%d", roomid), string(userBytes)).Err(); err != nil {
		log.Error(prefix+"redis.RPush", zap.Any("err", err))
		return
	}

	log.Debug(prefix, zap.Any("roomid", roomid), zap.Any("user", user))
	return
}

func RedisIfUserOnApplyMikeList(roomid uint64, userid uint64) (ret bool) {
	const prefix = "RedisIfUserOnApplyMikeList "

	ret = false
	allInfo, _ := RedisGetApplyMikeUserList(roomid)
	for _, value := range allInfo {
		if value.UserID == userid {
			ret = true
		}
	}

	log.Debug(prefix, zap.Any("roomid", roomid), zap.Any("userid", userid), zap.Any("ret", ret))
	return ret
}

func RedisRemoveUserFromApplyMikeList(roomid uint64, userid uint64) (err error) {
	const prefix = "RedisRemoveUserFromApplyMikeList "

	//由于排麦可以选择用户上麦，所以需要先把数据全取出来，然后再删
	var redisVal []string
	if redisVal, err = RedisClient.LRange(fmt.Sprintf("apply_taling_users_%d", roomid), 0, -1).Result(); err != nil && err != redis.Nil {
		log.Error(prefix+"redis.LRange", zap.Any("err", err))
		return
	}

	for _, value := range redisVal {
		user := &entity.User{}
		if err = json.Unmarshal([]byte(value), user); err != nil {
			log.Warn(prefix+"json.Unmarshal", zap.Any("value", value), zap.Any("err", err))
			return
		}

		if user.UserID == userid {
			if err = RedisClient.LRem(fmt.Sprintf("apply_taling_users_%d", roomid), 0, string(value)).Err(); err != nil {
				log.Error(prefix+"Redis.LRem", zap.Any("err", err))
				return
			}

			log.Info(prefix, zap.Any("roomid", roomid), zap.Any("user", value))
		}

	}

	log.Debug(prefix, zap.Any("roomid", roomid), zap.Any("userid", userid))
	return
}

// LPoP
func RedisGetOneApplyMikeUser(roomid uint64) (user *entity.User, err error) {
	const prefix = "RedisGetOneApplyMikeUser "
	user = &entity.User{}

	var redisVal string
	if redisVal, err = RedisClient.LPop(fmt.Sprintf("apply_taling_users_%d", roomid)).Result(); err != nil && err != redis.Nil {
		log.Error(prefix+"Redis.LPop", zap.Any("err", err))
		return
	}

	if len(redisVal) == 0 {
		log.Warn(prefix, zap.Any("roomid", roomid), zap.Any("redisVal", redisVal))
		return nil, nil
	}

	if err = json.Unmarshal([]byte(redisVal), user); err != nil {
		log.Warn(prefix+"json.Unmarshal ", zap.Any("err", err))
		return nil, err
	}

	log.Debug(prefix, zap.Any("roomid", roomid), zap.Any("user", user))
	return
}

func RedisGetApplyMikeUserCount(roomid uint64) (num int64, err error) {
	const prefix = "RedisGetApplyMikeUserCount "

	if num, err = RedisClient.LLen(fmt.Sprintf("apply_taling_users_%d", roomid)).Result(); err != nil {
		log.Error(prefix+"redis.LLen", zap.Any("err", err))
	}

	log.Debug(prefix, zap.Any("roomid", roomid), zap.Any("num", num))
	return
}

func RedisGetApplyMikeUserList(roomid uint64) (userlist []*entity.User, err error) {
	const prefix = "RedisGetApplyMikeUserList "

	var redisVal []string
	if redisVal, err = RedisClient.LRange(fmt.Sprintf("apply_taling_users_%d", roomid), 0, -1).Result(); err != nil && err != redis.Nil {
		log.Error(prefix+"redis.LRange", zap.Any("err", err))
		return
	}

	for _, value := range redisVal {
		user := &entity.User{}
		if err = json.Unmarshal([]byte(value), user); err != nil {
			log.Warn(prefix+"json.Unmarshal", zap.Any("value", value), zap.Any("err", err))
			return nil, err
		}

		userlist = append(userlist, user)
	}

	log.Debug(prefix, zap.Any("roomid", roomid))
	return
}

func RedisClearApplyMikeList(roomid uint64) (err error) {
	const prefix = "RedisClearApplyMikeList "

	if err = RedisClient.Del(fmt.Sprintf("apply_taling_users_%d", roomid)).Err(); err != nil {
		log.Error(prefix+"redis.Del", zap.Any("err", err))
		return
	}

	log.Info(prefix, zap.Any("roomid", roomid))
	return
}

//对主播粉丝数的缓存，key--value的方式，ttl默认为5分钟
func RedisSetUserFansCount(uid uint64, count uint32) (err error) {
	const prefix = "RedisSetUserFansCount "

	if err = RedisClient.Set(fmt.Sprintf("owner_fans_%d", uid), count, 5*time.Minute).Err(); err != nil {
		log.Error(prefix+"redis.Set", zap.Any("err", err), zap.Any("uid", uid), zap.Any("count", count))
		return
	}

	log.Debug(prefix, zap.Any("uid", uid), zap.Any("count", count))
	return
}

func RedisGetUserFansCount(uid uint64) (count uint32, err error) {
	const prefix = "RedisGetUserFansCount "

	var value string
	if value, err = RedisClient.Get(fmt.Sprintf("owner_fans_%d", uid)).Result(); err != nil && err != redis.Nil {
		log.Error(prefix+"redis.Get", zap.Any("err", err), zap.Any("uid", uid))
		return
	}

	intCount, _ := strconv.Atoi(value)
	log.Debug(prefix, zap.Any("uid", uid), zap.Any("count", intCount))
	return uint32(intCount), nil
}

//主播和房间的映射
func RedisSetUidToRoomId(uid uint64, roomid uint64) (err error) {
	const prefix = "RedisSetUidToRoomId "

	if err = RedisClient.Set(fmt.Sprintf("uid2room_%d", uid), roomid, 0).Err(); err != nil {
		log.Error(prefix+"redis.Set", zap.Any("err", err), zap.Any("uid", uid), zap.Any("roomid", roomid))
		return
	}

	log.Info(prefix, zap.Any("uid", uid), zap.Any("roomid", roomid))
	return
}

//just for test
func RedisChangeRoomMaixuSwitch(roomid uint64) (err error) {
	const prefix = "RedisChangeRoomMaixuSwitch "

	log.Info(prefix, zap.Any("roomid", roomid))
	return
}

//主持人退出房间需要定时上报java，java要求这样的
func RedisAddHostLeaveTooLongTask(room *entity.Room, host uint64) (err error) {
	const prefix = "RedisAddHostLeaveTooLongTask "

	task := &entity.HostEnterLeaveTask{
		RoomId:    room.Roomid,
		RoomType:  room.RoomType2,
		HostId:    host,
		LeaveTime: uint64(time.Now().Unix()),
	}
	taskBytes, err := json.Marshal(task)
	if err != nil {
		log.Error(prefix+"json.Marshal(taskBytes)", zap.Any("err", err))
		return
	}
	err = RedisClient.HSet(fmt.Sprintf("HostEnterLeaveTask"), fmt.Sprintf("%d", room.Roomid), string(taskBytes)).Err()
	if err != nil {
		log.Error(prefix+"redis.HSet", zap.Any("err", err))
		return
	}

	log.Info(prefix, zap.Any("task", task))
	return
}

func RedisRemoveHostLeaveTooLongTask(roomid uint64) (err error) {
	const prefix = "RedisRemoveHostLeaveTooLongTask "

	if _, err = RedisClient.HDel(fmt.Sprintf("HostEnterLeaveTask"), fmt.Sprintf("%d", roomid)).Result(); err != nil {
		log.Error(prefix+"redis.HDel", zap.Any("err", err))
		return
	}

	log.Info(prefix, zap.Any("roomid", roomid))
	return
}

func RedisGetAllHostLeaveTooLongTasks() (tasks []*entity.HostEnterLeaveTask, err error) {
	const prefix = "RedisGetAllHostLeaveTooLongTasks "

	var redisVal map[string]string
	if redisVal, err = RedisClient.HGetAll(fmt.Sprintf("HostEnterLeaveTask")).Result(); err != nil && err != redis.Nil {
		log.Error(prefix+"redis.HDel", zap.Any("err", err))
		return
	}

	for _, value := range redisVal {
		task := &entity.HostEnterLeaveTask{}
		if err = json.Unmarshal([]byte(value), task); err != nil {
			log.Warn(prefix, zap.Any("value", value), zap.Any("err", err))
			continue
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

//用户特权信息
func RedisGetUserHeadgear(uid uint64) (result string, err error) {
	const prefix = "RedisGetUserHeadgear "

	result, err = RedisClient.HGet("user_headgear", fmt.Sprintf("%d", uid)).Result()
	if err != nil && err != redis.Nil {
		log.Error(prefix+"redis.HGet", zap.Any("uid", uid), zap.Any("err", err))
	}

	log.Debug(prefix, zap.Any("uid", uid), zap.Any("result", result))
	return
}

func RedisSetUserHeadgear(uid uint64, info string) (err error) {
	const prefix = "RedisSetUserHeadgear "

	err = RedisClient.HSet("user_headgear", fmt.Sprintf("%d", uid), info).Err()
	if err != nil {
		log.Error(prefix+"redis.HSet", zap.Any("err", err))
	}

	log.Debug(prefix, zap.Any("uid", uid), zap.Any("info", info))
	return
}

func RedisDelUserHeadgear(uid uint64) (err error) {
	const prefix = "RedisDelUserHeadgear "

	if err = RedisClient.HDel(fmt.Sprintf("user_headgear"), fmt.Sprintf("%d", uid)).Err(); err != nil {
		log.Error(prefix+"redis.HDel", zap.Any("err", err))
	}

	log.Debug(prefix, zap.Any("uid", uid))
	return
}
