package dao

import (
	"bilin/chattag/config"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"errors"
	"github.com/go-redis/redis"
)

var (
	redisNotInitErr = errors.New("redis not init")
	RedisClient     redis.UniversalClient
)

func InitRedisDao(ac *config.AppConfig) error {
	if ac != nil {
		RedisClient = redis.NewUniversalClient(&redis.UniversalOptions{
			MasterName: ac.MasterName,
			Addrs:      ac.SentinelAddrs,
			//DialTimeout:  1 * time.Second,
			//ReadTimeout:  1 * time.Second,
			//WriteTimeout: 1 * time.Second,
			//PoolSize:     10,
			//PoolTimeout:  5 * time.Second,
		})
		if err := RedisClient.Ping().Err(); err != nil {
			appzaplog.Error("InitRedisDao Ping err", zap.Error(err))
			return err
		}
		return nil
	}
	return redisNotInitErr
}

//func AddUserTag(tagid, fromuid, touid int64) error {
//	key := fmt.Sprintf("tagnum_%d_%d_%d", tagid, fromuid, touid)
//	RedisClient.GetSet()
//}
