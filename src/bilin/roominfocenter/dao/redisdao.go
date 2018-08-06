package dao

import (
	"bilin/roominfocenter/config"
	"errors"
	"github.com/go-redis/redis"
)

const (
// RedisAddr 是在dbms上申请的。TODO 放到配置文件
//RedisAddr = "183.36.122.50:4019"
)

var (
	RedisClient     *redis.Client
	redisNotInitErr = errors.New("redis not init")
)

func InitRedisDao() error {
	if conf := config.GetAppConfig(); conf != nil {
		RedisClient = redis.NewClient(&redis.Options{
			Addr:     conf.RoomInfoRedis,
			Password: "", // no password set
			DB:       0,  // use default DB
		})
		return nil
	}
	return redisNotInitErr
}
