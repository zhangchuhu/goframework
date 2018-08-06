package handler

import (
	"bilin/bcserver/bccommon"
	"bilin/bizroomcenter/service"
	"bilin/protocol"
	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/kissgo/httpmetrics"
	"context"
	"fmt"
	"strings"
	"time"
)

const (
	UNLOCKROOM = 0
	LOCKROOM   = 1
)

var _ bilin.BizRoomCenterServantServer = &BizRoomCenterPbObj{}

type BizRoomCenterPbObj struct {
}

func NewBizRoomCenterPbObj() (o *BizRoomCenterPbObj) {
	service.MysqlInit()

	return &BizRoomCenterPbObj{}
}

//获取房间的基本业务信息
func (this *BizRoomCenterPbObj) GetBizRoomInfo(ctx context.Context, req *bilin.GetBizRoomInfoReq) (resp *bilin.GetBizRoomInfoResp, err error) {
	const prefix = "BatchGetBizRoomInfo "
	resp = &bilin.GetBizRoomInfoResp{Commonret: bccommon.SUCCESSMESSAGE}

	defer func(now time.Time) {
		httpmetrics.DefReport(strings.TrimSpace(prefix), int64(resp.Commonret.Ret), now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	resp.Bizroominfo, _ = service.MysqlGetBizRoomInfo(req.Roomid)

	log.Info(prefix, zap.Any("req", req), zap.Any("resp", resp))
	return resp, nil
}

func (this *BizRoomCenterPbObj) BatchGetBizRoomInfo(ctx context.Context, req *bilin.BatchGetBizRoomInfoReq) (resp *bilin.BatchGetBizRoomInfoResp, err error) {
	const prefix = "BatchGetBizRoomInfo "
	resp = &bilin.BatchGetBizRoomInfoResp{Commonret: bccommon.SUCCESSMESSAGE}

	defer func(now time.Time) {
		httpmetrics.DefReport(strings.TrimSpace(prefix), int64(resp.Commonret.Ret), now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	log.Info(prefix, zap.Any("req", req), zap.Any("resp", resp))
	return resp, nil
}

//设置房间密码
func (this *BizRoomCenterPbObj) SetRoomPassword(ctx context.Context, req *bilin.SetRoomPasswordReq) (resp *bilin.SetRoomPasswordResp, err error) {
	const prefix = "SetRoomPassword "
	resp = &bilin.SetRoomPasswordResp{Commonret: bccommon.SUCCESSMESSAGE}

	defer func(now time.Time) {
		httpmetrics.DefReport(strings.TrimSpace(prefix), int64(resp.Commonret.Ret), now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	mysqlErr := service.MysqlSetBizRoomInfo(&bilin.BizRoomInfo{Roomid: req.Roomid, Lockstatus: LOCKROOM, Roompwd: req.Password})
	if mysqlErr != nil {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_LOCKUNLOCK_FAILED, fmt.Sprintf("服务器异常"))
		log.Error(prefix, zap.Any("req", req), zap.Any("resp", resp))
		return
	}

	log.Info(prefix, zap.Any("req", req), zap.Any("resp", resp))
	return resp, nil
}

func (this *BizRoomCenterPbObj) RemoveRoomPassword(ctx context.Context, req *bilin.RemoveRoomPasswordReq) (resp *bilin.RemoveRoomPasswordResp, err error) {
	const prefix = "RemoveRoomPassword "
	resp = &bilin.RemoveRoomPasswordResp{Commonret: bccommon.SUCCESSMESSAGE}

	defer func(now time.Time) {
		httpmetrics.DefReport(strings.TrimSpace(prefix), int64(resp.Commonret.Ret), now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	mysqlErr := service.MysqlDelBizRoomInfo(req.Roomid)
	if mysqlErr != nil {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_LOCKUNLOCK_FAILED, fmt.Sprintf("服务器异常"))
		log.Error(prefix, zap.Any("req", req), zap.Any("resp", resp))
		return
	}

	log.Info(prefix, zap.Any("req", req), zap.Any("resp", resp))
	return resp, nil
}

//获取房间锁定状态
func (this *BizRoomCenterPbObj) GetBizRoomLockStatus(ctx context.Context, req *bilin.GetBizRoomLockStatusReq) (resp *bilin.GetBizRoomLockStatusResp, err error) {
	const prefix = "GetBizRoomLockStatus "
	resp = &bilin.GetBizRoomLockStatusResp{Commonret: bccommon.SUCCESSMESSAGE, Lockstatus: UNLOCKROOM}

	defer func(now time.Time) {
		httpmetrics.DefReport(strings.TrimSpace(prefix), int64(resp.Commonret.Ret), now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	data, _ := service.MysqlGetBizRoomInfo(req.Roomid)
	if data != nil {
		resp.Lockstatus = LOCKROOM
	}

	log.Info(prefix, zap.Any("req", req), zap.Any("resp", resp))
	return resp, nil
}

func (this *BizRoomCenterPbObj) BatchGetBizRoomLockStatus(ctx context.Context, req *bilin.BatchGetBizRoomLockStatusReq) (resp *bilin.BatchGetBizRoomLockStatusResp, err error) {
	const prefix = "BatchGetBizRoomLockStatus "
	resp = &bilin.BatchGetBizRoomLockStatusResp{Commonret: bccommon.SUCCESSMESSAGE}

	defer func(now time.Time) {
		httpmetrics.DefReport(strings.TrimSpace(prefix), int64(resp.Commonret.Ret), now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	infos, _ := service.MysqlGetAllBizRoomInfos()
	for _, roomid := range req.Roomids {
		for _, item := range infos {
			if item.Roomid == roomid {
				resp.Roomids = append(resp.Roomids, roomid)
			}
		}
	}

	log.Info(prefix, zap.Any("req", req), zap.Any("resp", resp))
	return resp, nil
}

//获取所有锁定的房间列表
func (this *BizRoomCenterPbObj) GetAllLockedRooms(ctx context.Context, req *bilin.GetAllLockedRoomsReq) (resp *bilin.GetAllLockedRoomsResp, err error) {
	const prefix = "GetAllLockedRooms "
	resp = &bilin.GetAllLockedRoomsResp{Commonret: bccommon.SUCCESSMESSAGE}

	defer func(now time.Time) {
		httpmetrics.DefReport(strings.TrimSpace(prefix), int64(resp.Commonret.Ret), now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	infos, _ := service.MysqlGetAllBizRoomInfos()
	for _, item := range infos {
		resp.Roomids = append(resp.Roomids, item.Roomid)
	}

	log.Info(prefix, zap.Any("req", req), zap.Any("resp", resp))
	return resp, nil
}
