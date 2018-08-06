package entity

const (
	UNLOCKROOM = 0
	LOCKROOM   = 1
)

type BizRoomInfo struct {
	RoomId     uint64 `json:"roomid"`
	LockStatus int32  `json:"lockstatus"`
}
