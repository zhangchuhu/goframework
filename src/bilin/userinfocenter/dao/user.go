package dao

import (
	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"strconv"
	"time"
)

type UserInfo struct {
	UserId             uint64    `gorm:"primary_key;column:ID"`
	AvatarId           uint64    `gorm:"not null;column:BIG_HEAD_IMG_ID"`
	NickName           string    `gorm:"column:NICK_NAME"`
	Sex                uint32    `gorm:"column:SEX"`
	Sign               string    `gorm:"column:SIGN"`
	City               string    `gorm:"column:CITY"`
	Birthday           time.Time `gorm:"column:BIRTHDAY" sql:"DEFAULT:nil"`
	SHOW_GLAMOUR_VALUE uint64    `gorm:"column:SHOW_GLAMOUR_VALUE"`
	ShowSex            uint32    `gorm:"column:SHOW_SEX"`
}

func GetUserInfo(uid uint64) (*UserInfo, error) {

	if UserDB == nil {
		log.Warn("no user database connection available")
		return nil, NoAvailabelDB
	}

	user := UserInfo{}
	index := uid % 100 //分表
	table_name := fmt.Sprintf("USER_%d", index)
	condition := "ID = " + strconv.FormatUint(uid, 10)
	db_ := UserDB.Table(table_name).First(&user, condition)
	if db_.RecordNotFound() {
		return nil, nil
	}

	if db_.Error != nil {
		log.Error("GetUserInfo fail", zap.Uint64("uid", uid), zap.Error(db_.Error))
		return nil, db_.Error
	}

	return &user, nil
}

func GetAvatatrUsers(index uint64, count uint64) ([]UserInfo, error) {

	if UserDB == nil {
		log.Warn("no user database connection available")
		return nil, NoAvailabelDB
	}

	if count > 60 {
		count = 60
	}
	users := make([]UserInfo, 0, count)
	index = index % 100
	table_name := fmt.Sprintf("USER_%d", index)
	db_ := UserDB.Table(table_name).Limit(count).Where("BIG_HEAD_IMG_ID > 0").Find(&users)
	//db_ := UserDB.Table(table_name).Limit(count).Find(&users)
	if db_.RecordNotFound() {
		return nil, nil
	}

	if db_.Error != nil {
		if IsTableNotExist(db_.Error) {
			log.Warn("query mysql fail", zap.Error(db_.Error))
			return nil, nil //????, ??????, ?????????
		}
		log.Error("GetUserInfo fail", zap.Uint64("index", index), zap.Error(db_.Error))
		return nil, db_.Error
	}

	return users, nil
}
