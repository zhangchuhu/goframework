// @author kordenlu
// @创建时间 2018/03/29 17:13
// 功能描述:

package entity

import (
	"bilin/protocol"
	"time"
)

type Room struct {
	UniqueId           int64                                  `json:"uniqueid"`
	Roomid             uint64                                 `json:"roomid"`
	Owner              uint64                                 `json:"owner"`
	Status             bilin.BaseRoomInfo_ROOMSTATUS          `json:"status"`
	RoomType           bilin.BaseRoomInfo_ROOMTYPE            `json:"roomtype"`
	LinkStatus         bilin.BaseRoomInfo_LINKSTATUS          `json:"linkstatus"`
	Title              string                                 `json:"title"`
	RoomType2          int32                                  `json:"roomType2"`
	RoomCategoryID     int32                                  `json:"roomCategoryID"`
	RoomPendantLevel   int32                                  `json:"roomPendantLevel"`
	HostBilinID        int64                                  `json:"hostBilinID"`
	StartTime          uint64                                 `json:"starttime"`
	EndTime            uint64                                 `json:"endtime"`
	AutoLink           bilin.BaseRoomInfo_AUTOLINK            `json:"autolink"`
	Maixuswitch        bilin.BaseRoomInfo_MAIXUSWITCH         `json:"maixuswitch"`
	From               string                                 `json:"from"`
	Karaokeswitch      bilin.BaseRoomInfo_KARAOKESWITCH       `json:"karaoke_switch"`
	Relationlistswitch bilin.BaseRoomInfo_RELATIONLISTESWITCH `json:"relationlist_switch"`
	LockStatus         uint32                                 `json:"lock_status"`
}

func NewRoom(roomid uint64) *Room {
	return &Room{
		Roomid:             roomid,
		Owner:              0,
		Status:             bilin.BaseRoomInfo_OPEN,
		RoomType:           bilin.BaseRoomInfo_ROOMTYPE_THREE,
		LinkStatus:         bilin.BaseRoomInfo_CLOSELINK,
		Title:              "",
		StartTime:          uint64(time.Now().Unix()),
		AutoLink:           bilin.BaseRoomInfo_CLOSEAUTOTOMIKE,
		Maixuswitch:        bilin.BaseRoomInfo_CLOSEMAIXU, // 默认关闭开关，等运营通知打开的时候需要修改这里默认值为开  todo
		Karaokeswitch:      bilin.BaseRoomInfo_CLOSEKARAOKE,
		Relationlistswitch: bilin.BaseRoomInfo_OPENRELATIONLIST,
	}
}

func (room *Room) SetLinkstatus(newStatus bilin.BaseRoomInfo_LINKSTATUS) {
	room.LinkStatus = newStatus
}

func (room *Room) GetLinkStatus() bilin.BaseRoomInfo_LINKSTATUS {
	return room.LinkStatus
}

func (room *Room) SetAutoLink(newStatus bilin.BaseRoomInfo_AUTOLINK) {
	room.AutoLink = newStatus
}

func (room *Room) GetAutoLink() bilin.BaseRoomInfo_AUTOLINK {
	if room.RoomType == bilin.BaseRoomInfo_ROOMTYPE_RADIO { //电台模板不能自动连麦
		return bilin.BaseRoomInfo_CLOSEAUTOTOMIKE
	}

	return room.AutoLink
}
