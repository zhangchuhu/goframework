package entity

type LivingRecordInfo struct {
	Roomid     int64  `json:"roomid"`
	Owner      string `json:"owner"`
	Title      string `json:"title"`
	RoomType2  string `json:"roomType2"`
	Starttime  string `json:"starttime"`
	Endtime    string `json:"endtime"`
	LivingTime string `json:"livingTime"`
}

type AllLivingRecordInfoResp struct {
	Result    int32               `json:"result"`
	ErrorDesc string              `json:"errordesc"`
	Data      []*LivingRecordInfo `json:"items"`
}
