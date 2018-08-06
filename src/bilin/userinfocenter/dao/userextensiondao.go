package dao

import (
	"fmt"
	"strconv"
)

type USEREXTENSION struct {
	FK_USER_ID uint64 `gorm:"column:FK_USER_ID"`
	CITY       int64  `gorm:"column:CITY"`
}

func GetUserExtension(uid uint64) (*USEREXTENSION, error) {
	index := uid % 100 //分表
	table_name := fmt.Sprintf("USER_EXTENSION_%d", index)

	var ret USEREXTENSION
	condition := "FK_USER_ID = " + strconv.FormatUint(uid, 10)
	db_ := UserDB.Table(table_name).First(&ret, condition)
	if db_.RecordNotFound() {
		return nil, nil
	}
	return &ret, db_.Error
}
