/*
 * Copyright (c) 2018-07-03.
 * Author: kordenlu
 */

package dao_test

import (
	"bilin/chattag/config"
	"bilin/chattag/dao"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"reflect"
	"testing"
)

func init() {
	if err := dao.InitRedisDao(&config.AppConfig{
		SentinelAddrs: []string{"grouptest1006-bj-sentinel.yy.com:20006"},
	}); err != nil {
		appzaplog.Error("InitRedisDao err", zap.Error(err))
	}
}

func TestAddUsedUserTopic(t *testing.T) {
	var (
		fromUid  = int64(100)
		toUid    = int64(200)
		topic    = "truthtopic"
		topicIds = []int64{10}
	)
	err := dao.AddUsedUserTopic(fromUid, toUid, topic, topicIds)
	if err != nil {
		t.Error(err)
	}
	info, err := dao.GetUsedUserTopic(fromUid, toUid, topic)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(info, topicIds) {
		t.Errorf("%v not equesl %v", info, topicIds)
	}
}
