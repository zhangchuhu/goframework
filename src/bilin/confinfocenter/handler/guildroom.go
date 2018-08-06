package handler

import (
	"bilin/confinfocenter/dao"
	"bilin/protocol"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/kissgo/httpmetrics"
	"context"
	"errors"
	"time"
)

// 工会频道关系
func (this *ConfInfoServantObj) GuildRoomS(ctx context.Context, r *bilin.GuildRoomSReq) (*bilin.GuildRoomSResp, error) {
	code := int64(0)
	defer func(now time.Time) {
		httpmetrics.DefReport("GuildRoomS", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	resp := &bilin.GuildRoomSResp{}
	room := dao.GuildRoom{}
	if r.Info != nil {
		room.RoomID = r.Info.Roomid
		room.GuildID = r.Info.Guildid
		room.ID = uint(r.Info.Id)
	}
	rec, err := room.Get()
	if err != nil {
		code = GetGuildRoomFailed
		return resp, err
	}

	for _, v := range rec {
		resp.Info = append(resp.Info, &bilin.GuildRoomInfo{
			Roomid:  v.RoomID,
			Guildid: v.GuildID,
			Id:      int64(v.ID),
		})
	}
	return resp, nil
}

func (this *ConfInfoServantObj) DelGuildRoom(ctx context.Context, r *bilin.DelGuildRoomReq) (*bilin.DelGuildRoomResp, error) {
	code := int64(0)
	defer func(now time.Time) {
		httpmetrics.DefReport("DelGuildRoom", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	resp := &bilin.DelGuildRoomResp{}

	if r.Info.Id == 0 {
		appzaplog.Warn("DelCategoryHostRec all", zap.Any("req", r))
		return resp, errors.New("DelCategoryHostRec all")
	}

	rec := dao.GuildRoom{}
	rec.ID = uint(r.Info.Id)
	if err := rec.Del(); err != nil {
		code = DelGuildRoomFailed
		return resp, err
	}
	return resp, nil
}

func (this *ConfInfoServantObj) CreateGuildRoom(ctx context.Context, r *bilin.CreateGuildRoomReq) (*bilin.CreateGuildRoomResp, error) {
	code := int64(0)
	defer func(now time.Time) {
		httpmetrics.DefReport("CreateGuildRoom", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	resp := &bilin.CreateGuildRoomResp{}
	rec := dao.GuildRoom{
		GuildID: r.Info.Guildid,
		RoomID:  r.Info.Roomid,
	}
	if err := rec.Create(); err != nil {
		code = CreateGuildRoomFailed
		return resp, err
	}
	return resp, nil
}
