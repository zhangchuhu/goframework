package handler

import (
	"bilin/bcserver/bccommon"
	"bilin/protocol"
	"bilin/relationlist/controller"
	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/kissgo/httpmetrics"
	"context"
	"strings"
	"time"
)

type RelationListPbObj struct {
}

func NewRelationListPbObj() *RelationListPbObj {

	controller.InitMedalConfig()

	//start timer
	timerHandler := controller.NewTimerManager()
	go timerHandler.Start(1 * time.Second)

	return &RelationListPbObj{}
}

//用户上下麦操作
func (this *RelationListPbObj) RSUserMikeOption(ctx context.Context, req *bilin.RSUserMikeOptionReq) (resp *bilin.RSUserMikeOptionResp, err error) {
	const prefix = "RSUserMikeOption "
	resp = &bilin.RSUserMikeOptionResp{Commonret: bccommon.SUCCESSMESSAGE}

	defer func(now time.Time) {
		httpmetrics.DefReport(strings.TrimSpace(prefix), int64(resp.Commonret.Ret), now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	if req.Opt == bilin.RSUserMikeOptionReq_ONMIKE {
		controller.UserOnMike(req.Owner, req.Header.Userid, time.Now().Unix())
	} else {
		controller.UserOffMike(req.Owner, req.Header.Userid)
	}

	log.Info(prefix, zap.Any("req", req), zap.Any("resp", resp))
	return resp, nil
}

//查询用户勋章
func (this *RelationListPbObj) GetUserRelationMedal(ctx context.Context, req *bilin.GetUserRelationMedalReq) (resp *bilin.GetUserRelationMedalResp, err error) {
	const prefix = "GetUserRelationMedal "
	resp = &bilin.GetUserRelationMedalResp{Commonret: bccommon.SUCCESSMESSAGE}

	defer func(now time.Time) {
		httpmetrics.DefReport(strings.TrimSpace(prefix), int64(resp.Commonret.Ret), now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	resp.Medalid, resp.Medalname, resp.MedalUrl = controller.GetUserMedalInfo(req.Owner, req.Header.Userid)
	log.Info(prefix, zap.Any("req", req), zap.Any("resp", resp))
	return resp, nil
}
