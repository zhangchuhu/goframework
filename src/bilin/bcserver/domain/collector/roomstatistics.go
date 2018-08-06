package collector

import (
	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"fmt"
	"os"
	"time"
)

const (
	RoomStatisticsFileName = "broadcast_meeting_stat"
)

var (
	StatPath string
)

func RoomStatisticsInit(appname string) (err error) {
	StatPath = "/data/bilin/" + appname + "/stat/logs"

	if err = os.MkdirAll(StatPath, 0755); err != nil {
		log.Error("MkdirAll", zap.Error(err), zap.String("path", StatPath))
	}

	return
}

func writeHourStat(message string) (err error) {
	const prefix = "writeHourStat "
	t := time.Now()
	file_suffix := fmt.Sprintf("%d-%02d-%02d-%02d", t.Year(), t.Month(), t.Day(), t.Hour())
	fileName := StatPath + "/" + RoomStatisticsFileName + "_" + file_suffix
	fd, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Error(prefix+"OpenFile", zap.Error(err), zap.String("filename", fileName), zap.String("message", message))
		return
	}
	defer fd.Close()
	if _, err := fd.Write([]byte(message + "\n")); err != nil {
		log.Error(prefix+"Write", zap.Error(err), zap.String("filename", fileName), zap.String("message", message))
	}
	log.Debug(prefix + "end")
	return
}

//第一个进入房间，表示房间被创建
func CreateRoomStat(roomid uint64, uid uint64) (err error) {
	const prefix = "CreateRoomStat "

	currentTime := time.Now()
	fmtTime := fmt.Sprintf("%d%02d%02d%02d%02d%02d", currentTime.Year(), currentTime.Month(), currentTime.Day(), currentTime.Hour(), currentTime.Minute(), currentTime.Second())
	writeHourStat(fmt.Sprintf("%s,LIVE_BROADCAST={uid=%d&rid=%d&step=CB}", fmtTime, uid, roomid))

	log.Debug(prefix + "end")
	return
}

//用户进入频道
func EnterRoomStat(roomid uint64, uid uint64, role uint32) (err error) {
	const prefix = "EnterRoomStat "

	currentTime := time.Now()
	fmtTime := fmt.Sprintf("%d%02d%02d%02d%02d%02d", currentTime.Year(), currentTime.Month(), currentTime.Day(), currentTime.Hour(), currentTime.Minute(), currentTime.Second())
	writeHourStat(fmt.Sprintf("%s,LIVE_BROADCAST_USER={rid=%d&uid=%d&step=JLB&identify=%d}", fmtTime, roomid, uid, role))

	log.Debug(prefix + "end")
	return
}

//用户离开频道
func ExitRoomStat(roomid uint64, uid uint64, role uint32, beginTime int64, endTime int64, join_type int) (err error) {
	const prefix = "ExitRoomStat "

	currentTime := time.Now()
	fmtTime := fmt.Sprintf("%d%02d%02d%02d%02d%02d", currentTime.Year(), currentTime.Month(), currentTime.Day(), currentTime.Hour(), currentTime.Minute(), currentTime.Second())
	writeHourStat(fmt.Sprintf("%s,LIVE_BROADCAST_USER={rid=%d&uid=%d&step=ELB&identify=%d&begin_time=%d&end_time=%d&join_type=%d}",
		fmtTime, roomid, uid, role, beginTime, endTime, join_type))

	log.Debug(prefix + "end")
	return
}
