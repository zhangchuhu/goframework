package dao

import (
	// "bilin/bcserver/domain/entity"
	// "bilin/protocol"
	"bilin/protocol/userinfocenter"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	// "encoding/json"
	"github.com/go-redis/redis"
	"github.com/golang/protobuf/proto"
	"strconv"
	"time"
)

const (
	BILIN_USER_KEY_PRE            = "bilin_userinfocenter:"
	BILIN_OPEN_STATUS_KEY_PRE     = "bilin_useropenstatus:"
	USER_INFO_EXPIRE_TIME         = 6 * 3600 * time.Second
	OPEN_STATUS_TRUE_EXPIRE_TIME  = 7 * 24 * 3600 * time.Second //一个月
	OPEN_STATUS_FALSE_EXPIRE_TIME = 1 * 24 * 3600 * time.Second //一个月
	CHECK_USER_STATUS             = 0
	NOT_FOUND                     = -1
)

func SetCacheUserInfo(uid uint64, userInfo *userinfocenter.UserInfo) error {
	if RedisClient == nil {
		appzaplog.Error("redis not init")
		return redisNotInitErr
	}
	key := BILIN_USER_KEY_PRE + strconv.FormatUint(uid, 10)
	pbBytes, err := proto.Marshal(userInfo)
	if err != nil {
		appzaplog.Error("Marshal userinfo fail,", zap.Error(err), zap.Any("userInfo", *userInfo))
		return err
	}

	redisVal, err := RedisClient.Set(key, pbBytes, USER_INFO_EXPIRE_TIME).Result()
	if err != nil && err != redis.Nil {
		appzaplog.Error("redis.Set", zap.Error(err))
		return err
	}

	appzaplog.Debug("redis.Set", zap.String("redis key", key), zap.Any("userInfo", *userInfo), zap.Any("redis redisVal", redisVal))
	return nil
}

func GetCacheUserInfo(uid uint64) (*userinfocenter.UserInfo, error) {
	if RedisClient == nil {
		appzaplog.Error("redis not init")
		return nil, redisNotInitErr
	}
	key := BILIN_USER_KEY_PRE + strconv.FormatUint(uid, 10)
	redisVal, err := RedisClient.Get(key).Result()
	if err != nil && err != redis.Nil {
		appzaplog.Error("redis.Get eror", zap.String("redis rekey", key), zap.Error(err))
		return nil, err
	}

	if err != nil && err == redis.Nil {
		appzaplog.Debug("redis.Get key no exist", zap.String("redis rekey", key))
		return nil, nil
	}

	var user userinfocenter.UserInfo
	if err = proto.Unmarshal([]byte(redisVal), &user); err != nil {
		appzaplog.Warn("Unmarshal failed", zap.String("key", key), zap.Any("redis redisVal", redisVal), zap.Error(err))
		return nil, err
	}

	appzaplog.Debug("redis.Get", zap.String("redis key", key), zap.Any("user", user))
	return &user, nil

}

func DelelteCacheUserInfo(uid uint64) error {
	if RedisClient == nil {
		appzaplog.Error("redis not init")
		return redisNotInitErr
	}
	key := BILIN_USER_KEY_PRE + strconv.FormatUint(uid, 10)
	redisVal, err := RedisClient.Del(key).Result()
	if err != nil && err != redis.Nil {
		appzaplog.Error("redis.Get eror", zap.String("redis rekey", key), zap.Error(err))
		return err
	}

	if err != nil && err == redis.Nil {
		appzaplog.Debug("redis.Get key no exist", zap.String("redis rekey", key))
	}
	appzaplog.Debug("redis.Del", zap.String("redis key", key), zap.Uint64("uid", uid), zap.Any("redisVal", redisVal))

	return nil
}

func SetCacheOpenStatus(uid uint64, status int32, version, clientType, ip string) error {
	if RedisClient == nil {
		appzaplog.Error("SetCacheOpenStatus, redis not init")
		return redisNotInitErr
	}
	key := BILIN_OPEN_STATUS_KEY_PRE + strconv.FormatUint(uid, 10) + "_" + version + "_" + clientType + "_" + ip

	var ex time.Duration
	if status == CHECK_USER_STATUS {
		ex = OPEN_STATUS_TRUE_EXPIRE_TIME
	} else {
		ex = OPEN_STATUS_FALSE_EXPIRE_TIME
	}

	redisVal, err := RedisClient.Set(key, status, ex).Result()
	if err != nil && err != redis.Nil {
		appzaplog.Error("SetCacheOpenStatus, redis.Set", zap.Error(err))
		return err
	}

	appzaplog.Debug("SetCacheOpenStatus, redis.Set", zap.String("redis key", key), zap.Int32("status", status), zap.Any("redis redisVal", redisVal))
	return nil
}

func GetCacheOpenStatus(uid uint64, version, clientType, ip string) (int32, error) {
	if RedisClient == nil {
		appzaplog.Error("GetCacheOpenStatus redis not init")
		return NOT_FOUND, redisNotInitErr
	}
	key := BILIN_OPEN_STATUS_KEY_PRE + strconv.FormatUint(uid, 10) + "_" + version + "_" + clientType + "_" + ip
	redisVal, err := RedisClient.Get(key).Result()
	if err != nil && err != redis.Nil {
		appzaplog.Error("GetCacheOpenStatus redis.Get eror", zap.String("redis rekey", key), zap.Error(err))
		return NOT_FOUND, err
	}

	if err != nil && err == redis.Nil {
		appzaplog.Debug("GetCacheOpenStatus redis.Get key no exist", zap.String("redis rekey", key))
		return NOT_FOUND, nil
	}

	status, err := strconv.Atoi(redisVal)
	if err != nil {
		appzaplog.Debug("GetCacheOpenStatus atoi err", zap.String("redis rekey", key), zap.Any("redisVal", redisVal))
		return NOT_FOUND, nil
	}

	appzaplog.Debug("GetCacheOpenStatus redis.Get", zap.String("redis key", key), zap.Any("redisVal", redisVal))
	return int32(status), nil

}
