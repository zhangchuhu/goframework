package dao

import (
	//log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	//"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"errors"
	"strings"
)

var NoAvailabelDB = errors.New("no database connection available")

func IsTableNotExist(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "Table") && strings.Contains(err.Error(), "doesn't exist")
}
