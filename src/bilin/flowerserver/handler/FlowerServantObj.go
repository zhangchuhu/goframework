package handler

import (
	"bilin/common/onlinepush"
	"bilin/protocol"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/golang/protobuf/proto"

	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/tars"
)

const (
	FlowerLimit  = 5
	FlowerExpire = 86460
)

var (
	redisClient *redis.Client
)

type AppConfig struct {
	RedisAddr      string
	OnlinePushURL  string
	OnlineQueryURL string
}

type FlowerServantObj struct {
}

func NewFlowerServantObj(conf AppConfig) *FlowerServantObj {
	onlinepush.URL = conf.OnlinePushURL
	appzaplog.Info("Set onlinepush.URL", zap.Any("url", onlinepush.URL))

	redisClient = redis.NewClient(&redis.Options{
		Addr:     conf.RedisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	appzaplog.Info("Set redis addr", zap.Any("redis", conf.RedisAddr))

	pong, err := redisClient.Ping().Result()
	appzaplog.Info("redis PING", zap.Any("pong", pong), zap.Any("err", err))

	return &FlowerServantObj{}
}

func getRedisKey(uid uint32) string {
	return fmt.Sprintf("%d", uid) + time.Now().UTC().Format("2006-01-02")
}

func GetUid(ctx context.Context) (uint32, error) {
	if m, ok := tars.FromOutgoingContext(ctx); ok {
		uid, err := strconv.ParseUint(m["uid"], 10, 0)
		if err != nil {
			appzaplog.Error("GetUid", zap.Error(err))
			return 0, err
		}

		return uint32(uid), nil
	}

	return 0, fmt.Errorf("no uid found in content")
}

func (this *FlowerServantObj) QueryUsableFlowerCount(ctx context.Context, r *bilin.QueryUsableFlowerCountRequest) (*bilin.QueryUsableFlowerCountRespone, error) {
	uid, err := GetUid(ctx)
	if err != nil {
		appzaplog.Error("QueryUsableFlowerCount fail", zap.Error(err))
		return &bilin.QueryUsableFlowerCountRespone{
			Count: uint32(0),
		}, err
	}

	value, err := redisClient.Get(getRedisKey(uid)).Result()

	//redis err
	if err != nil {
		//if not existed
		if err == redis.Nil {
			redisClient.Set(getRedisKey(uid), 0, FlowerExpire).Result()
			value = "0"
		} else {
			appzaplog.Error("QueryUsableFlowerCount redisClient.Get fail", zap.String("key", getRedisKey(uid)), zap.Error(err))
			return &bilin.QueryUsableFlowerCountRespone{
				Count: uint32(0),
			}, err
		}
	}

	tmpV, err := strconv.Atoi(value)
	if err != nil {
		tmpV = FlowerLimit
	}

	flower := FlowerLimit - tmpV
	if tmpV > FlowerLimit {
		flower = 0
	}

	appzaplog.Info("QueryUsableFlowerCount", zap.String("key", getRedisKey(uid)), zap.Int("flower count", flower))
	return &bilin.QueryUsableFlowerCountRespone{
		Count: uint32(flower),
	}, nil
}

func (this *FlowerServantObj) SendFlower(ctx context.Context, req *bilin.SendFlowerRequest) (*bilin.SendFlowerRespone, error) {
	uid, err := GetUid(ctx)
	if err != nil {
		appzaplog.Error("SendFlower fail", zap.Error(err))
		return &bilin.SendFlowerRespone{
			Result: 1,
			Count:  uint32(0),
		}, nil
	}

	flws, err := redisClient.Incr(getRedisKey(uid)).Result()
	//redis err
	if err != nil {
		appzaplog.Error("SendFlower redisClient.Incr fail", zap.String("key", getRedisKey(uid)), zap.Error(err))

		return &bilin.SendFlowerRespone{
			Result: 1,
			Count:  uint32(0),
		}, nil
	}

	//no chance now
	if flws > FlowerLimit {
		appzaplog.Info("SendFlower fail. no flower left", zap.Uint32("uid", uid), zap.Uint32("toId", req.ToUser))

		return &bilin.SendFlowerRespone{
			Result: 1,
			Count:  uint32(0),
		}, nil
	}

	//send flower done, send unicast to reciever
	bc := bilin.SendFloweBC{
		FromUserid: uid,
		Count:      1,
	}

	go unicast(req.ToUser, &bc, bilin.MaxType_FLOWER_MSG, bilin.MinType_FLOWER_SENDFLOWERBROCAST_MINTYPE)

	appzaplog.Info("SendFlower done", zap.Uint32("from uid", uid), zap.Uint32("to uid", req.ToUser), zap.Int32("flower left", int32(FlowerLimit-flws)))

	// 客户端已经调用 http 接口增加了花数!
	//clientcenter.AddFlower(1, int64(req.ToUser))

	return &bilin.SendFlowerRespone{
		Result: 0,
		Count:  uint32(FlowerLimit - uint32(flws)),
	}, nil
}

func contains(s []int64, e int64) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func unicast(uid uint32, msg proto.Message, maxtype bilin.MaxType, mintype bilin.MinType_FLOWER) error {
	const prefix = "unicast "

	var body bilin.BcMessageBody
	body.Type = int32(mintype)
	var err error
	if body.Data, err = proto.Marshal(msg); err != nil {
		appzaplog.Error(prefix+"[-]proto.Marshal failed", zap.Any("err", err))
		return err
	}
	pushBody, err := proto.Marshal(&body)
	if err != nil {
		appzaplog.Error(prefix+"[-]proto.Marshal failed", zap.Any("err", err))
		return err
	}

	var uids []int64
	uids = append(uids, int64(uid))
	multiMsg := bilin.MultiPush{
		Msg: &bilin.ServerPush{
			MessageType: int32(maxtype),
			PushBuffer:  pushBody,
		},
		UserIDs: uids,
	}

	var offline []int64
	offline, err = onlinepush.PushToUser(multiMsg)
	if err != nil {
		appzaplog.Error("[-]PushNotifyToUser failed push", zap.Any("err", err))
		return err
	}
	if contains(offline, int64(uid)) {
		err = fmt.Errorf("uid %d is offline", uid)
		appzaplog.Error("[-]PushNotifyToUser failed push", zap.Any("err", err))
		return err
	}

	appzaplog.Debug("[+]PushNotifyToUser success push", zap.Any("uids", uids), zap.Any("minType", bilin.MinType_BC_name[int32(mintype)]))
	return nil
}
