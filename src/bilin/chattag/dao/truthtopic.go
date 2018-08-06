package dao

import (
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"github.com/jinzhu/gorm"
)

type TruthTopic struct {
	gorm.Model
	Topic string
}

//todo move to tar service
func GetTruthTopicNotIn(ids []int64, limit int) ([]TruthTopic, error) {
	var ret []TruthTopic
	if err := hujiaoChatTagDB.Not(ids).Limit(limit).Find(&ret).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) || IsTableNotExistErr(err) {
			return nil, nil
		}
		appzaplog.Error("GetTopicNotIn err", zap.Error(err), zap.Int64s("ids", ids))
		return nil, err
	}
	return ret, nil
}

//todo 当表变大时，需要更改这种random方式
func RandTruthTopic(limit int64) ([]TruthTopic, error) {
	var ret []TruthTopic
	if err := hujiaoChatTagDB.Limit(limit).Order(gorm.Expr("rand()")).Find(&ret).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) || IsTableNotExistErr(err) {
			return nil, nil
		}
		appzaplog.Error("RandTruthTopic err", zap.Error(err))
		return nil, err
	}
	return ret, nil
}
