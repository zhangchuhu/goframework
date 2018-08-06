package main

import (
	"bilin/protocol"
	"context"
	"fmt"
	"testing"

	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/tars/servant"
)

func TestMatchRandomCallReq(t *testing.T) {

	comm := servant.NewPbCommunicator()
	//comm.SetProperty("locator", "bilin.matchserver.MatchServantObj@tcp -h 183.36.111.89 -p 13333")
	objName := fmt.Sprintf("bilin.matchserver.MatchServantObj@tcp -h 183.36.111.89 -p 13333")
	client := bilin.NewMatchServantClient(objName, comm)

	for a := 0; a <= 1; a++ {
		var ctxmap = make(map[string]string)
		ctxmap["uid"] = "17795535"
		ctx := servant.NewOutgoingContext(context.TODO(), ctxmap)

		req := &bilin.MatchRandomCallReq{}
		req.Sex = 0
		req.MatchType = 0
		req.Province = "gd"

		resp, err := client.MatchRandomCall(ctx, req)
		if err != nil {
			t.Error("EnterBroRoom err", err)
		}

		t.Logf("resp msg:%v", resp)
	}

	//	{
	//		var ctxmap = make(map[string]string)
	//		ctxmap["uid"] = "17795743"
	//		ctx := servant.NewOutgoingContext(context.TODO(), ctxmap)

	//		req := &bilin.MatchRandomCallReq{}
	//		req.Sex = 0
	//		req.MatchType = 0
	//		req.Province = "gd"

	//		resp, err := client.MatchRandomCall(ctx, req)
	//		if err != nil {
	//			appzaplog.Error("MatchRandomCall err", zap.Error(err))

	//		}
	//		appzaplog.Debug("resp msg", zap.Any("resp", resp.String()))
	//	}

	//	{
	//		var ctxmap = make(map[string]string)
	//		ctxmap["uid"] = "17795745"
	//		ctx := servant.NewOutgoingContext(context.TODO(), ctxmap)

	//		req := &bilin.MatchRandomCallReq{}
	//		req.Sex = 0
	//		req.MatchType = 0
	//		req.Province = "gd"

	//		resp, err := client.MatchRandomCall(ctx, req)
	//		if err != nil {
	//			appzaplog.Error("MatchRandomCall err", zap.Error(err))

	//		}
	//		appzaplog.Debug("resp msg", zap.Any("resp", resp.String()))
	//	}

	//{
	//	var ctxmap = make(map[string]string)
	//	//ctxmap["uid"] = "17795746"
	//	ctxmap["uid"] = "17795817"
	//
	//	ctx := servant.NewOutgoingContext(context.TODO(), ctxmap)
	//
	//	req := &bilin.MatchRandomCallReq{}
	//	req.Sex = 1
	//	req.MatchType = 0
	//	req.Province = "gd"
	//
	//	resp, err := client.MatchRandomCall(ctx, req)
	//	if err != nil {
	//		appzaplog.Error("MatchRandomCall err", zap.Error(err))
	//
	//	}
	//	appzaplog.Debug("resp msg", zap.Any("resp", resp.String()))
	//}

}

func TestCancleMatchRandomReq(t *testing.T) {

	comm := servant.NewPbCommunicator()
	comm.SetProperty("locator", "tars.tarsregistry.QueryObj@tcp -h 58.215.138.213 -p 17890")
	objName := fmt.Sprintf("bilin.matchserver.MatchServantObj")
	client := bilin.NewMatchServantClient(objName, comm)

	var ctxmap = make(map[string]string)
	ctxmap["uid"] = "17795742"
	ctx := servant.NewOutgoingContext(context.TODO(), ctxmap)

	req := &bilin.CancleMatchRandomReq{}
	req.Sex = 0
	req.MatchType = 1
	req.Province = "gd"

	resp, err := client.CancleMatchRandom(ctx, req)
	if err != nil {
		appzaplog.Error("MatchRandomCall err", zap.Error(err))

	}

	appzaplog.Debug("resp msg", zap.Any("resp", resp))
}

func TestSelectMatchingResultReq(t *testing.T) {

	comm := servant.NewPbCommunicator()
	comm.SetProperty("locator", "tars.tarsregistry.QueryObj@tcp -h 58.215.138.213 -p 17890")
	objName := fmt.Sprintf("bilin.matchserver.MatchServantObj")
	client := bilin.NewMatchServantClient(objName, comm)

	var ctxmap = make(map[string]string)
	ctxmap["uid"] = "17795742"
	ctx := servant.NewOutgoingContext(context.TODO(), ctxmap)

	req := &bilin.SelectMatchingResultReq{}
	req.Matchid = "132"
	req.Uid = 12

	resp, err := client.SelectMatchingResult(ctx, req)
	if err != nil {
		appzaplog.Error("SelectMatchingResult err", zap.Error(err))

	}

	appzaplog.Debug("resp msg", zap.Any("resp", resp))
}

func TestTalkingRequest(t *testing.T) {

	comm := servant.NewPbCommunicator()
	comm.SetProperty("locator", "tars.tarsregistry.QueryObj@tcp -h 58.215.138.213 -p 17890")
	objName := fmt.Sprintf("bilin.matchserver.MatchServantObj")
	client := bilin.NewMatchServantClient(objName, comm)

	var ctxmap = make(map[string]string)
	ctxmap["uid"] = "17795742"
	ctx := servant.NewOutgoingContext(context.TODO(), ctxmap)

	req := &bilin.TalkingRequest{}
	req.RequestUid = 456
	req.UnicastUid = 132
	req.Operation = 0

	resp, err := client.Talking(ctx, req)
	if err != nil {
		appzaplog.Error("SelectMatchingResult err", zap.Error(err))

	}

	appzaplog.Debug("resp msg", zap.Any("resp", resp))
}

func TestGetComfortWordRequest(t *testing.T) {

	comm := servant.NewPbCommunicator()
	comm.SetProperty("locator", "tars.tarsregistry.QueryObj@tcp -h 58.215.138.213 -p 17890")
	objName := fmt.Sprintf("bilin.matchserver.MatchServantObj")
	client := bilin.NewMatchServantClient(objName, comm)

	var ctxmap = make(map[string]string)
	ctxmap["uid"] = "17795742"
	ctx := servant.NewOutgoingContext(context.TODO(), ctxmap)

	req := &bilin.GetComfortWordRequest{}
	req.Uid = 456

	resp, err := client.GetComfortWord(ctx, req)
	if err != nil {
		appzaplog.Error("SelectMatchingResult err", zap.Error(err))

	}

	appzaplog.Debug("resp msg", zap.Any("resp", resp))
}

func TestGetRandomAvatarReq(t *testing.T) {

	comm := servant.NewPbCommunicator()
	comm.SetProperty("locator", "tars.tarsregistry.QueryObj@tcp -h 58.215.138.213 -p 17890")
	objName := fmt.Sprintf("bilin.matchserver.MatchServantObj")
	client := bilin.NewMatchServantClient(objName, comm)

	var ctxmap = make(map[string]string)
	ctxmap["uid"] = "17795742"
	ctx := servant.NewOutgoingContext(context.TODO(), ctxmap)

	req := &bilin.GetRandomAvatarReq{}
	req.Uid = 456
	req.Sex = 0

	resp, err := client.GetRandomAvatar(ctx, req)
	if err != nil {
		appzaplog.Error("SelectMatchingResult err", zap.Error(err))

	}

	appzaplog.Debug("resp msg", zap.Any("resp", resp))
}

//func main() {

//	var ctxmap = make(map[string]string)
//	ctxmap["uid"] = "17795742"
//	ctx := servant.NewOutgoingContext(context.TODO(), ctxmap)

//	req := &bilin.MatchRandomCallReq{}
//	req.Sex = 0
//	req.MatchType = 1
//	req.Province = "gd"

//	resp, err := client.MatchRandomCall(context.TODO(), req)
//	if err != nil {
//		appzaplog.Error("MatchRandomCall err", zap.Error(err))
//		return
//	}

//	req := &bilin.CancleMatchRandomReq{}
//	req.Sex = 0
//	req.MatchType = 1
//	req.Province = "gd"
//	resp, err := client.CancleMatchRandom(context.TODO(), req)
//	if err != nil {
//		appzaplog.Error("MatchRandomCall err", zap.Error(err))
//		return
//	}

//	req := &bilin.SelectMatchingResultReq{}
//	req.Matchid = "132"
//	req.Uid = 12
//	resp, err := client.SelectMatchingResult(context.TODO(), req)
//	if err != nil {
//		appzaplog.Error("SelectMatchingResult err", zap.Error(err))
//		return
//	}

//	req := &bilin.TalkingRequest{}
//	req.RequestUid = 456
//	req.UnicastUid = 132
//	req.Operation = 1
//	resp, err := client.Talking(context.TODO(), req)
//	if err != nil {
//		appzaplog.Error("Talking err", zap.Error(err))
//		return
//	}

//	req := &bilin.GetComfortWordRequest{}
//	req.Uid = 456
//	resp, err := client.GetComfortWord(context.TODO(), req)
//	if err != nil {
//		appzaplog.Error("GetComfortWord err", zap.Error(err))
//		return
//	}

//	req := &bilin.GetRandomAvatarReq{}
//	req.Uid = 456
//	req.Sex = 0

//	resp, err := client.GetRandomAvatar(ctx, req)
//	if err != nil {
//		appzaplog.Error("GetRandomAvatar err", zap.Error(err))
//		return
//	}

//	appzaplog.Debug("resp msg", zap.Any("resp", resp))
//}
