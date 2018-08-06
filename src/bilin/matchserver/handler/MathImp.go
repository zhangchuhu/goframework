package handler

import (
	"bilin/common/onlinepush"
	"bilin/common/onlinequery"
	"bilin/common/thriftpool"
	"bilin/protocol"
	"encoding/json"
	"math/rand"
	"sort"
	"strconv"
	"time"

	"github.com/go-redis/redis"

	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
)

type UserItem struct {
	Uid       uint32 `json:"uid"`       // uid
	Sex       int    `json:"sex"`       // 性别 0:男性 1:女性
	Role      int    `json:"role"`      // 角色 0:普票用户 1：白名单
	MatchType int    `json:"matchType"` // 匹配类型，0:异性 1:同性
	Province  string `json:"province"`  // 地址
	Timestamp int64  `json:"timestamp"` // 时间戳
}

type UserList []UserItem

// UserList 排序算法
func (list UserList) Swap(i, j int)      { list[i], list[j] = list[j], list[i] }
func (list UserList) Len() int           { return len(list) }
func (list UserList) Less(i, j int) bool { return list[i].Timestamp < list[j].Timestamp }
func sortMapByValue(m map[string]string) UserList {
	p := make(UserList, len(m))
	i := 0
	for _, v := range m {
		var user UserItem
		if err := json.Unmarshal([]byte(v), &user); err != nil {
			log.Error("sortMapByValue", zap.Error(err))
		}
		p[i] = user
		i = i + 1
	}
	sort.Sort(p)
	return p
}

// ******************************************************

const (
	Male   = 0 // 男性
	Female = 1 // 女性

	NormalUser = 0 // 普通用户
	WhiteUser  = 1 // 白名单用户

	MatchSex   = 0 // 异性匹配
	MatchNoSex = 1 // 同性匹配

	ProTimeOut = 5000 // 省内超时时间
)

var (
	Conf        AppConfig
	redisClient *redis.Client
	hotLine     thriftpool.Pool
	callRecord  thriftpool.Pool
	spamLevel   thriftpool.Pool
	hotLineData thriftpool.Pool
	cb          MatchCallBackHanlder
)

type OnlineItem struct {
	Male      uint32 `json:"male"`      // 男性当前人数
	Female    uint32 `json:"female"`    // 女性当前人数
	Online    uint32 `json:"online"`    // 当前人数
	Timestamp int64  `json:"timestamp"` // 修改时间戳ms

	RealMale       int64
	RealFemale     int64
	RealOnline     int64
	WaitingMaleO   int64
	WaitingFemaleO int64
	WaitingMaleS   int64
	WaitingFemaleS int64
}

const (
	// 异性匹配
	FemaleWhiteKey = "FemaleWhite" //全国的女白

	FemaleKey    = "Female"
	MaleWhiteKey = "MaleWhite"
	MaleKey      = "Male"

	// 同性匹配
	SameMaleWhiteKey = "SameMaleWhite"
	SameMaleKey      = "SameMale"
	SameFemaleKey    = "SameFemale"
)

type AppConfig struct {
	SentinelAddr      []string
	RedisAddr         string
	OnlinePushURL     string
	OnlineQueryURL    string
	JavaThriftAddr    []string
	ActTaskThriftAddr []string
}

func Init() {
	rand.Seed(time.Now().Unix())

	onlinequery.URL = Conf.OnlineQueryURL
	log.Info("set onlinequery.URL", zap.Any("value", onlinequery.URL))

	onlinepush.URL = Conf.OnlinePushURL
	log.Info("set onlinepush.URL", zap.Any("value", onlinepush.URL))

	if len(Conf.SentinelAddr) == 0 {
		redisClient = redis.NewClient(&redis.Options{
			Addr: Conf.RedisAddr,
		})
		log.Info("set redis addr", zap.Any("value", Conf.RedisAddr))
	} else {
		redisClient = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    Conf.RedisAddr,
			SentinelAddrs: Conf.SentinelAddr,
		})
		log.Info("set redis addr", zap.Any("master", Conf.RedisAddr), zap.Any("sentinel", Conf.SentinelAddr))
	}

	pong, err := redisClient.Ping().Result()
	log.Info("redis PING", zap.Any("pong", pong), zap.Any("err", err))

	hotLine, err = thriftpool.NewChannelPool(0, 1000, CreateHotLineServiceConn)
	log.Info("set thrift addr", zap.Any("hotLine", Conf.JavaThriftAddr), zap.Any("err", err))

	callRecord, err = thriftpool.NewChannelPool(0, 1000, CreateCallRecordServiceConn)
	log.Info("set thrift addr", zap.Any("callRecord", Conf.JavaThriftAddr), zap.Any("err", err))

	spamLevel, err = thriftpool.NewChannelPool(0, 1000, CreateSpamLevelServiceConn)
	log.Info("set thrift addr", zap.Any("meeting", Conf.JavaThriftAddr), zap.Any("err", err))

	hotLineData, err = thriftpool.NewChannelPool(0, 1000, CreateHotLineDataServiceConn)
	log.Info("set thrift addr", zap.Any("hotLineData", Conf.ActTaskThriftAddr), zap.Any("err", err))

	AddComfortWord(0, "心动，总是在下一个瞬间")
	AddComfortWord(1, "每一次相遇，都是一种缘分")
	AddComfortWord(2, "走了她，却迎来了整个世界")
	AddComfortWord(3, "错过，也是一次缘分的擦肩")
	AddComfortWord(4, "时间是个无赖，总藏着期待")
	AddComfortWord(5, "得之，我幸；不得，我命")
	AddComfortWord(6, "走过、路过，也曾错过")
	AddComfortWord(7, "多谢你的绝情，让我学会追寻")
	AddComfortWord(8, "有缘太短暂，比无缘还惨")
	AddComfortWord(9, "好运来，祝你好运来")
	AddComfortWord(10, "邪灵散开，桃花速来")
	AddComfortWord(11, "1234，换个姿势再来一次")
	AddComfortWord(12, "幸亏不是你，陪我到最后")
	AddComfortWord(13, "群众的眼光总是异样的")
	AddComfortWord(14, "真理总是掌握在少数人的手上")
	AddComfortWord(15, "男人哭吧哭吧不是罪")
	AddComfortWord(16, "呵，女人…")
	AddComfortWord(17, "有些人配着配着就来了")
	AddComfortWord(18, "哥等的不是你，是寂寞")
	AddComfortWord(19, "插科打诨都可以，来个妹子行不行")
	AddComfortWord(20, "我看着颗猕猴桃，眼泪突然被引爆")
	AddComfortWord(21, "据说换个头像，能提高拍拖概率哦")
}

//  新增一个用户
func newUser(uid uint32, matchType int, sex int, role int, province string) (string, error) {
	var err error
	var user UserItem
	user.Uid = uid
	user.MatchType = matchType
	user.Sex = sex
	user.Role = role
	user.Timestamp = time.Now().UnixNano() / 1e6
	user.Province = province
	value, err := json.Marshal(user)
	if err != nil {
		log.Error("newUser", zap.Error(err))
	}

	log.Info("newUser enter", zap.Any("value", string(value)))

	AddProvince(user.Province)

	//1、角色判断 UserItem._role //是否是女性白名单
	if user.Role == WhiteUser && user.Sex == Female {
		AddUserItemRedis(user, false)
		return string(value), err
	}

	//2、(1)判断匹配类型
	// UserItem.MatchType 同性：
	if user.MatchType == MatchNoSex {
		AddUserItemRedis(user, false)
		return string(value), err
	}

	Male, _ := GetProvinceMaleStatRedis(user.Province)
	Female, _ := GetProvinceFemaleStatRedis(user.Province)

	log.Info("newUser get number in province", zap.Any("Male", Male), zap.Any("Female", Female))
	// 屏蔽同省匹配，只做大区匹配
	AddUserItemRedis(user, false)
	return string(value), err

	// UserItem.MatchType 异性：
	if Female == 0 || Male == 0 || Female/Male >= 1/3 {
		AddUserItemRedis(user, true)
		return string(value), err
	} else {
		// 进入大区匹配
		AddUserItemRedis(user, false)
		return string(value), err
	}

	return string(value), err
}

func delUser(uid uint32, matchType int, sex int, role int, province string) (bool, error) {
	var user UserItem
	user.Uid = uid
	user.MatchType = matchType
	user.Sex = sex
	user.Role = role
	user.Timestamp = time.Now().UnixNano() / 1e6
	user.Province = province

	DelUserItemRedis(user, true)
	DelUserItemRedis(user, false)

	return true, nil
}

// 根据UserItem生成redis key isPro:是不是省内
func MakeRedisKey(user UserItem, isPro bool) string {
	// 女性白名单
	if user.Role == WhiteUser && user.Sex == Female {
		return "FemaleWhite"
	}
	var RedisKey string
	var RedisKeyPrefix string
	// 异性匹配
	if user.MatchType == MatchSex {
		if user.Role == WhiteUser && user.Sex == Male {
			RedisKeyPrefix = "MaleWhite"
		}
		if user.Role == NormalUser && user.Sex == Male {
			RedisKeyPrefix = "Male"
		}
		if user.Role == NormalUser && user.Sex == Female {
			RedisKeyPrefix = "Female"
		}
	} else {
		// 同性匹配
		if user.Role == WhiteUser && user.Sex == Male {
			RedisKeyPrefix = "SameMaleWhite"
		}
		if user.Role == NormalUser && user.Sex == Male {
			RedisKeyPrefix = "SameMale"
		}
		if user.Role == NormalUser && user.Sex == Female {
			RedisKeyPrefix = "SameFemale"
		}
	}

	// 是否省内
	if isPro {
		RedisKey = RedisKeyPrefix + "_" + user.Province
	} else {
		RedisKey = RedisKeyPrefix
	}

	return RedisKey
}

// 删除redis对应的user isPro:是不是省内
func DelUserItemRedis(user UserItem, isPro bool) {
	userkey := strconv.FormatUint(uint64(user.Uid), 10)
	redisKey := MakeRedisKey(user, isPro)
	if _, err := redisClient.HDel(redisKey, userkey).Result(); err != nil {
		log.Error("DelUserItemRedis", zap.Error(err))
	}
}

// 添加redis对应的user isPro:是不是省内
func AddUserItemRedis(user UserItem, isPro bool) {
	val, err := json.Marshal(user)
	if err != nil {
		log.Error("AddUserItemRedis", zap.Error(err))
		return
	}
	redisKey := MakeRedisKey(user, isPro)
	userkey := strconv.FormatUint(uint64(user.Uid), 10)
	if _, err = redisClient.HSet(redisKey, userkey, val).Result(); err != nil {
		log.Error("AddUserItemRedis", zap.Error(err))
	}
}

// 查找user是否在对应的redis队列中
func ExistUserItemRedis(uid uint32) (exist bool) {
	userkey := strconv.FormatUint(uint64(uid), 10)
	keys := []string{
		FemaleWhiteKey,
		FemaleKey,
		MaleWhiteKey,
		MaleKey,
		SameMaleWhiteKey,
		SameMaleKey,
		SameFemaleKey,
	}
	for _, key := range keys {
		val, err := redisClient.HGet(key, userkey).Result()
		if err == nil && val != "" {
			exist = true
			break
		}
	}
	return
}

func countHash(keys []string) (res int64) {
	for _, key := range keys {
		if cnt, err := redisClient.HLen(key).Result(); err == nil {
			res += cnt
		}
	}
	return
}

// 获取排队中的人数
func CountUserItemRedis() (femaleO, maleO, femaleS, maleS int64) {
	femaleO = countHash([]string{
		FemaleWhiteKey,
		FemaleKey,
	})
	maleO = countHash([]string{
		MaleWhiteKey,
		MaleKey,
	})
	femaleS = countHash([]string{
		SameFemaleKey,
	})
	maleS = countHash([]string{
		SameMaleWhiteKey,
		SameMaleKey,
	})
	return
}

// 获取单个key的list并排序
func GetUserListRedisKey(key string) UserList {
	val, err := redisClient.HGetAll(key).Result()
	if err != nil {
		log.Error("GetUserListRedisKey", zap.Error(err), zap.Any("result", val))
	}
	list := sortMapByValue(val)
	return list
}

// 获取两个key的list，时间排序和Append
func AppendUserListRedisKey(keyOne string, keyTwo string, merge bool) UserList {
	valueOne, err := redisClient.HGetAll(keyOne).Result()
	if err != nil {
		log.Error("AppendUserListRedisKey", zap.Error(err), zap.Any("valueOne", valueOne))
	}
	valueTwo, err := redisClient.HGetAll(keyTwo).Result()
	if err != nil {
		log.Error("AppendUserListRedisKey", zap.Error(err), zap.Any("valueTwo", valueTwo))
	}
	var list UserList
	if merge {
		value := make(map[string]string, len(valueOne)+len(valueTwo))
		for k, v := range valueOne {
			value[k] = v
		}
		for k, v := range valueTwo {
			value[k] = v
		}
		list = sortMapByValue(value)
	} else {
		listOne := sortMapByValue(valueOne)
		listTwo := sortMapByValue(valueTwo)
		list = append(listOne, listTwo...)
	}
	return list
}

// 获取省内人数统计
func GetProvinceMaleStatRedis(province string) (int64, error) {
	redisKey := "Male_" + province
	lenMale, _ := redisClient.HLen(redisKey).Result()

	redisKey = "MaleWhite_" + province
	lenMaleWhite, err := redisClient.HLen(redisKey).Result()

	return lenMale + lenMaleWhite, err
}

func GetProvinceFemaleStatRedis(province string) (int64, error) {
	redisKey := "Female_" + province
	return redisClient.HLen(redisKey).Result()
}

func DelRedis(Key string) {
	Value, _ := redisClient.HGetAll(Key).Result()
	for k, v := range Value {
		redisClient.HDel(Key, k)
		log.Info("DelRedis HDel ", zap.Any("redisKey", Key), zap.Any("key", k), zap.Any("value", v))
	}
}

// 清除redis 数据
func ClearRedisData() {

	DelRedis(FemaleWhiteKey)
	DelRedis(MaleWhiteKey)
	DelRedis(MaleKey)
	DelRedis(FemaleKey)
	DelRedis(SameMaleWhiteKey)
	DelRedis(SameMaleKey)
	DelRedis(SameFemaleKey)

	ProvinceMap, _ := GetProvinceAll()
	for proNo := range ProvinceMap {
		proKey := MaleWhiteKey + "_" + proNo
		DelRedis(proKey)
		proKey = MaleKey + "_" + proNo
		DelRedis(proKey)
		proKey = FemaleKey + "_" + proNo
		DelRedis(proKey)

		proKey = SameMaleWhiteKey + "_" + proNo
		DelRedis(proKey)
		proKey = SameMaleKey + "_" + proNo
		DelRedis(proKey)
		proKey = SameFemaleKey + "_" + proNo
		DelRedis(proKey)

		delete(ProvinceMap, proNo)
	}
}

// 省内超超时判断调到大区匹配处理
func doTimeOut(fromList UserList, isPro bool) UserList {
	// 判断超时情况
	var list UserList
	now := time.Now().UnixNano() / 1e6

	if isPro {
		for _, fromValue := range fromList {
			if duration := now - fromValue.Timestamp; duration > ProTimeOut {
				log.Info("move user from province to global",
					zap.Any("uid", fromValue.Uid),
					zap.Any("user", fromValue),
					zap.Any("duration", duration))
				DelUserItemRedis(fromValue, true)
				AddUserItemRedis(fromValue, false)
			} else {
				list = append(list, fromValue)
			}
		}
	} else {
		list = fromList
	}

	return list
}

// 匹配处理
func DoMatchSex(fromList UserList, toList UserList, matchCount int, isPro bool) {
	startTime := time.Now().UnixNano() / 1e6

	if isPro {
		fromList = doTimeOut(fromList, isPro)
		toList = doTimeOut(toList, isPro)
	}

	toListNo := 0
	for fromIndex, fromValue := range fromList {
		var matchList UserList = nil
		counter := 0
		for i := toListNo; i < len(toList); i++ {
			matchList = append(matchList, toList[i])
			toListNo++
			counter++
			if counter == matchCount {
				log.Info("DoMatchSex ok",
					zap.Any("fromIndex", fromIndex),
					zap.Any("toListNo", toListNo),
					zap.Any("counter", counter),
					zap.Any("fromValue", fromValue),
					zap.Any("matchList", matchList),
					zap.Any("fromList length", len(fromList)),
					zap.Any("toList length", len(toList)),
					zap.Any("matchCount", matchCount),
					zap.Any("province", isPro))
				for _, v := range matchList {
					DelUserItemRedis(v, isPro)
				}
				DelUserItemRedis(fromValue, isPro)
				cb.ReturnMatchOk(fromValue, matchList)
				break
			}
		}
		if counter >= 1 && counter < matchCount {
			log.Info("DoMatchSex ok not enough male",
				zap.Any("fromIndex", fromIndex),
				zap.Any("toListNo", toListNo),
				zap.Any("counter", counter),
				zap.Any("fromValue", fromValue),
				zap.Any("matchList", matchList),
				zap.Any("fromList length", len(fromList)),
				zap.Any("toList length", len(toList)),
				zap.Any("matchCount", matchCount),
				zap.Any("province", isPro))
			for _, v := range matchList {
				DelUserItemRedis(v, isPro)
			}
			DelUserItemRedis(fromValue, isPro)
			cb.ReturnMatchOk(fromValue, matchList)
			break
		}
	}

	endTime := time.Now().UnixNano() / 1e6
	if duration := endTime - startTime; duration >= 900 {
		log.Warn("DoMatchSex cost time too long",
			zap.Any("fromList length", len(fromList)),
			zap.Any("toList length", len(toList)),
			zap.Any("matchCount", matchCount),
			zap.Any("province", isPro),
			zap.Any("duration", duration))
	}
}

// 异性匹配流程
func HandleMatchSex() {
	startTime := time.Now().UnixNano() / 1e6

	ProvinceMap, _ := GetProvinceAll()
	// 同省匹配
	// 1、获取全国女性白名单，遍历每个省的男白，与之匹配
	for proNo := range ProvinceMap {
		fwl := GetUserListRedisKey(FemaleWhiteKey)
		proKey := MaleWhiteKey + "_" + proNo
		mwl := GetUserListRedisKey(proKey)
		DoMatchSex(fwl, mwl, 1, true)
	}

	// 2、遍历每个省的男白+男性(低优先级)和省里的女性+女白(低优先级) 匹配
	for proNo := range ProvinceMap {
		proMaleWhiteKey := MaleWhiteKey + "_" + proNo
		proMaleKey := MaleKey + "_" + proNo
		maleList := AppendUserListRedisKey(proMaleWhiteKey, proMaleKey, true)
		proFemaleKey := FemaleKey + "_" + proNo
		femaleList := AppendUserListRedisKey(proFemaleKey, FemaleWhiteKey, false)
		DoMatchSex(femaleList, maleList, 3, true)
	}

	// 大区匹配
	// 1、获取全国女性白名单和全国男白，与之匹配
	fwl := GetUserListRedisKey(FemaleWhiteKey)
	mwl := GetUserListRedisKey(MaleWhiteKey)
	DoMatchSex(fwl, mwl, 1, false)

	// 2、全国的男白+男性(低优先级)和全国的女性+女白(低优先级) 匹配
	maleList := AppendUserListRedisKey(MaleWhiteKey, MaleKey, true)
	femaleList := AppendUserListRedisKey(FemaleKey, FemaleWhiteKey, false)
	DoMatchSex(femaleList, maleList, 3, false)

	endTime := time.Now().UnixNano() / 1e6
	if duration := endTime - startTime; duration >= 900 {
		log.Warn("HandleMatchSex cost time too long",
			zap.Any("duration", duration))
	}
}

/////////////////////////////////////////////////////////////////
// 处理匹配
func shuffle(list UserList) {
	for n := len(list); n > 0; n-- {
		i := rand.Intn(n)
		list[n-1], list[i] = list[i], list[n-1]
	}
}

func DoMatchNoSex(fromList UserList, matchCount int, isPro bool) {
	if isPro {
		fromList = doTimeOut(fromList, isPro)
	}

	shuffle(fromList)

	counter := 0
	matchList := make(UserList, matchCount)
	for _, toValue := range fromList {
		matchList[counter] = toValue
		counter++
		if counter == matchCount {
			log.Info("DoMatchNoSex ok",
				zap.Any("matchCount", matchCount),
				zap.Any("province", isPro),
				zap.Any("matchList", matchList))
			counter = 0
			for _, v := range matchList {
				DelUserItemRedis(v, isPro)
			}
			cb.ReturnMatchNoSexOk(matchList)
		}
	}
}

//同性匹配流程
func HandleMatchNoSex() {
	ProvinceMap, _ := GetProvinceAll()
	// 同省匹配
	for proNo := range ProvinceMap {
		proMaleWhiteKey := SameMaleWhiteKey + "_" + proNo
		proMaleKey := SameMaleKey + "_" + proNo
		maleList := AppendUserListRedisKey(proMaleWhiteKey, proMaleKey, true)
		DoMatchNoSex(maleList, 2, true)
	}

	for proNo := range ProvinceMap {
		proKey := SameFemaleKey + "_" + proNo
		femaleList := GetUserListRedisKey(proKey)
		DoMatchNoSex(femaleList, 2, true)
	}

	// 大区匹配
	maleList := AppendUserListRedisKey(SameMaleWhiteKey, SameMaleKey, true)
	DoMatchNoSex(maleList, 2, false)

	femaleList := GetUserListRedisKey(SameFemaleKey)
	DoMatchNoSex(femaleList, 2, false)
}

type MatchCallBackIf interface {
	ReturnMatchOk(user UserItem, matchList UserList)
	ReturnMatchNoSexOk(matchList UserList)
}

// 匹配处理协程
func HandleMatchTimer() {
	startTime := time.Now().UnixNano() / 1e6
	for {
		nowTime := time.Now().UnixNano() / 1e6

		if nowTime-startTime >= 1000 {
			startTime = time.Now().UnixNano() / 1e6
			HandleMatchSex()
			HandleMatchNoSex()
		} else {
			time.Sleep(1000 * time.Millisecond)
		}
	}
}

// 心跳检测协程
func HandleTalkingHeartTimer() {
	startTime := time.Now().UnixNano() / 1e6
	for {
		nowTime := time.Now().UnixNano() / 1e6

		if nowTime-startTime >= 1000 {
			startTime = time.Now().UnixNano() / 1e6
			DoTalkingHeart()
		} else {
			time.Sleep(1000 * time.Millisecond)
		}
	}
}

// 在线人数统计协程
func HandleOnlineStatTimer() {
	startTime := time.Now().UnixNano() / 1e6
	for {
		nowTime := time.Now().UnixNano() / 1e6

		if nowTime-startTime >= 1000 {
			startTime = time.Now().UnixNano() / 1e6
			DoOnlinePush()
		} else {
			time.Sleep(30000 * time.Millisecond)
		}
	}
}

func DoOnlineStatPush(userset map[int64]struct{}) {
	if len(userset) == 0 {
		return
	}
	// 在线人数统计
	stat, _ := GetOnlineStat()

	WriteOnlineCount(stat.RealOnline, stat.RealMale, stat.RealFemale, stat.WaitingMaleO, stat.WaitingFemaleO, stat.WaitingMaleS, stat.WaitingFemaleS)

	var uids []int64
	bc := bilin.BroadcastOnlineUserCount{stat.Online, stat.Male, stat.Female}

	for uid := range userset {
		uids = append(uids, uid)
	}

	unicast(uids, &bc, bilin.MaxType_MATCH_MSG, bilin.MinType_MATCH_BROADCASTONLINEUSERCOUNT_MINTYPE)
}

func addUser(userset map[int64]struct{}, userlist UserList) {
	for _, v := range userlist {
		userset[int64(v.Uid)] = struct{}{}
	}
}

func DoOnlinePush() {
	var ul UserList
	users := make(map[int64]struct{})

	ProvinceMap, _ := GetProvinceAll()

	for proNo := range ProvinceMap {

		proKey := MaleWhiteKey + "_" + proNo
		ul = GetUserListRedisKey(proKey)
		addUser(users, ul)

		proMaleKey := MaleKey + "_" + proNo
		ul = GetUserListRedisKey(proMaleKey)
		addUser(users, ul)

		proFemaleKey := FemaleKey + "_" + proNo
		ul = GetUserListRedisKey(proFemaleKey)
		addUser(users, ul)
	}

	ul = GetUserListRedisKey(FemaleWhiteKey)
	addUser(users, ul)

	ul = GetUserListRedisKey(MaleWhiteKey)
	addUser(users, ul)

	ul = GetUserListRedisKey(MaleKey)
	addUser(users, ul)

	ul = GetUserListRedisKey(FemaleKey)
	addUser(users, ul)

	// 同省
	for proNo := range ProvinceMap {
		proMaleWhiteKey := SameMaleWhiteKey + "_" + proNo
		proMaleKey := SameMaleKey + "_" + proNo
		ul = AppendUserListRedisKey(proMaleWhiteKey, proMaleKey, false)
		addUser(users, ul)
	}

	for proNo := range ProvinceMap {
		proKey := SameFemaleKey + "_" + proNo
		ul = GetUserListRedisKey(proKey)
		addUser(users, ul)
	}

	ul = AppendUserListRedisKey(SameMaleWhiteKey, SameMaleKey, false)
	addUser(users, ul)

	ul = GetUserListRedisKey(SameFemaleKey)
	addUser(users, ul)

	DoOnlineStatPush(users)
}
