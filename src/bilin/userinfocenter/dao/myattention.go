//我的关注
package dao

import (
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"errors"
	"fmt"
)

type MY_ATTENTION_USER struct {
	FROM_USER_ID uint64 `gorm:"column:FROM_USER_ID"`
	TO_USER_ID   uint64 `gorm:"column:TO_USER_ID"`
	ENABLED      int64  `gorm:"column:ENABLED"`
	UPDATE_TIME  int64  `gorm:"column:UPDATE_TIME"`
	CREATE_ON    int64  `gorm:"column:CREATE_ON"`
}

func MyAttentionNum(uid uint64) (count uint64, err error) {
	if AttentionDB == nil {
		appzaplog.Warn("no attention database connection available", zap.Uint64("uid", uid))
		err = errors.New("no attention database connection available")
		return
	}

	index := uid % 100 //分表
	table_name := fmt.Sprintf("MY_ATTENTION_USER_%d", index)
	cond := fmt.Sprintf("FROM_USER_ID = %d and ENABLED = %d", uid, 1)
	if err = AttentionDB.Table(table_name).Where(cond).Count(&count).Error; err != nil {
		appzaplog.Error("query db err", zap.Error(err), zap.Uint64("uid", uid))
		return
	}
	return
}
