package service

import (
	"bilin/relationlist/config"
	"github.com/go-redis/redis"

	"bilin/protocol/userinfocenter"
	"bilin/relationlist/entity"
	"bilin/relationlist/service/redis-lock"
	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	RelationGuestMikeOpt   = "relation_guest_mike_opt_"   //记录用户上下麦时间戳相关信息
	RelationDailyMikeValue = "relation_daily_mike_value_" //记录用户通过上麦计算和主播的亲密值,每天最大值为10

	//统计出结果数据，供客户端使用
	RelationDailyStatisticsValue  = "relation_daily_statistics_value_"  //日榜
	RelationWeeklyStatisticsValue = "relation_weekly_statistics_value_" //周榜
	RelationTotalStatisticsValue  = "relation_total_statistics_value_"  //总榜

	//勋章相关信息存redis
	RelationMedalID = "relation_medal_id" //勋章ID
)

const (
	RELATION_LIST_LOCK_KEY  = "relationlist.lock"
	RELATION_LIST_LOCK_TIME = 30 * time.Second

	DISPATCH_MEDAL_LOCK_KEY  = "dispatchmedal.lock"
	DISPATCH_MEDAL_LOCK_TIME = 60 * time.Second
)

var (
	RedisClient *redis.Client
	syncLock    sync.Mutex
)

func RedisInit() {
	if len(config.GetAppConfig().SentinelAddr) == 0 {
		RedisClient = redis.NewClient(&redis.Options{
			Addr: config.GetAppConfig().RedisAddr,
		})
		log.Info("set redis addr", zap.Any("value", config.GetAppConfig().RedisAddr))
	} else {
		RedisClient = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    config.GetAppConfig().RedisAddr,
			SentinelAddrs: config.GetAppConfig().SentinelAddr,
			DialTimeout:   1 * time.Second,
			ReadTimeout:   1 * time.Second,
			WriteTimeout:  1 * time.Second,
			PoolSize:      10,
			PoolTimeout:   5 * time.Second,
		})
		log.Info("set redis addr", zap.Any("master", config.GetAppConfig().RedisAddr), zap.Any("sentinel", config.GetAppConfig().SentinelAddr))
	}

	pong, err := RedisClient.Ping().Result()
	log.Info("redis PING", zap.Any("pong", pong), zap.Any("err", err))

	log.Info("RedisInit connect redis success")
}

func getCurrentDate() string {
	return strings.Replace(time.Now().String()[0:10], "-", "", -1)
}

func RedisLock(lockKey string, lockTime time.Duration) (redislock *lock.Locker, err error) {
	const prefix = "RedisLock "
	// Obtain a new lock with default settings
	redislock, err = lock.Obtain(RedisClient, lockKey, &lock.Options{LockTimeout: lockTime})
	if err != nil {
		log.Error(prefix+"Obtain failed", zap.Any("err", err))
		return
	} else if redislock == nil {
		log.Error(prefix+"Obtain failed", zap.Any("err", err), zap.Any("redislock", redislock))
		return redislock, fmt.Errorf("%s", "redislock is nil")
	}

	log.Info("RedisLock", zap.Any("redislock", redislock), zap.Any("lockKey", lockKey), zap.Any("lockTime", lockTime))
	return
}

func RedisUnLock(redislock *lock.Locker, lockKey string) {
	redislock.Unlock()

	log.Info("RedisUnLock", zap.Any("redislock", redislock), zap.Any("lockKey", lockKey))
}

//scan all owners by match relation_guest_mike_opt_
func RedisScanOwners(position uint64) (uids []uint64, cursor uint64, err error) {
	const prefix = "RedisScanOwners "

	var keys []string
	keys, cursor, err = RedisClient.Scan(position, RelationGuestMikeOpt+"*", 1000).Result()
	if err != nil && err != redis.Nil {
		log.Error(prefix+"redis.HScan", zap.Any("err", err))
		return
	}

	for _, item := range keys {
		s := strings.Split(item, "_")
		uid, _ := strconv.Atoi(s[len(s)-1])
		uids = append(uids, uint64(uid))
	}

	log.Info(prefix, zap.Any("uids", uids))
	return
}

func RedisGetUserInfo(uid uint64) (userInfo *userinfocenter.UserInfo, err error) {
	const prefix = "RedisGetUserInfo "

	var value string
	if value, err = RedisClient.Get(fmt.Sprintf("relation_userinfo_%d", uid)).Result(); err != nil && err != redis.Nil {
		log.Error(prefix+"redis.Get", zap.Any("err", err), zap.Any("uid", uid))
		return nil, err
	}

	userInfo = &userinfocenter.UserInfo{}
	if err = json.Unmarshal([]byte(value), userInfo); err != nil {
		log.Warn(prefix, zap.Any("value", value), zap.Any("err", err))
		return nil, err
	}

	log.Info(prefix, zap.Any("userInfo", userInfo))
	return
}

func RedisSetUserInfo(userInfo *userinfocenter.UserInfo) (err error) {
	const prefix = "RedisSetUserInfo "

	infoBytes, err := json.Marshal(userInfo)
	if err != nil {
		log.Error(prefix+"json.Marshal(infoBytes)", zap.Any("err", err))
		return
	}

	if err = RedisClient.Set(fmt.Sprintf("relation_userinfo_%d", userInfo.Uid), string(infoBytes), 5*time.Minute).Err(); err != nil {
		log.Error(prefix+"redis.Set", zap.Any("err", err), zap.Any("userInfo", userInfo))
		return
	}

	log.Info(prefix, zap.Any("userInfo", userInfo))
	return
}

//嘉宾上麦 以主播的uid作为key，因为所有的关系都是围绕主播展开的
func RedisUserOnMike(owner uint64, guest_uid uint64, onMikeTime int64) (err error) {
	const prefix = "RedisUserOnMike "

	if owner == guest_uid {
		//主播开播，需要删除列表，初始化一下
		if err = RedisClient.Del(RelationGuestMikeOpt + fmt.Sprintf("%d", owner)).Err(); err != nil {
			log.Error(prefix+"redis.Del", zap.Any("err", err))
		}

		log.Info(prefix+"redis del key: relation_guest_mike_opt_", zap.Any("owner", owner), zap.Any("guest_uid", guest_uid))
		return
	}

	err = RedisClient.HSet(RelationGuestMikeOpt+fmt.Sprintf("%d", owner), fmt.Sprintf("%d", guest_uid), onMikeTime).Err()
	if err != nil {
		log.Error(prefix+"redis.HSet", zap.Any("err", err))
		return
	}

	log.Info(prefix, zap.Any("owner", owner), zap.Any("guest_uid", guest_uid), zap.Any("onMikeTime", onMikeTime))
	return
}

func RedisExpireUserOnMike(owner uint64) (err error) {
	const prefix = "RedisExpireUserOnMike "

	//设置一下TTL 防止脏数据长时间霸占内存
	err = RedisClient.Expire(RelationGuestMikeOpt+fmt.Sprintf("%d", owner), time.Hour*5).Err()
	if err != nil {
		log.Error(prefix+"redis Set TTL", zap.Any("err", err))
		return
	}

	log.Info(prefix, zap.Any("owner", owner))
	return
}

//嘉宾下麦
func RedisUserOffMike(owner uint64, guest_uid uint64) (err error) {
	const prefix = "RedisUserOffMike "

	if owner == guest_uid {
		//主播下播，需要删除列表
		if err = RedisClient.Del(RelationGuestMikeOpt + fmt.Sprintf("%d", owner)).Err(); err != nil {
			log.Error(prefix+"redis.Del", zap.Any("err", err))
		}

		log.Info(prefix+"redis del key: relation_guest_mike_opt_", zap.Any("owner", owner), zap.Any("guest_uid", guest_uid))
		return
	}

	if _, err = RedisClient.HDel(RelationGuestMikeOpt+fmt.Sprintf("%d", owner), fmt.Sprintf("%d", guest_uid)).Result(); err != nil {
		log.Error(prefix+"redis.HDel", zap.Any("err", err))
		return
	}

	log.Info(prefix, zap.Any("owner", owner), zap.Any("guest_uid", guest_uid))
	return
}

//定时任务获取嘉宾麦上信息==》时间戳
func RedisGetUserMikeInfoByOwner(owner uint64) (result map[uint64]int64, err error) {
	const prefix = "RedisGetUserMikeInfoByOwner "

	var redisVal map[string]string
	if redisVal, err = RedisClient.HGetAll(RelationGuestMikeOpt + fmt.Sprintf("%d", owner)).Result(); err != nil {
		log.Error(prefix+"redis.HGetAll", zap.Any("err", err))
		return
	}

	result = make(map[uint64]int64)
	for key, value := range redisVal {
		i, e := strconv.Atoi(key)
		if e != nil {
			continue
		}
		j, e := strconv.Atoi(value)
		if e != nil {
			continue
		}
		result[uint64(i)] = int64(j)
	}

	log.Info(prefix, zap.Any("owner", owner), zap.Any("result", result))
	return
}

//attention : 该数据通过ttl清掉，不要手动去删数据
//记录用户通过上麦计算和主播的亲密值，按天计算，每天上限为10
func RedisSetDailyRelationByMike(owner uint64, guest_uid uint64, relationVal int64) (err error) {
	const prefix = "RedisSetDailyRelationByMike "

	redisKey := RelationDailyMikeValue + fmt.Sprintf("%d", owner) + getCurrentDate()
	err = RedisClient.HSet(redisKey, fmt.Sprintf("%d", guest_uid), relationVal).Err()
	if err != nil {
		log.Error(prefix+"redis.HSet", zap.Any("err", err))
		return
	}

	//设置一下TTL 防止脏数据长时间霸占内存
	err = RedisClient.Expire(redisKey, time.Hour*24).Err()
	if err != nil {
		log.Error(prefix+"redis Set TTL", zap.Any("err", err))
		return
	}

	log.Info(prefix, zap.Any("owner", owner), zap.Any("guest_uid", guest_uid), zap.Any("relationVal", relationVal))
	return
}

//获取用户当天通过上麦和主播互动产生的亲密度 如果大于=10  则不能再增加了
func RedisGetDailyRelationByMike(owner uint64, guest_uid uint64) (relationVal int64, err error) {
	const prefix = "RedisGetDailyRelationByMike "

	redisKey := RelationDailyMikeValue + fmt.Sprintf("%d", owner) + getCurrentDate()
	val, err := RedisClient.HGet(redisKey, fmt.Sprintf("%d", guest_uid)).Result()
	if err != nil {
		log.Error(prefix+"redis.HSet", zap.Any("err", err))
		return 0, err
	}

	if len(val) == 0 {
		log.Warn(prefix+"user not find in redis", zap.Any("owner", owner), zap.Any("guest_uid", guest_uid))
		return 0, nil
	}

	var tmpVal int
	tmpVal, err = strconv.Atoi(val)
	relationVal = int64(tmpVal)
	log.Info(prefix, zap.Any("owner", owner), zap.Any("guest_uid", guest_uid), zap.Any("relationVal", relationVal))
	return
}

//用户勋章相关信息存redis，ttl = 5 Minutes
func RedisSetOwnerMedals(owner uint64, medals map[uint64]int32) (err error) {
	const prefix = "RedisSetUserMedalInfo "

	syncLock.Lock()
	defer syncLock.Unlock()

	redisKey := RelationMedalID
	redisval := make(map[string]interface{})
	for guest_uid, medalId := range medals {
		redisval[fmt.Sprintf("%d", guest_uid)] = fmt.Sprintf("%d", medalId)
	}
	err = RedisClient.HMSet(RelationDailyMikeValue, redisval).Err()
	if err != nil {
		log.Error(prefix+"redis.HSet", zap.Any("err", err))
		return
	}

	//设置一下TTL
	err = RedisClient.Expire(redisKey, time.Minute*5).Err()
	if err != nil {
		log.Error(prefix+"redis Set TTL", zap.Any("err", err))
		return
	}

	log.Info(prefix, zap.Any("owner", owner), zap.Any("medals", medals))
	return
}

func RedisGetOwnerMedals(owner uint64) (result map[uint64]int32, err error) {
	const prefix = "RedisGetOwnerMedals "

	redisKey := RelationMedalID
	var redisVal map[string]string
	if redisVal, err = RedisClient.HGetAll(redisKey).Result(); err != nil {
		log.Error(prefix+"redis.HGetAll", zap.Any("err", err))
		return
	}

	result = make(map[uint64]int32)
	for key, value := range redisVal {
		i, e := strconv.Atoi(key)
		if e != nil {
			continue
		}
		j, e := strconv.Atoi(value)
		if e != nil {
			continue
		}
		result[uint64(i)] = int32(j)
	}

	log.Info(prefix, zap.Any("owner", owner), zap.Any("result", result))
	return
}

//下面都是一些临时缓存，如日榜，7日榜，总榜。。。供客户端获取
//日榜
func RedisAddDailyStaticsticsRelationList(owner uint64, guest_uid uint64, relationVal int64) (err error) {
	const prefix = "RedisAddDailyStaticsticsRelationList "

	redisKey := RelationDailyStatisticsValue + fmt.Sprintf("%d", owner) + getCurrentDate()
	err = RedisClient.ZIncrBy(redisKey, float64(relationVal), fmt.Sprintf("%d", guest_uid)).Err()
	if err != nil {
		log.Error(prefix+"redis.HSet", zap.Any("err", err))
		return
	}

	//设置一下TTL 防止脏数据长时间霸占内存
	err = RedisClient.Expire(redisKey, time.Hour*24).Err()
	if err != nil {
		log.Error(prefix+"redis Set TTL", zap.Any("err", err))
		return
	}

	log.Info(prefix, zap.Any("owner", owner), zap.Any("guest_uid", guest_uid), zap.Any("relationVal", relationVal))
	return
}

func RedisGetDailyStatisticsRelationList(owner uint64) (result *entity.RelationStatistics, err error) {
	const prefix = "RedisGetDailyStatisticsRelationList "
	result = &entity.RelationStatistics{AnchorInfo: &entity.UserRelationInfo{UserID: owner, RelationVal: 0}}

	var redisVal []redis.Z
	redisKey := RelationDailyStatisticsValue + fmt.Sprintf("%d", owner) + getCurrentDate()

	//取全量，因为要计算总数
	if redisVal, err = RedisClient.ZRevRangeWithScores(redisKey, 0, -1).Result(); err != nil && err != redis.Nil {
		log.Error(prefix+"redis.ZRange", zap.Any("err", err))
		return
	}

	for _, val := range redisVal {
		uid, e := strconv.Atoi(val.Member.(string))
		if e != nil {
			continue
		}

		item := &entity.UserRelationInfo{UserID: uint64(uid), RelationVal: int64(val.Score)}
		result.RelationList = append(result.RelationList, item)
		result.AnchorInfo.RelationVal += int64(val.Score)
	}

	log.Info(prefix, zap.Any("owner", owner), zap.Any("redisVal", redisVal), zap.Any("result", result))
	return
}

//周榜
func RedisSetWeeklyStatisticsRelationList(owner uint64, relation_list *entity.RelationStatistics) (err error) {
	const prefix = "RedisSetWeeklyStatisticsRelationList "

	redisKey := RelationWeeklyStatisticsValue + fmt.Sprintf("%d", owner) + getCurrentDate()
	jsonBytes, err := json.Marshal(relation_list)
	if err != nil {
		log.Error(prefix+"json.Marshal(room)", zap.Any("err", err))
		return
	}
	if err = RedisClient.Set(redisKey, string(jsonBytes), 48*time.Hour).Err(); err != nil {
		log.Error(prefix+"redis.Set", zap.Any("err", err), zap.Any("owner", owner), zap.Any("relation_list", relation_list))
		return
	}
	log.Info(prefix, zap.Any("owner", owner), zap.Any("relation_list", relation_list))
	return
}

func RedisGetWeeklyStatisticsRelationList(owner uint64) (result *entity.RelationStatistics, err error) {
	const prefix = "RedisGetWeeklyStatisticsRelationList "

	redisKey := RelationWeeklyStatisticsValue + fmt.Sprintf("%d", owner) + getCurrentDate()

	var redisVal string
	redisVal, err = RedisClient.Get(redisKey).Result()
	if err != nil && err != redis.Nil {
		log.Error(prefix+"redis.Get", zap.Any("err", err), zap.Any("owner", owner))
		return
	}

	if len(redisVal) == 0 {
		log.Info(prefix+"not find in redis", zap.Any("owner", owner))
		return nil, nil
	}

	if err = json.Unmarshal([]byte(redisVal), result); err != nil {
		log.Warn(prefix+"Unmarshal failed", zap.Any("owner", owner))
		return nil, err
	}

	log.Info(prefix, zap.Any("owner", owner), zap.Any("result", result))
	return
}

//总榜  mysql和redis同步更新,每次写数据之后更新expire 保存7天
//每个主播下面所有用户和主播的亲密度的总量，包含送礼和上麦的亲密度之和
func RedisAddOwnerTotalRelation(owner uint64, guest_uid uint64, relationVal int64) (err error) {
	const prefix = "RedisAddOwnerTotalRelation "

	redisKey := RelationTotalStatisticsValue + fmt.Sprintf("%d", owner)
	err = RedisClient.ZIncrBy(redisKey, float64(relationVal), fmt.Sprintf("%d", guest_uid)).Err()
	if err != nil {
		log.Error(prefix+"redis.HSet", zap.Any("err", err))
		return
	}

	//设置一下TTL 防止脏数据长时间霸占内存
	err = RedisClient.Expire(redisKey, time.Hour*24*7).Err()
	if err != nil {
		log.Error(prefix+"redis Set TTL", zap.Any("err", err))
		return
	}

	log.Info(prefix, zap.Any("owner", owner), zap.Any("guest_uid", guest_uid), zap.Any("relationVal", relationVal))
	return
}

func RedisGetTotalStatisticsRelationList(owner uint64) (result *entity.RelationStatistics, err error) {
	const prefix = "RedisGetTotalStatisticsRelationList "
	result = &entity.RelationStatistics{AnchorInfo: &entity.UserRelationInfo{UserID: owner, RelationVal: 0}}

	var redisVal []redis.Z
	redisKey := RelationTotalStatisticsValue + fmt.Sprintf("%d", owner)
	if redisVal, err = RedisClient.ZRevRangeWithScores(redisKey, 0, -1).Result(); err != nil && err != redis.Nil {
		log.Error(prefix+"redis.ZRange", zap.Any("err", err))
		return
	}

	for _, val := range redisVal {
		uid, e := strconv.Atoi(val.Member.(string))
		if e != nil {
			continue
		}

		item := &entity.UserRelationInfo{UserID: uint64(uid), RelationVal: int64(val.Score)}
		result.RelationList = append(result.RelationList, item)
		result.AnchorInfo.RelationVal += int64(val.Score)
	}

	log.Info(prefix, zap.Any("owner", owner), zap.Any("result", result))
	return
}
