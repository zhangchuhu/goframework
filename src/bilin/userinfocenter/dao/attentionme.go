/*
我的粉丝
*/
package dao

import (
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"errors"
	"fmt"
)

type ATTENTION_ME_USER struct {
}

func AttentionMeNum(uid uint64) (count uint64, err error) {
	if AttentionDB == nil {
		appzaplog.Warn("no attention database connection available", zap.Uint64("uid", uid))
		err = errors.New("no attention database connection available")
		return
	}

	index := uid % 100 //分表
	table_name := fmt.Sprintf("ATTENTION_ME_USER_%d", index)
	cond := fmt.Sprintf("FROM_USER_ID = %d and ENABLED = %d", uid, 1)
	if err = AttentionDB.Table(table_name).Where(cond).Count(&count).Error; err != nil {
		appzaplog.Error("query db err", zap.Error(err), zap.Uint64("uid", uid))
		return
	}
	return
}
