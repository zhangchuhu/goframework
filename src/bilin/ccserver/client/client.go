package main

import (
	"bilin/protocol"
	"context"
	"fmt"

	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/tars/servant"
)

func main() {
	comm := servant.NewPbCommunicator()
	comm.SetProperty("locator", "tars.tarsregistry.QueryObj@tcp -h 183.36.111.89 -p 17890:tcp -h 58.215.138.213 -p 17890")
	objName := fmt.Sprintf("bilin.dubboproxy.CommonProxy")
	client := bilin.NewDubboProxyClient(objName, comm)

	rsp, err := client.Invoke(context.TODO(), &bilin.DPInvokeReq{
		Service: "com.bilin.user.center.service.IUserCenterService",
		Method:  "queryUserBilinIds",
		Args: []*bilin.DPInvokeArg{
			{
				Type:  "java.util.List",
				Value: `[40373825, 40373826]`,
			},
		},
	})
	fmt.Printf("err: %v\n", err)
	if rsp != nil {
		fmt.Printf("rsp: type %q exception %v value:\n", rsp.Type, rsp.ThrewException)
		fmt.Printf(rsp.Value)
		fmt.Printf("\n")
	}

	rsp, err = client.Invoke(context.TODO(), &bilin.DPInvokeReq{
		Service: "com.yy.bilin.newmms.api.service.ISpamUserLevelService",
		Method:  "getUserSpamLevel",
		Args: []*bilin.DPInvokeArg{
			{
				Type:  "java.lang.Long",
				Value: `40373825`,
			},
		},
	})
	fmt.Printf("err: %v\n", err)
	if rsp != nil {
		fmt.Printf("rsp: type %q exception %v value:\n", rsp.Type, rsp.ThrewException)
		fmt.Printf(rsp.Value)
		fmt.Printf("\n")
	}

	rsp, err = client.Invoke(context.TODO(), &bilin.DPInvokeReq{
		Service: "com.bilin.user.account.service.IUserLoginService",
		Method:  "verifyUserAccessToken",
		Args: []*bilin.DPInvokeArg{
			{
				Type:  "long",
				Value: "17796200",
			},
			{
				Type:  "java.lang.String",
				Value: "70cf08f0cee680ec0b37b9710d7740d1",
			},
		},
	})
	fmt.Printf("err: %v\n", err)
	if rsp != nil {
		fmt.Printf("rsp: type %q exception %v value:\n", rsp.Type, rsp.ThrewException)
		fmt.Printf(rsp.Value)
		fmt.Printf("\n")
	}

	//test1(client)
	//test2(client)
	//test3(client)
	//test4(client)
}

func test1(client bilin.CCServantClient) {
	resp, err := client.GetRandomCallNumberClient(context.TODO(), &bilin.GetRandomCallNumberClientReq{})
	if err != nil {
		appzaplog.Error("GetRandomCallNumberClient err", zap.Error(err))
		return
	}
	fmt.Printf("GetRandomCallNumberClient: %+v\n", resp)
}

func test2(client bilin.CCServantClient) {
	resp, err := client.GetUserCurrentRoom(context.TODO(), &bilin.GetUserCurrentRoomReq{
		Header: &bilin.Header{
			Userid: 17795524,
		},
	})
	if err != nil {
		appzaplog.Error("GetUserCurrentRoom err", zap.Error(err))
		return
	}
	fmt.Printf("GetUserCurrentRoom: %+v\n", resp)
}

func test3(client bilin.CCServantClient) {
	resp, err := client.SendMessageToUser(context.TODO(), &bilin.SendMessageToUserReq{
		Header: &bilin.Header{
			Userid: 17795524,
			Roomid: 400000384,
		},
		ToUserID: []int64{
			17795556,
			2064,
		},
		Data: []byte("Hello!"),
	})
	if err != nil {
		appzaplog.Error("SendMessageToUser err", zap.Error(err))
		return
	}
	fmt.Printf("SendMessageToUser: %+v\n", resp)
}

func test4(client bilin.CCServantClient) {
	resp, err := client.GenerateRoom(context.TODO(), &bilin.GenerateRoomReq{})
	if err != nil {
		appzaplog.Error("GenerateRoom err", zap.Error(err))
		return
	}
	fmt.Printf("GenerateRoom: %+v\n", resp)
}
