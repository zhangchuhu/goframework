package handler

import (
	"bilin/bcserver/bccommon"
	"bilin/bcserver/domain/collector"
	"bilin/bcserver/domain/entity"
	"bilin/bcserver/domain/service"
	"bilin/protocol"
	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/kissgo/httpmetrics"
	"context"
	"fmt"
	"strings"
	"time"
)

// 主持人锁定/解锁房间
func (this *BCServantObj) LockUnlockRoomOperation(ctx context.Context, req *bilin.LockUnlockRoomOperationReq) (resp *bilin.LockUnlockRoomOperationResp, err error) {
	const prefix = "LockUnlockRoomOperation "
	roomid := req.Header.Roomid
	userid := req.Header.Userid
	resp = &bilin.LockUnlockRoomOperationResp{Commonret: bccommon.SUCCESSMESSAGE}
	log.Info(prefix+"begin", zap.Any("req", req))

	defer func(now time.Time) {
		httpmetrics.DefReport(strings.TrimSpace(prefix), int64(resp.Commonret.Ret), now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	var room *entity.Room
	var user *entity.User
	if room, user, err = this.CommonCheckAuth(roomid, userid); err != nil {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_ILLEGAL_MESSAGE, bccommon.COMMERRORILLEGALDESC)
		log.Warn(prefix+" failed", zap.Any("req", req), zap.Any("resp", resp))
		return resp, nil
	}

	if room.RoomType2 == service.OFFICAIL_ROOM {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_ILLEGAL_MESSAGE, fmt.Sprintf("官频不能做此操作"))
		log.Warn(prefix+"failed, official room not allowed", zap.Any("req", req))
		return
	}

	//检查用户权限
	if userid != room.Owner {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_LOCKUNLOCK_NO_RIGHT, fmt.Sprintf("用户没有权限"))
		log.Warn(prefix+"failed, permission denied", zap.Any("req", req), zap.Any("Role", user.Role))
		return
	}

	var bizErr error
	if req.Opt == entity.UNLOCKROOM {
		bizErr = collector.BizRemoveRoomPassword(room)
	} else {
		if len(req.Pwd) < 4 {
			resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_LOCKUNLOCK_FAILED, fmt.Sprintf("密码格式不对"))
			log.Warn(prefix+"failed, error room password", zap.Any("req", req))
			return
		}
		bizErr = collector.BizSetRoomPassword(room, req.Pwd)
	}

	if bizErr != nil {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_LOCKUNLOCK_FAILED, fmt.Sprintf("服务器开小差了，再试试呗~"))
		log.Error(prefix+"failed, lock unlock error", zap.Any("req", req), zap.Any("err", err))
		return
	}

	log.Info("[+]"+prefix+"success", zap.Any("req", req), zap.Any("resp", resp))
	return
}
