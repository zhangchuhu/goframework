package handler

import (
	"context"
	"testing"

	"bilin/bcserver/config"
	"bilin/bcserver/domain/entity"
	"bilin/protocol"
	"time"
)

var (
	s       *BCServantObj
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
		RedisAddr:           "183.36.122.50:4019",
		JavaThriftAddr:      []string{"221.228.91.178:9090"},
		ActTaskThriftAddr:   "112.25.80.21:22333",
		MsgFilterThriftAddr: "221.228.105.21:19202",
		PushProxyAddr:       "221.228.105.21:11111",
		MysqlAddr:           "bilin:ZG7qEsNi2@tcp(183.36.124.123:6304)/bilin_hongbao?charset=utf8",
	}
)

func init() {
	config.SetTestAppConfig(appconfig)
	s = NewBCServantObj()
}

func TestEnterBroRoom(t *testing.T) {
	resp, err := s.EnterBroRoom(context.TODO(), &bilin.EnterBroRoomReq{
		Header: &bilin.Header{
			Roomid: roomid,
			Userid: uid1,
		},
	})
	resp, err = s.EnterBroRoom(context.TODO(), &bilin.EnterBroRoomReq{
		Header: &bilin.Header{
			Roomid: roomid,
			Userid: uid2,
		},
	})
	resp, err = s.EnterBroRoom(context.TODO(), &bilin.EnterBroRoomReq{
		Header: &bilin.Header{
			Roomid: roomid,
			Userid: uid3,
		},
	})
	resp, err = s.EnterBroRoom(context.TODO(), &bilin.EnterBroRoomReq{
		Header: &bilin.Header{
			Roomid: roomid,
			Userid: uid4,
		},
	})
	resp, err = s.EnterBroRoom(context.TODO(), &bilin.EnterBroRoomReq{
		Header: &bilin.Header{
			Roomid: roomid,
			Userid: uid5,
		},
	})
	resp, err = s.EnterBroRoom(context.TODO(), &bilin.EnterBroRoomReq{
		Header: &bilin.Header{
			Roomid: roomid,
			Userid: uid6,
		},
	})
	resp, err = s.EnterBroRoom(context.TODO(), &bilin.EnterBroRoomReq{
		Header: &bilin.Header{
			Roomid: roomid,
			Userid: uid7,
		},
	})
	resp, err = s.EnterBroRoom(context.TODO(), &bilin.EnterBroRoomReq{
		Header: &bilin.Header{
			Roomid: roomid,
			Userid: uid8,
		},
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", *resp)
}

func TestNewRoom(t *testing.T) {
	room := entity.NewRoom(roomid)

	t.Log(room)
}

func TestEnterBroRoom6(t *testing.T) {
	resp, err := s.EnterBroRoom(context.TODO(), &bilin.EnterBroRoomReq{
		Header: &bilin.Header{
			Roomid: roomid,
			Userid: uid6,
		},
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", *resp)
}

func TestExitBroRoom4(t *testing.T) {
	resp, err := s.ExitBroRoom(context.TODO(), &bilin.ExitBroRoomReq{
		Header: &bilin.Header{
			Roomid: 400000368,
			Userid: 17795536,
		},
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", *resp)

}

func TestHostEnterBroRoom(t *testing.T) {
	resp, err := s.EnterBroRoom(context.TODO(), &bilin.EnterBroRoomReq{
		Header: &bilin.Header{
			Roomid: roomid,
			Userid: hostuid,
		},
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", *resp)
}

func TestCommonCheckAuth(t *testing.T) {
	room, user, err := s.CommonCheckAuth(123456, 567891)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", *room)

	t.Logf("%+v", *user)
}

func TestPingBroRoom(t *testing.T) {
	resp, err := s.PingBroRoom(context.TODO(), &bilin.PingBroRoomReq{
		Header: &bilin.Header{
			Roomid: roomid,
			Userid: userid,
		},
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", *resp)

}

func TestChangeBroRoomLinkStatus(t *testing.T) {
	resp, err := s.ChangeBroRoomLinkStatus(context.TODO(), &bilin.ChangeBroRoomLinkStatusReq{
		Header: &bilin.Header{
			Roomid: roomid,
			Userid: hostuid,
		},
		Linkstatus: bilin.BaseRoomInfo_OPENLINK,
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", *resp)

}

func TestHostChangeBroRoomLinkStatus(t *testing.T) {
	resp, err := s.ChangeBroRoomLinkStatus(context.TODO(), &bilin.ChangeBroRoomLinkStatusReq{
		Header: &bilin.Header{
			Roomid: roomid,
			Userid: hostuid,
		},
		Linkstatus: bilin.BaseRoomInfo_CLOSELINK,
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", *resp)

}

func TestHostChangeBroRoomAutoToMikeStatus(t *testing.T) {
	resp, err := s.ChangeBroRoomAutoToMikeStatus(context.TODO(), &bilin.ChangeBroRoomAutoToMikeStatusReq{
		Header: &bilin.Header{
			Roomid: roomid,
			Userid: hostuid,
		},
		Autolink: bilin.BaseRoomInfo_CLOSEAUTOTOMIKE,
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", *resp)

}

func TestExitBroRoom(t *testing.T) {
	resp, err := s.ExitBroRoom(context.TODO(), &bilin.ExitBroRoomReq{
		Header: &bilin.Header{
			Roomid: roomid,
			Userid: uid1,
		},
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", *resp)

}

func TestHostExitBroRoom(t *testing.T) {
	resp, err := s.ExitBroRoom(context.TODO(), &bilin.ExitBroRoomReq{
		Header: &bilin.Header{
			Roomid: roomid,
			Userid: hostuid,
		},
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", *resp)

}

// AudienceLinkOperation 观众请求麦位、取消麦位
func TestAudienceLinkOperation(t *testing.T) {
	resp, err := s.AudienceLinkOperation(context.TODO(), &bilin.AudienceLinkOperationReq{
		Header: &bilin.Header{
			Roomid: roomid,
			Userid: uid1,
		},
		Linkop:     bilin.AudienceLinkOperationReq_LINK,
		Micknumber: 3,
	})
	time.Sleep(2 * time.Second)
	resp, err = s.AudienceLinkOperation(context.TODO(), &bilin.AudienceLinkOperationReq{
		Header: &bilin.Header{
			Roomid: roomid,
			Userid: uid2,
		},
		Linkop:     bilin.AudienceLinkOperationReq_LINK,
		Micknumber: 1,
	})
	time.Sleep(2 * time.Second)
	resp, err = s.AudienceLinkOperation(context.TODO(), &bilin.AudienceLinkOperationReq{
		Header: &bilin.Header{
			Roomid: roomid,
			Userid: uid3,
		},
		Linkop:     bilin.AudienceLinkOperationReq_LINK,
		Micknumber: 1,
	})
	time.Sleep(2 * time.Second)
	resp, err = s.AudienceLinkOperation(context.TODO(), &bilin.AudienceLinkOperationReq{
		Header: &bilin.Header{
			Roomid: roomid,
			Userid: uid4,
		},
		Linkop:     bilin.AudienceLinkOperationReq_LINK,
		Micknumber: 1,
	})
	time.Sleep(2 * time.Second)
	resp, err = s.AudienceLinkOperation(context.TODO(), &bilin.AudienceLinkOperationReq{
		Header: &bilin.Header{
			Roomid: roomid,
			Userid: uid5,
		},
		Linkop:     bilin.AudienceLinkOperationReq_LINK,
		Micknumber: 1,
	})
	time.Sleep(2 * time.Second)
	resp, err = s.AudienceLinkOperation(context.TODO(), &bilin.AudienceLinkOperationReq{
		Header: &bilin.Header{
			Roomid: roomid,
			Userid: uid6,
		},
		Linkop:     bilin.AudienceLinkOperationReq_LINK,
		Micknumber: 1,
	})
	time.Sleep(2 * time.Second)
	resp, err = s.AudienceLinkOperation(context.TODO(), &bilin.AudienceLinkOperationReq{
		Header: &bilin.Header{
			Roomid: roomid,
			Userid: uid7,
		},
		Linkop:     bilin.AudienceLinkOperationReq_LINK,
		Micknumber: 1,
	})
	time.Sleep(2 * time.Second)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", *resp)

}

func TestAudienceUnLinkOperation(t *testing.T) {
	resp, err := s.AudienceLinkOperation(context.TODO(), &bilin.AudienceLinkOperationReq{
		Header: &bilin.Header{
			Roomid: roomid,
			Userid: uid3,
		},
		Linkop:     bilin.AudienceLinkOperationReq_UNLINK,
		Micknumber: 3,
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", *resp)
}

// 直播间禁麦和开麦  主持人抱听众上下麦
func TestMikeOperation(t *testing.T) {
	resp, err := s.MikeOperation(context.TODO(), &bilin.MikeOperationReq{
		Header: &bilin.Header{
			Roomid: roomid,
			Userid: hostuid,
		},
		Opt:     bilin.MikeOperationReq_UNMIKE,
		Mikeidx: 1,
		Userid:  uid5,
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", *resp)

}
func TestGetBroRoomPreparedAudience(t *testing.T) {
	resp, err := s.GetBroRoomPreparedAudience(context.TODO(), &bilin.GetBroRoomPreparedAudienceReq{
		Header: &bilin.Header{
			Roomid: roomid,
			Userid: hostuid,
		},
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", *resp)

}

func TestKickUser(t *testing.T) {
	resp, err := s.KickUser(context.TODO(), &bilin.KickUserReq{
		Header: &bilin.Header{
			Roomid: roomid,
			Userid: hostuid,
		},
		Kickeduserid: uid1,
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", *resp)

}

func TestBroRoomPraise(t *testing.T) {
	resp, err := s.BroRoomPraise(context.TODO(), &bilin.BroRoomPraiseReq{
		Header: &bilin.Header{
			Roomid: roomid,
			Userid: hostuid,
		},
		PraiseCount: 1000,
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", *resp)

}

func TestMuteUser(t *testing.T) {
	resp, err := s.MuteUser(context.TODO(), &bilin.MuteUserReq{
		Header: &bilin.Header{
			Roomid: roomid,
			Userid: hostuid,
		},
		Muteuserid: uid2,
		Opt:        bilin.MuteUserReq_MUTE,
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", *resp)

}

func TestMuteResult(t *testing.T) {
	resp, err := s.MuteResult(context.TODO(), &bilin.MuteResultReq{
		Header: &bilin.Header{
			Roomid: roomid,
			Userid: uid2,
		},
		Opt: bilin.MuteUserReq_MUTE,
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", *resp)

}

func TestForbiddenUser(t *testing.T) {
	resp, err := s.ForbiddenUser(context.TODO(), &bilin.ForbiddenUserReq{
		Header: &bilin.Header{
			Roomid: roomid,
			Userid: hostuid,
		},
		Forbiddenuserid: uid2,
		Opt:             false,
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", *resp)

}

func TestChangeBroRoomType(t *testing.T) {
	resp, err := s.ChangeBroRoomType(context.TODO(), &bilin.ChangeBroRoomTypeReq{
		Header: &bilin.Header{
			Roomid: roomid,
			Userid: hostuid,
		},
		Roomtype: bilin.BaseRoomInfo_ROOMTYPE_UNKNOW,
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", *resp)

}

func TestSendRoomMessage(t *testing.T) {
	resp, err := s.SendRoomMessage(context.TODO(), &bilin.SendRoomMessageReq{
		Header: &bilin.Header{
			Roomid: roomid,
			Userid: hostuid,
		},
		Data: []byte(""),
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", *resp)

}
