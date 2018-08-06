package dao

import (
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"strconv"
)

type USER_BL_ID struct {
	USER_ID  uint64 `gorm:"column:USER_ID"`
	BILIN_ID uint64 `gorm:"column:BILIN_ID"`
}

func (USER_BL_ID) TableName() string {
	return "USER_BL_ID"
}

func GetUserBiID(uid uint64) (*USER_BL_ID, error) {
	ret := &USER_BL_ID{}
	cond := "USER_ID = " + strconv.FormatUint(uid, 10)
	if err := UserDB.First(ret, cond).Error; err != nil {
		appzaplog.Error("GetUserBiID err", zap.Error(err), zap.Uint64("uid", uid))
		return ret, err
	}
	return ret, nil
}

func BatchUserBiID(uids []uint64) ([]USER_BL_ID, error) {
	ret := []USER_BL_ID{}
	if len(uids) == 0 {
		return ret, nil
	}
	cond := "USER_ID IN ( "
	for k, v := range uids {
		cond += strconv.FormatUint(v, 10)
		if k+1 < len(uids) {
			cond += ","
		}
	}
	cond += " )"
	if err := UserDB.Find(&ret, cond).Error; err != nil {
		appzaplog.Error("GetUserBiID err", zap.Error(err), zap.Uint64s("uids", uids))
		return ret, err
	}
	return ret, nil
}

func BatchUserID(bilinids []uint64) ([]USER_BL_ID, error) {
	ret := []USER_BL_ID{}
	if len(bilinids) == 0 {
		return ret, nil
	}
	cond := "BILIN_ID IN ( "
	for k, v := range bilinids {
		cond += strconv.FormatUint(v, 10)
		if k+1 < len(bilinids) {
			cond += ","
		}
	}
	cond += " )"
	if err := UserDB.Find(&ret, cond).Error; err != nil {
		appzaplog.Error("BatchUserID err", zap.Error(err), zap.Uint64s("bilinids", bilinids))
		return ret, err
	}
	return ret, nil
}
