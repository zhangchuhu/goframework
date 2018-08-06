package dao

import (
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"fmt"
	"strconv"
	"strings"
)

type LivingRecord struct {
	ID                 int64
	LivingStartTime    string  `gorm:"column:dt"`
	HOSTID             string  `gorm:"column:uid"`
	LivingTime         int64   `gorm:"column:kaibo_dr"`
	AudienceNum        int64   `gorm:"column:user_cnt"`
	MikeUserNum        int64   `gorm:"column:shangmai"`
	OneMinuteOutInRate float64 `gorm:"column:run_out"`
	RoomID             int64   `gorm:"column:pindao_id"`
	GiftNum            int64   `gorm:"column:liwu"`
	//AverageStayTime
}

func (LivingRecord) TableName() string {
	return "dm_official_anchor_quality_supervision_day"
}

// startTime,endTime format as YYYY-MM-DD
func GetLivingRecordByHostID(hostid []uint64, startTime, endTime string) (livingrecord []LivingRecord, err error) {
	cond := fmt.Sprintf("dt >= \"%s\" and dt <= \"%s\" and uid in ", startTime, endTime)
	var strUids []string
	for _, v := range hostid {
		strUids = append(strUids, strconv.FormatUint(v, 10))
	}
	uidcon := "(" + strings.Join(strUids, ",") + ")"
	cond += uidcon
	fmt.Println(cond)
	if err = BIDataDB.Where(cond).Order("dt").Find(&livingrecord).Error; err != nil {
		appzaplog.Error("Get Guild err", zap.Error(err),
			zap.Uint64s("hostid", hostid),
			zap.String("startTime", startTime),
			zap.String("endTime", endTime),
		)
		return
	}
	appzaplog.Debug("Get Guild", zap.Any("resp", livingrecord))
	return
}

func inCondition(key string, elements []uint64) string {
	var strUids []string
	for _, v := range elements {
		strUids = append(strUids, strconv.FormatUint(v, 10))
	}
	return " and " + key + " in (" + strings.Join(strUids, ",") + ")"
}

func GetLivingRecordByHostIDAndRoomID(hostid []uint64, roomid []uint64, startTime, endTime string) (livingrecord []LivingRecord, err error) {
	if len(hostid) == 0 || len(roomid) == 0 {
		return []LivingRecord{}, nil
	}
	cond := fmt.Sprintf("dt >= \"%s\" and dt <= \"%s\" ", startTime, endTime)
	cond += inCondition("uid", hostid)
	cond += inCondition("pindao_id", roomid)

	fmt.Println(cond)
	if err = BIDataDB.Where(cond).Order("dt").Find(&livingrecord).Error; err != nil {
		appzaplog.Error("Get Guild err", zap.Error(err),
			zap.Uint64s("hostid", hostid),
			zap.String("startTime", startTime),
			zap.String("endTime", endTime),
		)
		return
	}
	appzaplog.Debug("Get Guild", zap.Any("resp", livingrecord))
	return
}
