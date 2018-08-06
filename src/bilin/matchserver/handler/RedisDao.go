package handler

import (
	cc "bilin/ccserver/handler"
	"bilin/common/onlinequery"
	"bilin/common/thriftpool"
	"bilin/protocol"
	"bilin/thrift/gen-go/hotline"
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
)

const (
	FemaleWhiteListKey = "FemaleWhiteList" //女白
	MaleWhiteListKey   = "MaleWhiteList"   //男白
	MaleFailTimes      = 5
	MatchIdKey         = "MatchIdKey"
	TalkingHeartKey    = "TalkingHeartKey"
	UserOnlineKey      = "UserOnlineKey" // 用户在线基本信息，保存用户的匹配时的json信息

	FirstPlayKey = "FirstPlayKey" // 用户第一次玩随机匹配
	UserPlayKey  = "UserPlayKey"  // 所有玩随机匹配用户

	TalkingUserKey = "TalkingUserKey" // 正在谈话的对像

	RobotKey = "RobotKey" // 机器人用户

	ProvinceKey = "ProvinceKey" // 省份

	ComfortWordKey       = "ComfortWordKey"       // 安慰语
	UserComfortWordNoKey = "UserComfortWordNoKey" // 安慰语序号

	TalkingHeartTimeout = 17000
	OnlineCountDelay    = 15000
)

var (
	lastOnline         uint32
	lastOnlineMale     uint32
	lastOnlineFemale   uint32
	lastOnlineTime     int64
	lastRealOnline     int64
	lastRealMale       int64
	lastRealFemale     int64
	lastWaitingMaleO   int64
	lastWaitingFemaleO int64
	lastWaitingMaleS   int64
	lastWaitingFemaleS int64
)

// 女性白名单
func AddFemaleWhite(uid uint32) (bool, error) {
	userkey := strconv.FormatUint(uint64(uid), 10)
	redisClient.HSet(FemaleWhiteListKey, userkey, MaleFailTimes).Result()
	return true, nil
}

func DelFemaleWhite(uid uint32) (bool, error) {
	userkey := strconv.FormatUint(uint64(uid), 10)
	redisClient.HDel(FemaleWhiteListKey, userkey).Result()
	return true, nil
}

// 男性白名单
func AddMaleWhite(uid uint32) (bool, error) {
	userkey := strconv.FormatUint(uint64(uid), 10)
	cnt, err := redisClient.HIncrBy(MaleWhiteListKey, userkey, 1).Result()
	if err != nil {
		log.Error("AddMaleWhite", zap.Error(err))
	}
	if cnt >= MaleFailTimes {
		log.Info("AddMaleWhite", zap.Any("uid", uid), zap.Any("count", cnt))
	}
	return true, nil
}

func DelMaleWhite(uid uint32) (bool, error) {
	userkey := strconv.FormatUint(uint64(uid), 10)
	if _, err := redisClient.HDel(MaleWhiteListKey, userkey).Result(); err != nil {
		log.Error("DelMaleWhite", zap.Error(err))
	}
	return true, nil
}

func GetMaleWhite(uid uint32) (string, error) {
	userkey := strconv.FormatUint(uint64(uid), 10)
	value, err := redisClient.HGet(MaleWhiteListKey, userkey).Result()
	return value, err
}

// 判断用户是否是白名单
func IsWhite(uid uint32, sex int) (int, error) {
	userkey := strconv.FormatUint(uint64(uid), 10)
	var value string
	var err error
	if sex == Female {
		value, err = redisClient.HGet(FemaleWhiteListKey, userkey).Result()
	} else {
		value, err = redisClient.HGet(MaleWhiteListKey, userkey).Result()
	}

	count, _ := strconv.Atoi(value)

	log.Info("IsWhite",
		zap.Any("uid", uid),
		zap.Any("sex", sex),
		zap.Any("value", value),
		zap.Any("count", count),
		zap.Any("err", err))
	if count >= MaleFailTimes {
		return 1, nil
	} else {
		return 0, nil
	}
}

// matchid
// 生成唯一的matchid
func GenerateMatchid() string {
	val, err := redisClient.HIncrBy(MatchIdKey, "matchkey", 1).Result()
	if err != nil {
		log.Error("GenerateMatchid", zap.Error(err))
	}
	matchid := strconv.FormatUint(uint64(val), 10)
	return matchid
}

// 添加matchId对应的组队json信息
func AddMatchIdValue(matchid string, value string) (bool, error) {
	if _, err := redisClient.HSet(MatchIdKey, matchid, value).Result(); err != nil {
		log.Error("AddMatchIdValue", zap.Error(err))
		return false, err
	}
	return true, nil
}

// 删除matchId对应的组队json信息
func DelMatchIdValue(matchid string) {
	redisClient.HDel(MatchIdKey, matchid)
}

// 获取matchId对应的组队json信息
func GetMatchIdValue(matchid string) (string, error) {
	return redisClient.HGet(MatchIdKey, matchid).Result()
}

func DelMatchIdExpired() {
	val, err := redisClient.HGet(MatchIdKey, "matchkey").Result()
	if err != nil {
		return
	}
	matchid, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return
	}
	matches, _ := redisHGetAll(MatchIdKey, 50000)
	for k := range matches {
		id, err := strconv.ParseInt(k, 10, 64)
		if err != nil {
			continue
		}
		if matchid-id > 100000 {
			redisClient.HDel(MatchIdKey, k)
		}
	}
}

//////////////////////////////////////////////////////////////
// 第一次玩随机匹配记录
func AddFirstPlay(uid uint32) (bool, error) {
	userkey := strconv.FormatUint(uint64(uid), 10)
	startTime := time.Now().UnixNano() / 1e6
	timeValue := strconv.FormatUint(uint64(startTime), 10)
	if _, err := redisClient.HSet(FirstPlayKey, userkey, timeValue).Result(); err != nil {
		log.Error("AddFirstPlay", zap.Error(err))
	}
	return true, nil
}

// 删除第一次玩随机匹配记录
func DelFirstPlay(uid uint32) (bool, error) {
	userkey := strconv.FormatUint(uint64(uid), 10)
	redisClient.HDel(FirstPlayKey, userkey).Result()
	return true, nil
}

// 获取第一次玩随机匹配记录
func GetFirstPlay(uid uint32) (string, error) {
	userkey := strconv.FormatUint(uint64(uid), 10)
	value, err := redisClient.HGet(FirstPlayKey, userkey).Result()
	return value, err
}

///////////////////////////////////////////////////////////
// 心跳
func AddTalkingHeart(uid uint32) (bool, error) {
	userkey := strconv.FormatUint(uint64(uid), 10)
	startTime := time.Now().UnixNano() / 1e6
	timeValue := strconv.FormatUint(uint64(startTime), 10)
	if _, err := redisClient.HSet(TalkingHeartKey, userkey, timeValue).Result(); err != nil {
		log.Error("AddTalkingHeart", zap.Error(err))
	}
	return true, nil
}

// 删除心跳信息
func DelTalkingHeart(uid uint32) (bool, error) {
	userkey := strconv.FormatUint(uint64(uid), 10)
	if _, err := redisClient.HDel(TalkingHeartKey, userkey).Result(); err != nil {
		log.Error("DelTalkingHeart", zap.Error(err))
	}
	return true, nil
}

// 心跳超时处理
func DoTalkingHeart() (bool, error) {
	val, err := redisClient.HGetAll(TalkingHeartKey).Result()
	if err != nil {
		log.Error("DoTalkingHeart", zap.Error(err), zap.Any("TalkingHeartKey", val))
	}
	startTime := time.Now().UnixNano() / 1e6
	for k, v := range val {
		timeMs, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			log.Error("DoTalkingHeart", zap.Error(err), zap.Any("timeMs", timeMs))
		}
		if startTime-timeMs >= int64(TalkingHeartTimeout) {

			// 通知对方通话取消
			uid, err := strconv.ParseUint(k, 10, 32)
			if err != nil {
				log.Error("DoTalkingHeart", zap.Error(err))
				continue
			}
			peeruid, cid, err := GetTalkingUser(uint32(uid))
			if peeruid != 0 {
				bc := bilin.TalkingAction{}
				bc.Operation = 1
				bc.CancelReason = 1

				var uids []int64
				uids = append(uids, int64(uid))
				uids = append(uids, int64(peeruid))
				unicast(uids, &bc, bilin.MaxType_MATCH_MSG, bilin.MinType_MATCH_TALKACTION_MINTYPE)
				log.Info("DoTalkingHeart notify talking peer to cancel",
					zap.Any("uid", uid),
					zap.Any("peeruid", peeruid))

				// 删除正在通话对象
				DelTalkingUser(uint32(uid))

				// 任务系统
				thriftpool.Invoke(HotLineDataService, hotLineData, func(client interface{}) (err error) {
					var ret int32
					c := client.(*hotline.DataServiceClient)
					tasks := []*hotline.TaskReq{
						{int64(uid), "cashRetainTask", "2", int64(cid)},
						{int64(peeruid), "cashRetainTask", "2", int64(cid)},
					}
					ret, err = c.Cancel(context.TODO(), tasks)
					if err != nil {
						log.Error("Talking Start/Cancel Task", zap.Any("err", err), zap.Any("ret", ret))
					}
					return
				})

				// 生成话单
				WriteCallRecord(CallOpEnd, int64(uid), int64(peeruid), int64(cid), CallTypeUnknown)
			}

			// 删除心跳
			DelTalkingHeart(uint32(uid))

			// 获取用户基本信息
			val, err := GetOnlineUser(uint32(uid))
			if err != nil {
				log.Warn("DoTalkingHeart GetOnlineUser", zap.Error(err), zap.Any("uid", uid))
				continue
			}
			DelOnlineUser(uint32(uid))

			var user UserItem
			if err := json.Unmarshal([]byte(val), &user); err != nil {
				log.Error("DoTalkingHeart json.Unmarshal", zap.Error(err), zap.Any("uid", uid))
			}
			// 删除用户redis信息
			DelUserItemRedis(user, true)
			DelUserItemRedis(user, false)

			log.Info("DoTalkingHeart remove user timed out", zap.Any("uid", uid), zap.Any("value", val))
		}
	}
	DelUserPlayExpired()
	DelRobotExpired()
	DelMatchIdExpired()
	return true, nil
}

/////////////////////////////////////////////////////////////
// 用户在线基本信息，保存用户的匹配时的json信息
func AddOnlineUser(uid uint32, value string) (bool, error) {
	userkey := strconv.FormatUint(uint64(uid), 10)
	if _, err := redisClient.HSet(UserOnlineKey, userkey, value).Result(); err != nil {
		log.Error("AddOnlineUser", zap.Error(err), zap.Any("uid", uid), zap.Any("value", value))
		return false, err
	}

	// 如果在机器人那里就删除
	DelRobot(uid)

	return true, nil
}

func DelOnlineUser(uid uint32) (bool, error) {
	userkey := strconv.FormatUint(uint64(uid), 10)

	if _, err := redisClient.HDel(UserOnlineKey, userkey).Result(); err != nil {
		log.Error("DelOnlineUser", zap.Error(err))
	}

	// 用户下线之后，添加机器人
	val, err := GetUserPlay(uid)
	if err != nil {
		log.Error("DelOnlineUser", zap.Error(err))
	} else {
		AddRobot(uid, val)
	}

	return true, nil
}

func GetOnlineUser(uid uint32) (string, error) {
	userkey := strconv.FormatUint(uint64(uid), 10)
	return redisClient.HGet(UserOnlineKey, userkey).Result()
}

func DelOnlineUserAll() {
	Value, _ := redisClient.HGetAll(UserOnlineKey).Result()
	for k := range Value {
		redisClient.HDel(UserOnlineKey, k)
	}
}

func GetOnlineStat() (OnlineItem, error) {
	var online OnlineItem
	now := time.Now().UnixNano() / 1e6

	if now-lastOnlineTime > OnlineCountDelay {
		// Get real male count
		val, err := redisClient.HGetAll(UserOnlineKey).Result()
		if err != nil {
			log.Error("GetOnlineStat", zap.Error(err), zap.Any("UserOnlineKey", val))
		}
		online.Online = uint32(len(val))
		for _, v := range val {
			var user UserItem
			if err := json.Unmarshal([]byte(v), &user); err != nil {
				log.Error("GetOnlineStat", zap.Error(err))
			}
			if user.Sex == Female {
				online.Female += 1
				online.RealFemale += 1
			} else {
				online.Male += 1
				online.RealMale += 1
			}
		}
		if online.Female == 0 {
			online.Female = 1
		}
		// Revise male count
		ratio := float64(online.Male) / float64(online.Female)
		coefficient := 1.0
		switch {
		case ratio <= 3:
			// keep ratio unchanged
		case ratio > 3 && ratio <= 4:
			coefficient = 0.8
		case ratio > 4 && ratio <= 5:
			coefficient = 0.65
		case ratio > 5 && ratio <= 6:
			coefficient = 0.6
		case ratio > 6 && ratio <= 7:
			coefficient = 0.55
		case ratio > 7 && ratio <= 8:
			coefficient = 0.5
		case ratio > 8:
			ratio = 4.0
		}
		ratio = ratio * coefficient
		online.Male = uint32(ratio * float64(online.Female))
		if online.Male == 0 {
			online.Male = 1
		}
		// Get real online count
		var cnt int
		if cnt, err = onlinequery.UserCount(); err != nil {
			cnt = 6235 // fall back count when service unavailable
			log.Warn("onlinequery.UserCount() fail", zap.Any("err", err), zap.Any("fall back count", cnt))
		} else {
			online.RealOnline = int64(cnt)
		}
		lastOnline = uint32(cc.CalcOnlineDisplay(cnt))
		lastOnlineMale = online.Male
		lastOnlineFemale = online.Female
		lastOnlineTime = now
		lastRealOnline = online.RealOnline
		lastRealMale = online.RealMale
		lastRealFemale = online.RealFemale
		// Get waiting count
		lastWaitingFemaleO, lastWaitingMaleO, lastWaitingFemaleS, lastWaitingMaleS = CountUserItemRedis()
	}
	online.Online = lastOnline
	online.Male = lastOnlineMale
	online.Female = lastOnlineFemale
	online.RealOnline = lastRealOnline
	online.RealMale = lastRealMale
	online.RealFemale = lastRealFemale
	online.WaitingFemaleO = lastWaitingFemaleO
	online.WaitingMaleO = lastWaitingMaleO
	online.WaitingFemaleS = lastWaitingFemaleS
	online.WaitingMaleS = lastWaitingMaleS

	return online, nil
}

//////////////////////////////////////////////////////////
// 所有用户基本信息，保存用户的匹配时的json信息
func AddUserPlay(uid uint32, value string) (bool, error) {
	userkey := strconv.FormatUint(uint64(uid), 10)
	if _, err := redisClient.HSet(UserPlayKey, userkey, value).Result(); err != nil {
		log.Error("AddUserPlay", zap.Error(err))
	}
	return true, nil
}

func GetUserPlayHLen() (int64, error) {
	return redisClient.HLen(UserPlayKey).Result()
}

func DelUserPlayExpired() {
	if hlen, err := GetUserPlayHLen(); err != nil || hlen < 500 {
		return
	}
	var (
		user UserItem
	)
	now := time.Now().UnixNano() / 1e6
	val, _ := redisHGetAll(UserPlayKey, 5000)
	for k, v := range val {
		if err := json.Unmarshal([]byte(v), &user); err != nil {
			log.Error("unmarshal error", zap.Error(err), zap.Any("data", v))
			continue
		}
		if now-user.Timestamp > 1200*1000 {
			redisClient.HDel(UserPlayKey, k)
		}
	}
}

func GetUserPlay(uid uint32) (string, error) {
	userkey := strconv.FormatUint(uint64(uid), 10)
	return redisClient.HGet(UserPlayKey, userkey).Result()
}

func GetUserPlayAll(maxcnt int) (res map[string]string, err error) {
	return redisHGetAll(UserPlayKey, maxcnt)
}

func redisHGetAll(hash string, maxcnt int) (res map[string]string, err error) {
	var (
		key    string
		keys   []string
		cursor uint64
	)
	res = make(map[string]string, maxcnt)
	for {
		keys, cursor, err = redisClient.HScan(hash, cursor, "", 100).Result()
		if err != nil {
			break
		}
		for i, k := range keys {
			switch i % 2 {
			case 0:
				key = k
			case 1:
				res[key] = k
			}
		}
		if cursor == 0 {
			break
		}
		if len(res) >= maxcnt {
			break
		}
	}
	return
}

//////////////////////////////////////////////////////////////

// 正在通话的对象。为了talking时异常退出发取消通话通知
func AddTalkingUser(uid, otheruid, cid, talktype uint32) (bool, error) {
	userkey := strconv.FormatUint(uint64(uid), 10)
	uservalue := strconv.FormatUint(uint64(otheruid), 10)
	channel := strconv.FormatUint(uint64(cid), 10)
	starttime := time.Now().UnixNano() / 1e6
	timevalue := strconv.FormatUint(uint64(starttime), 10)
	talktypestr := strconv.FormatUint(uint64(talktype), 10)
	if _, err := redisClient.HSet(TalkingUserKey, userkey, uservalue+","+channel+","+timevalue+","+talktypestr).Result(); err != nil {
		log.Error("AddTalkingUser", zap.Error(err))
	}
	if _, err := redisClient.HSet(TalkingUserKey, uservalue, userkey+","+channel+","+timevalue+","+talktypestr).Result(); err != nil {
		log.Error("AddTalkingUser", zap.Error(err))
	}
	return true, nil
}

func DelTalkingUser(uid uint32) (peeruid, cid, talktype uint32, begin, end int64, err error) {
	userkey := strconv.FormatUint(uint64(uid), 10)
	val, err := redisClient.HGet(TalkingUserKey, userkey).Result()
	if err != nil {
		log.Error("DelTalkingUser user does not exist", zap.Error(err), zap.Any("uid", uid))
		return
	}
	items := strings.Split(val, ",")
	if len(items) > 3 {
		if temp, err := strconv.ParseUint(items[3], 10, 32); err == nil {
			talktype = uint32(temp)
		}
	}
	if len(items) > 2 {
		if temp, err := strconv.ParseInt(items[2], 10, 64); err == nil {
			begin = temp
			end = time.Now().UnixNano() / 1e6
		}
	}
	if len(items) > 1 {
		if temp, err := strconv.ParseUint(items[1], 10, 32); err == nil {
			cid = uint32(temp)
		}
	}
	if len(items) > 0 {
		if temp, err := strconv.ParseUint(items[0], 10, 32); err == nil {
			peeruid = uint32(temp)
		}
		redisClient.HDel(TalkingUserKey, items[0])
	}
	redisClient.HDel(TalkingUserKey, userkey)
	return
}

func GetTalkingUser(uid uint32) (peeruid uint32, cid uint32, err error) {
	userkey := strconv.FormatUint(uint64(uid), 10)
	val, err := redisClient.HGet(TalkingUserKey, userkey).Result()
	if err != nil {
		log.Error("GetTalkingUser user does not exist", zap.Error(err), zap.Any("uid", uid))
		return
	}
	items := strings.Split(val, ",")
	if len(items) > 1 {
		if temp, err := strconv.ParseUint(items[1], 10, 32); err == nil {
			cid = uint32(temp)
		}
	}
	if len(items) > 0 {
		if temp, err := strconv.ParseUint(items[0], 10, 32); err == nil {
			peeruid = uint32(temp)
		}
	}
	return
}

//////////////////////////////////////////////////////////////////
// 机器人
// 机器人信息，保存用户的匹配时的json信息
func AddRobot(uid uint32, value string) (bool, error) {
	if uid == 0 {
		log.Error("AddRobot", zap.Any("uid", uid))
		return false, nil
	}
	var user UserItem
	err := json.Unmarshal([]byte(value), &user)
	if err != nil {
		log.Error("AddRobot", zap.Any("uid", uid), zap.Error(err))
		return false, err
	}

	if user.Sex == Male {
		user.Timestamp = time.Now().UnixNano() / 1e6
		newval, err := json.Marshal(user)
		if err != nil {
			log.Error("AddRobot", zap.Any("uid", uid), zap.Error(err))
			return false, err
		}
		value = string(newval)
		userkey := strconv.FormatUint(uint64(uid), 10)
		if _, err := redisClient.HSet(RobotKey, userkey, value).Result(); err != nil {
			log.Error("AddRobot", zap.Any("uid", uid), zap.Error(err))
			return false, err
		}
	}

	return true, nil
}

func DelRobot(uid uint32) (bool, error) {
	userkey := strconv.FormatUint(uint64(uid), 10)
	if _, err := redisClient.HDel(RobotKey, userkey).Result(); err != nil {
		log.Error("DelRobot", zap.Error(err))
	}
	return true, nil
}

func GetRobotHLen() (int64, error) {
	return redisClient.HLen(RobotKey).Result()
}

// 获取机器人的list并排序
func GetRobotList() UserList {
	val, err := redisHGetAll(RobotKey, 100)
	if err != nil {
		log.Error("GetRobotList", zap.Error(err))
		return nil
	}
	list := sortMapByValue(val)
	// 倒排，并且去掉最后60秒的机器人
	var robot UserList
	now := time.Now().UnixNano() / 1e6
	for i := len(list) - 1; i >= 0; i-- {
		if now-list[i].Timestamp > 60000 {
			robot = append(robot, list[i])
		}
	}
	return robot
}

func DelRobotExpired() {
	if hlen, err := GetRobotHLen(); err != nil || hlen < 500 {
		return
	}
	var (
		user UserItem
	)
	now := time.Now().UnixNano() / 1e6
	val, _ := redisHGetAll(RobotKey, 5000)
	for k, v := range val {
		if err := json.Unmarshal([]byte(v), &user); err != nil {
			log.Error("unmarshal error", zap.Error(err), zap.Any("data", v))
			continue
		}
		if now-user.Timestamp > 1200*1000 {
			redisClient.HDel(RobotKey, k)
		}
	}
}

//////////////////////////////////////////////////////////////////

// 省份信息
func AddProvince(province string) (bool, error) {
	startTime := time.Now().UnixNano() / 1e6
	timeValue := strconv.FormatUint(uint64(startTime), 10)
	redisClient.HSet(ProvinceKey, province, timeValue)
	return true, nil
}

func DelProvince(province string) (bool, error) {
	redisClient.HDel(ProvinceKey, province).Result()
	return true, nil
}

func GetProvinceHLen() (int64, error) {
	return redisClient.HLen(ProvinceKey).Result()
}

func GetProvinceAll() (map[string]string, error) {
	val, err := redisClient.HGetAll(ProvinceKey).Result()
	if err != nil {
		log.Error("GetProvinceAll", zap.Error(err), zap.Any("result", val))
	}
	return val, err
}

//////////////////////////////////////////////////////////////////

// 安慰语信息
func AddComfortWord(no uint32, word string) {
	key := strconv.FormatUint(uint64(no), 10)
	_, err := redisClient.HSet(ComfortWordKey, key, word).Result()
	if err != nil {
		log.Error("AddComfortWord", zap.Any("no", no), zap.Any("word", word), zap.Error(err))
	}
}

func DelComfortWord(no uint32) (bool, error) {
	userkey := strconv.FormatUint(uint64(no), 10)
	redisClient.HDel(ComfortWordKey, userkey).Result()
	return true, nil
}

func GetComfortWordHLen() (int64, error) {
	return redisClient.HLen(ComfortWordKey).Result()
}

func GetComfortWordAll() (map[string]string, error) {
	Value, err := redisClient.HGetAll(ComfortWordKey).Result()
	return Value, err
}

func GetComfortWord(no uint32) (string, error) {
	userkey := strconv.FormatUint(uint64(no), 10)
	return redisClient.HGet(ComfortWordKey, userkey).Result()
}

func GetUserComfortWordNo(uid uint32) (uint32, error) {
	userkey := strconv.FormatUint(uint64(uid), 10)
	incNo, _ := redisClient.HIncrBy(UserComfortWordNoKey, userkey, 1).Result()
	n, _ := GetComfortWordHLen()

	return uint32(incNo) % uint32(n), nil
}
