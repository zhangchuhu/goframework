package handler

import (
	"bilin/guildtars/dao"
	"bilin/protocol"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/kissgo/httpmetrics"
	"context"
	"time"
)

type GuildTarsObj struct {
}

func NewGuildTarsObj() *GuildTarsObj {
	return &GuildTarsObj{}
}

// 根据工会的owuid，查询对应的房间信息
func (this *GuildTarsObj) CategoryGuildRecByOwUid(ctx context.Context, r *bilin.CategoryGuildRecByOwUidReq) (*bilin.CategoryGuildRecByOwUidResp, error) {
	code := int64(0)
	defer func(now time.Time) {
		httpmetrics.DefReport("CategoryGuildRecByOwUid", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	room := dao.GuildRoom{
		GuildID: r.Owuid,
	}
	resp := &bilin.CategoryGuildRecByOwUidResp{}
	info, err := room.Get()
	if err != nil {
		appzaplog.Error("CategoryGuildRecByOwUid Get err", zap.Error(err))
		code = GetGuildRoomFailed
		return resp, err
	}
	roominfos, err := dao.GetGuildRec()
	if err != nil {
		appzaplog.Error("CategoryGuildRecByOwUid Get err", zap.Error(err))
		code = GetGuildRecFailed
		return resp, err
	}
	for _, v := range info {
		for _, roominfo := range roominfos {
			if uint64(v.RoomID) == roominfo.RoomID {
				resp.Info = append(resp.Info, &bilin.CategoryGuildRecInfo{
					Roomid: roominfo.RoomID,
					Typeid: roominfo.TypeId,
				})
				break
			}
		}
	}
	return resp, nil
}
