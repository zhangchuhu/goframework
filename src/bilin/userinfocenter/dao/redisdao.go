package dao

import (
	"bilin/userinfocenter/config"
	"errors"
	"github.com/go-redis/redis"
	"time"
)

var (
	redisNotInitErr = errors.New("redis not init")
	RedisClient     *redis.Client
)

func InitRedisDao() error {
	if conf := config.GetAppConfig(); conf != nil {
		RedisClient = redis.NewClient(&redis.Options{
			Addr:         conf.UserInfoRedis,
			Password:     "", // no password set
			DB:           0,  // use default DB
			DialTimeout:  1 * time.Second,
			ReadTimeout:  1 * time.Second,
			WriteTimeout: 1 * time.Second,
			PoolSize:     10,
			PoolTimeout:  5 * time.Second,
		})
		return nil
	}
	return redisNotInitErr
}
