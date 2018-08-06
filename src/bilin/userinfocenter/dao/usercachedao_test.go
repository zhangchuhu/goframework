package dao_test

import (
	"bilin/protocol/userinfocenter"
	"bilin/userinfocenter/dao"
	_ "code.yy.com/yytars/goframework/tars/servant"
	"github.com/go-redis/redis"
	"testing"
	"time"
)

func InitRedis() {
	dao.RedisClient = redis.NewClient(&redis.Options{
		Addr:         "127.0.0.1:6669",
		Password:     "", // no password set
		DB:           0,  // use default DB
		DialTimeout:  1 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		PoolSize:     10,
		PoolTimeout:  5 * time.Second,
	})
}

func TestSetCacheUserInfo(t *testing.T) {
	InitRedis()
	var user userinfocenter.UserInfo
	user.Uid = 1
	user.Sign = "hello world 10"
	if err := dao.SetCacheUserInfo(10, &user); err != nil {
		t.Error("SetCacheUserInfo error:" + err.Error())
	} else {
		t.Logf("SetCacheUserInfo success")
	}
}

func TestDelelteCacheUserInfo(t *testing.T) {
	InitRedis()
	if err := dao.DelelteCacheUserInfo(10); err != nil {
		t.Error("DelelteCacheUserInfo error:" + err.Error())
	} else {
		t.Logf("DelelteCacheUserInfo success")
	}
}

func TestGetCacheUserInfo(t *testing.T) {
	InitRedis()
	if user, err := dao.GetCacheUserInfo(10); err != nil {
		t.Error("GetCacheUserInfo error:" + err.Error())
	} else {
		t.Logf("GetCacheUserInfo success, 10 user:%v", user)
	}
}

func TestSetCacheOpenStatus(t *testing.T) {
	InitRedis()
	if err := dao.SetCacheOpenStatus(8888, 1, "5.0.0", "ios", "127.0.0.1"); err != nil {
		t.Error("SetCacheOpenStatus error:" + err.Error())
	} else {
		t.Logf("SetCacheOpenStatus success")
	}
}

func TestGetCacheOpenStatus(t *testing.T) {
	InitRedis()
	if status, err := dao.GetCacheOpenStatus(10001, "5.0.0", "ios", "127.0.0.1"); err != nil {
		t.Error("GetCacheOpenStatus error:" + err.Error())
	} else {
		t.Logf("GetCacheOpenStatus success, 10001 user:%v", status)
	}

	if status, err := dao.GetCacheOpenStatus(8888, "5.0.0", "ios", "127.0.0.1"); err != nil {
		t.Error("GetCacheOpenStatus error:" + err.Error())
	} else {
		t.Logf("GetCacheOpenStatus success, 8888 user:%v", status)
	}
}
