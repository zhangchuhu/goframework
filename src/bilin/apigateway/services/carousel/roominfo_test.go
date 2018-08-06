package carousel_test

import (
	_ "code.yy.com/yytars/goframework/tars/servant"
	"testing"
	//"time"
	"bilin/protocol"
	"code.yy.com/yytars/goframework/tars"
	"context"
)

func TestGetRoominfo(t *testing.T) {
	comm := tars.NewCommunicator()
	roomcenterclient := bilin.NewRoomInfoServantClient("bilin.roominfocenter.RoomInfoCenterServantObj@tcp -h 183.36.111.89 -p 12001  -t 60000", comm)
	info, err := roomcenterclient.LivingRoomsInfo(context.TODO(), &bilin.LivingRoomsInfoReq{})
	if err != nil {
		t.Error("LivingRoomsInfo error:" + err.Error())
	} else {
		t.Logf("LivingRoomsInfo success, info: %v", info)
	}
}

func TestBatchLivingRoomsInfoByHosts(t *testing.T) {
	comm := tars.NewCommunicator()
	roomcenterclient := bilin.NewRoomInfoServantClient("bilin.roominfocenter.RoomInfoCenterServantObj@tcp -h 183.36.111.89 -p 12001  -t 60000", comm)
	info, err := roomcenterclient.BatchLivingRoomsInfoByHosts(context.TODO(), &bilin.BatchLivingRoomsInfoByHostsReq{})
	if err != nil {
		t.Error("BatchLivingRoomsInfoByHosts error:" + err.Error())
	} else {
		t.Logf("BatchLivingRoomsInfoByHosts success, info: %v", info)
	}
}
