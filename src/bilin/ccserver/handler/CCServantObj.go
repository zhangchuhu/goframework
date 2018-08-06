package handler

import (
	"bilin/common/onlinepush"
	"bilin/common/onlinequery"
	"bilin/protocol"
	"context"
	"fmt"
	"math/rand"
	"time"

	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"github.com/golang/protobuf/proto"
)

const (
	minRoomSeq = 1000000000
	maxRoomSeq = 1999999999
)

var _ bilin.CCServantServer = &ccServantObj{}

type ccServantObj struct {
	roomSeq int64
}

func NewCCServantObj() *ccServantObj {
	rand.Seed(time.Now().Unix())
	return &ccServantObj{
		roomSeq: minRoomSeq + (time.Now().UnixNano()/1000000)&0xfffffff, // 9-digits decimal number
	}
}

func CalcOnlineDisplay(cnt int) (display int) {

	switch {
	case cnt <= 2315:
		display = 2315
	case cnt > 2315 && cnt <= 4901:
		display = cnt
	case cnt > 4901:
		display = int(float32(cnt)*3.8 - 11764)
	}
	if display > 22436 {
		display = 22436 + cnt%100
	}
	return
}

func (this *ccServantObj) GetRandomCallNumberClient(ctx context.Context, req *bilin.GetRandomCallNumberClientReq) (rsp *bilin.GetRandomCallNumberClientResp, err error) {
	const prefix = "GetRandomCallNumberClient "
	rsp = &bilin.GetRandomCallNumberClientResp{
		Cret:           retSuccess(),
		NumberOfClient: 6235,
	}
	var cnt int
	if cnt, err = onlinequery.UserCount(); err != nil {
		rsp.NumberOfClient = 6235 // fall back count when service unavailable
		log.Warn(prefix+"fail", zap.Any("err", err), zap.Any("rsp", rsp))
		return
	}
	rsp.NumberOfClient = int64(CalcOnlineDisplay(cnt))
	log.Debug(prefix+"ok", zap.Any("cnt", cnt), zap.Any("rsp", rsp))
	return
}

func (this *ccServantObj) GenerateRoom(ctx context.Context, req *bilin.GenerateRoomReq) (rsp *bilin.GenerateRoomResp, err error) {
	const prefix = "GenerateRoom "
	rsp = &bilin.GenerateRoomResp{
		Cret: retSuccess(),
	}
	rsp.RoomID = this.roomSeq
	this.roomSeq++
	if this.roomSeq > maxRoomSeq {
		this.roomSeq = minRoomSeq
	}
	log.Debug(prefix+"ok", zap.Any("rsp", rsp))
	return
}

func (this *ccServantObj) SendMessageToUser(ctx context.Context, req *bilin.SendMessageToUserReq) (rsp *bilin.SendMessageToUserResp, err error) {
	const prefix = "SendMessageToUser "
	rsp = &bilin.SendMessageToUserResp{
		Cret: retSuccess(),
	}
	var (
		uid1       int64
		uid2       []int64
		rid1       int64
		roomUsers  map[int64]int32
		passUsers  []int64
		dropUsers  []int64
		mpush      bilin.MultiPush
		cbody      bilin.CommonMessageBody
		cbodyBytes []byte
	)
	if req.Header == nil {
		rsp.Cret = retError("request header is nil")
		log.Error(prefix+"fail", zap.Any("rsp", rsp))
		return
	}
	uid1 = int64(req.Header.Userid)
	uid2 = req.ToUserID
	if uid1 <= 0 || len(uid2) == 0 {
		rsp.Cret = retError(fmt.Sprintf("invalid uid: from %v, to %v", uid1, uid2))
		log.Warn(prefix+"fail", zap.Any("rsp", rsp))
		return
	}
	if rid1, err = onlinequery.GetUserRoom(uid1); err != nil {
		rsp.Cret = retError(fmt.Sprintf("query user room failure: uid %v", uid1))
		log.Warn(prefix+"fail", zap.Any("err", err), zap.Any("rsp", rsp))
		return
	}
	if rid1 < 0 {
		rsp.Cret = retError(fmt.Sprintf("current user %v is not in a room (%v)", uid1, rid1))
		log.Warn(prefix+"fail", zap.Any("rsp", rsp))
		return
	}
	if roomUsers, err = onlinequery.GetRoomUser(rid1); err != nil {
		rsp.Cret = retError(fmt.Sprintf("query room %v user failure", rid1))
		log.Warn(prefix+"fail", zap.Any("err", err), zap.Any("rsp", rsp))
		return
	}
	for _, id := range uid2 {
		if _, ok := roomUsers[id]; ok {
			passUsers = append(passUsers, id)
		} else {
			dropUsers = append(dropUsers, id)
		}
	}
	cbody.Type = int32(bilin.MinType_CC_CLIENT_P2P_TUNNEL)
	cbody.Data = req.Data
	if cbodyBytes, err = proto.Marshal(&cbody); err != nil {
		rsp.Cret = retError("protobuf marshal failure")
		log.Warn(prefix+"fail", zap.Any("err", err), zap.Any("rsp", rsp))
		return
	}
	mpush.UserIDs = passUsers
	mpush.Msg = &bilin.ServerPush{
		MessageType: int32(bilin.MaxType_CC_MSG),
		PushBuffer:  cbodyBytes,
	}
	if _, err = onlinepush.PushToUser(mpush); err != nil {
		rsp.Cret = retError("push service failure")
		log.Warn(prefix+"fail", zap.Any("err", err), zap.Any("rsp", rsp))
		return
	}
	rsp.Cret.Desc = fmt.Sprintf("pass %v, drop %v", passUsers, dropUsers)
	log.Debug(prefix+"ok", zap.Any("rsp", rsp))
	return
}

func (this *ccServantObj) GetUserCurrentRoom(ctx context.Context, req *bilin.GetUserCurrentRoomReq) (rsp *bilin.GetUserCurrentRoomResp, err error) {
	const prefix = "GetUserCurrentRoom "
	rsp = &bilin.GetUserCurrentRoomResp{
		Cret: retSuccess(),
	}
	if req.Header == nil {
		rsp.Cret = retError("request header is nil")
		log.Error(prefix+"fail", zap.Any("rsp", rsp))
		return
	}
	var uid int64 = int64(req.Header.Userid)
	var rid int64
	if rid, err = onlinequery.GetUserRoom(uid); err != nil {
		rsp.Cret = retError(fmt.Sprintf("query user room failure: uid %v", uid))
		log.Warn(prefix+"fail", zap.Any("err", err), zap.Any("rsp", rsp))
		return
	}
	if rid < 0 {
		rsp.Cret = retError(fmt.Sprintf("current user %v is not in a room (%v)", uid, rid))
		log.Warn(prefix+"fail", zap.Any("rsp", rsp))
		return
	}
	rsp.RoomID = rid
	log.Debug(prefix+"ok", zap.Any("rsp", rsp))
	return
}
