package handler

import (
	"bilin/protocol"
	"bilin/roominfocenter/cache"
	"code.yy.com/yytars/goframework/kissgo/httpmetrics"
	"context"
	"time"
)

type RoomInfoCenterServantObj struct {
}

func NewRoomInfoCenterServantObj() *RoomInfoCenterServantObj {
	return &RoomInfoCenterServantObj{}
}

// 所有正在开播的直播间信息
func (p *RoomInfoCenterServantObj) LivingRoomsInfo(ctx context.Context, r *bilin.LivingRoomsInfoReq) (*bilin.LivingRoomsInfoResp, error) {
	defer func(now time.Time) {
		httpmetrics.DefReport("LivingRoomsInfo", 0, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	var (
		resp = &bilin.LivingRoomsInfoResp{
			Livingrooms: cache.GetRoomCache(),
		}
	)
	return resp, nil
}

func (p *RoomInfoCenterServantObj) BatchLivingRoomsInfoByHosts(ctx context.Context, r *bilin.BatchLivingRoomsInfoByHostsReq) (*bilin.BatchLivingRoomsInfoByHostsResp, error) {
	defer func(now time.Time) {
		httpmetrics.DefReport("BatchLivingRoomsInfoByHosts", 0, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	Livingrooms := cache.GetRoomCache()
	resp := &bilin.BatchLivingRoomsInfoByHostsResp{
		Livingrooms: make(map[uint64]*bilin.RoomInfo),
	}
	for _, host := range r.Hosts {
		for _, info := range Livingrooms {
			if host == info.Owner {
				resp.Livingrooms[info.Owner] = info
				break
			}
		}
	}
	return resp, nil
}

// 查询指定房间是否在开播
func (p *RoomInfoCenterServantObj) IsLiving(ctx context.Context, r *bilin.IsLivingReq) (*bilin.IsLivingResp, error) {
	defer func(now time.Time) {
		httpmetrics.DefReport("IsLiving", 0, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	Livingrooms := cache.GetRoomCache()
	resp := &bilin.IsLivingResp{}
	_, ok := Livingrooms[uint64(r.Roomid)]
	resp.Isliving = ok
	return resp, nil
}
