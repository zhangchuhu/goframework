package main

import (
	"log"
	"os"
	"strconv"

	"github.com/go-redis/redis"
)

type UserItem struct {
	Uid       uint32 `json:"uid"`       // uid
	Sex       int    `json:"sex"`       // 性别 0:男性 1:女性
	Role      int    `json:"role"`      // 角色 0:普票用户 1：白名单
	MatchType int    `json:"matchType"` // 匹配类型，0:异性 1:同性
	Province  string `json:"province"`  // 地址
	Timestamp int64  `json:"timestamp"` // 时间戳
}

var (
	redisAddr   string
	redisClient *redis.Client
	hash        string
	count       int
)

func init() {
	var err error

	if len(os.Args) < 4 {
		log.Fatalf("Call me as: %s redis-addr hash count", os.Args[0])
	}

	redisAddr = os.Args[1]
	hash = os.Args[2]
	if count, err = strconv.Atoi(os.Args[3]); err != nil {
		log.Fatalf("arg 3 (count) must be an integer")
	}
	if count <= 0 {
		log.Fatalf("arg 3 (count) must > 0")
	}

	//redisClient = redis.NewClient(&redis.Options{
	//	Addr: redisAddr,
	//})
	redisClient = redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName: redisAddr,
		SentinelAddrs: []string{
			"group1003-wx-sentinel.yy.com:20098",
			"group1003-sz-sentinel.yy.com:20098",
			"group1003-bj-sentinel.yy.com:20098",
		},
	})
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

func redisHDel(hash string, maxcnt int) {
	var (
		MatchIdKey = hash
		delCount   int
	)
	val, err := redisClient.HGet(MatchIdKey, "matchkey").Result()
	if err != nil {
		return
	}
	matchid, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return
	}
	matches, _ := redisHGetAll(MatchIdKey, maxcnt)
	for k := range matches {
		id, err := strconv.ParseInt(k, 10, 64)
		if err != nil {
			continue
		}
		if matchid-id > 100000 {
			redisClient.HDel(MatchIdKey, k)
			delCount++
		}
	}
	log.Printf("delete %q: %d / %d", hash, delCount, len(matches))
}

func main() {
	redisHDel(hash, count)
}
