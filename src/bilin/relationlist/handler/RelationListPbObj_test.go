package handler

import (
	"bilin/protocol"
	"bilin/relationlist/config"
	"bilin/relationlist/service"
	"context"
	"testing"
)

var (
	s       *RelationListPbObj
	roomid  uint64 = 400000367
	userid  uint64 = 10010
	uid1    uint64 = 1111
	uid2    uint64 = 2222
	uid3    uint64 = 3333
	uid4    uint64 = 4444
	uid5    uint64 = 5555
	uid6    uint64 = 6666
	uid7    uint64 = 7777
	uid8    uint64 = 8888
	uid9    uint64 = 9999
	hostuid uint64 = 17795535

	appconfig = &config.AppConfig{
		RedisAddr:    "183.36.122.50:4019",
		RabbitMqAddr: []string{},
		MysqlAddr:    "bilin:ZG7qEsNi2@tcp(183.36.124.123:6304)/bilin_hongbao?charset=utf8",
	}
)

func init() {
	config.SetTestAppConfig(appconfig)
	s = NewRelationListPbObj()

	service.RedisInit()
}

func TestRSUserMikeOption(t *testing.T) {
	resp, err := s.RSUserMikeOption(context.TODO(), &bilin.RSUserMikeOptionReq{
		Header: &bilin.Header{
			Roomid: roomid,
			Userid: uid2,
		},
		Owner: hostuid,
		Opt:   bilin.RSUserMikeOptionReq_ONMIKE,
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", *resp)
}
