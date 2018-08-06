package main

import (
	"bilin/protocol"
	"context"
	"fmt"
	"testing"

	"code.yy.com/yytars/goframework/tars/servant"
	"time"
)

var (
	comm                = servant.NewPbCommunicator()
	objName             = fmt.Sprintf("bilin.bcserver2.BCServantObj@tcp -h 58.215.138.213 -t 60000 -p 10020")
	client              = bilin.NewBCServantClient(objName, comm)
	onlineRoomid uint64 = 410261465
	onlineHostid uint64 = 40227988
)

func TestClient(t *testing.T) {

	for a := 0; a < 100; a++ {
		resp, err := client.EnterBroRoom(context.TODO(), &bilin.EnterBroRoomReq{
			&bilin.Header{
				Userid: uint64(a + 10000),
				Roomid: 400000367,
			},
			"test",
			bilin.USERFROM_BROADCAST,
		})
		if err != nil {
			t.Error("EnterBroRoom err", err)
			return
		}
		t.Logf("resp msg:%v", resp)

		time.Sleep(time.Second)
	}

	//
	//resp1, err := client.ExitBroRoom(context.TODO(), &bilin.ExitBroRoomReq{
	//	&bilin.Header{
	//		Userid: 17795537,
	//		Roomid: 400000367,
	//	},
	//})
	//if err != nil {
	//	t.Error("ExitBroRoom err", err)
	//	return
	//}
	//t.Logf("resp msg:%v", resp1)
}

func TestOnlineHostEnterBroRoom(t *testing.T) {
	resp, err := client.EnterBroRoom(context.TODO(), &bilin.EnterBroRoomReq{
		Header: &bilin.Header{
			Roomid: onlineRoomid,
			Userid: onlineHostid,
		},
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", *resp)
}

func TestOnlineAudienceEnterBroRoom(t *testing.T) {
	resp, err := client.EnterBroRoom(context.TODO(), &bilin.EnterBroRoomReq{
		Header: &bilin.Header{
			Roomid: onlineRoomid,
			Userid: 111111,
		},
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", *resp)
}

func TestChangeBroRoomLinkStatus(t *testing.T) {
	resp, err := client.ChangeBroRoomLinkStatus(context.TODO(), &bilin.ChangeBroRoomLinkStatusReq{
		Header: &bilin.Header{
			Roomid: onlineRoomid,
			Userid: onlineHostid,
		},
		Linkstatus: bilin.BaseRoomInfo_OPENLINK,
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", *resp)

}

func TestAudienceEnterBroRoom(t *testing.T) {
	resp, err := client.EnterBroRoom(context.TODO(), &bilin.EnterBroRoomReq{
		Header: &bilin.Header{
			Roomid: 400000367,
			Userid: 17795944,
		},
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", *resp)
}

func TestHostEnterBroRoom(t *testing.T) {
	resp, err := client.EnterBroRoom(context.TODO(), &bilin.EnterBroRoomReq{
		Header: &bilin.Header{
			Roomid: 400000367,
			Userid: 17795535,
		},
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", *resp)
}

func TestHostExitBroRoom(t *testing.T) {
	resp, err := client.ExitBroRoom(context.TODO(), &bilin.ExitBroRoomReq{
		Header: &bilin.Header{
			Roomid: onlineRoomid,
			Userid: onlineHostid,
		},
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", *resp)
}

func TestAudienceExitBroRoom(t *testing.T) {
	resp, err := client.ExitBroRoom(context.TODO(), &bilin.ExitBroRoomReq{
		Header: &bilin.Header{
			Roomid: 400000367,
			Userid: 17795537,
		},
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", *resp)
}

func TestHostChangeBroRoomLinkStatus(t *testing.T) {
	resp, err := client.ChangeBroRoomLinkStatus(context.TODO(), &bilin.ChangeBroRoomLinkStatusReq{
		Header: &bilin.Header{
			Roomid: 400000367,
			Userid: 17795535,
		},
		Linkstatus: bilin.BaseRoomInfo_OPENLINK,
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", *resp)

}

func TestSendRoomMessage(t *testing.T) {
	resp, err := client.SendRoomMessage(context.TODO(), &bilin.SendRoomMessageReq{
		Header: &bilin.Header{
			Roomid: onlineRoomid,
			Userid: 111111,
		},
		Data: []byte(""),
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", *resp)

}

func TestAudienceLinkOperation(t *testing.T) {
	resp, err := client.AudienceLinkOperation(context.TODO(), &bilin.AudienceLinkOperationReq{
		Header: &bilin.Header{
			Roomid: onlineRoomid,
			Userid: 111111,
		},
		Linkop:     bilin.AudienceLinkOperationReq_LINK,
		Micknumber: 1,
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", *resp)

}

func TestHostChangeBroRoomAutoToMikeStatus(t *testing.T) {
	resp, err := client.ChangeBroRoomAutoToMikeStatus(context.TODO(), &bilin.ChangeBroRoomAutoToMikeStatusReq{
		Header: &bilin.Header{
			Roomid: onlineRoomid,
			Userid: onlineHostid,
		},
		Autolink: bilin.BaseRoomInfo_OPENAUTOTOMIKE,
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", *resp)

}

func TestLockUnlockRoomOperation(t *testing.T) {
	resp, err := client.LockUnlockRoomOperation(context.TODO(), &bilin.LockUnlockRoomOperationReq{
		Header: &bilin.Header{
			Roomid: 400000367,
			Userid: 17795535,
		},
		Opt: 0,
		Pwd: "12345678",
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", *resp)

}
