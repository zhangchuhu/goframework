package client

import (
	"bilin/protocol"
	"code.yy.com/yytars/goframework/tars/servant"
	"context"
	"fmt"
	"testing"
)

var (
	comm    = servant.NewPbCommunicator()
	objName = fmt.Sprintf("bilin.operationManagement.OptManagementPbObj@tcp -h 58.215.138.213 -t 60000 -p 10038")
	client  = bilin.NewOperationManagementServantClient(objName, comm)
)

func TestClient(t *testing.T) {

	req := &bilin.ActDistributionHeadgearRequest{Hinfo: &bilin.HeadgearInfo{
		Uid:        9999999,
		Headgear:   "https://bilinoperationmanagement.bs2dl.yy.com/2fec7fa157e60a8b33c636f9522083fe.jpg",
		Effecttime: 1528989011,
		Expiretime: 1528999011,
	}}
	resp, err := client.ActDistributionHeadgear(context.TODO(), req)

	fmt.Println(resp, err)
	//t.Logf("resp msg:%v", resp, err)

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
