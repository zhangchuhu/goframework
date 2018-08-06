package service

import (
	"bilin/relationlist/config"
	"testing"
	"time"
)

var (
	roomid    uint64 = 123456
	userid    uint64 = 567891
	owner     uint64 = 17795535
	now              = time.Now().Unix()
	appconfig        = &config.AppConfig{
		//测试环境redis
		RedisAddr: "183.36.122.50:4019",
	}
)

func init() {
	config.SetTestAppConfig(appconfig)
	RedisInit()
}

func TestRedisScanOwners(t *testing.T) {
	var cursor uint64 = 0
	owners, cursor, err := RedisScanOwners(cursor)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(owners)
}

func TestRedisGetDailyStatisticsRelationList(t *testing.T) {
	ret, err := RedisGetDailyStatisticsRelationList(owner, 0, 1)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(ret)
}
