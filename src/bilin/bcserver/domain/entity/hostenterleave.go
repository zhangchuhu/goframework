package entity

//主持人进房间，从队列中删除，主持人退出失败时，需要加入到队列中
type HostEnterLeaveTask struct {
	RoomId    uint64 `json:"roomid"`
	RoomType  int32  `json:"roomtype"`
	HostId    uint64 `json:"hostid"`
	LeaveTime uint64 `json:"leave_time"`
}
