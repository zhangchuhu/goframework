/*
 * Copyright (c) 2018-07-03.
 * Author: kordenlu
 * 功能描述: 用户使用过的pua和真心话话题
 */

package dao

import (
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"fmt"
	"strconv"
	"time"
)

func usedUserTopicKey(fromUserId, toUserId int64, topicName string) string {
	if fromUserId > toUserId {
		return fmt.Sprintf("usedtopic_%s_%d_%d", topicName, fromUserId, toUserId)
	} else {
		return fmt.Sprintf("usedtopic_%s_%d_%d", topicName, toUserId, fromUserId)
	}
}

func AddUsedUserTopic(fromUserId, toUserId int64, topicName string, topicId []int64) error {
	if RedisClient == nil {
		appzaplog.Error("AddUsedUserTopic RedisClient nil",
			zap.Int64("fromUserId", fromUserId),
			zap.Int64("toUserId", toUserId),
			zap.Int64s("topicid", topicId),
			zap.String("topicname", topicName),
		)
		return fmt.Errorf("RedisClient ni")
	}
	if len(topicId) == 0 {
		return nil
	}
	var topicIdStr []string
	for _, v := range topicId {
		topicIdStr = append(topicIdStr, strconv.FormatInt(v, 10))
	}
	err := RedisClient.SAdd(usedUserTopicKey(fromUserId, toUserId, topicName), topicIdStr).Err()
	if err != nil {
		appzaplog.Error("AddUsedUserTopic SAdd err", zap.Error(err),
			zap.Int64("fromUserId", fromUserId),
			zap.Int64("toUserId", toUserId),
			zap.Int64s("topicid", topicId),
			zap.String("topicname", topicName),
		)
		return err
	}
	err = RedisClient.Expire(usedUserTopicKey(fromUserId, toUserId, topicName), time.Hour*3).Err()
	if err != nil {
		appzaplog.Error("AddUsedUserTopic Expire err",
			zap.Int64("fromUserId", fromUserId),
			zap.Int64("toUserId", toUserId),
			zap.Int64s("topicid", topicId),
			zap.String("topicname", topicName),
			zap.Error(err))
		return err
	}
	return nil
}

func GetUsedUserTopic(fromUserId, toUserId int64, topicName string) ([]int64, error) {
	if RedisClient == nil {
		appzaplog.Error("GetUsedUserTopic RedisClient nil",
			zap.Int64("fromUserId", fromUserId),
			zap.Int64("toUserId", toUserId),
			zap.String("topicname", topicName),
		)
		return nil, fmt.Errorf("RedisClient ni")
	}
	usedtopic, err := RedisClient.SMembers(usedUserTopicKey(fromUserId, toUserId, topicName)).Result()
	if err != nil {
		appzaplog.Error("GetUsedUserTopic SMembers err",
			zap.Int64("fromUserId", fromUserId),
			zap.Int64("toUserId", toUserId),
			zap.String("topicname", topicName),
		)
		return nil, err
	}
	var ret []int64
	for _, v := range usedtopic {
		topic, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			appzaplog.Error("GetUsedUserTopic SMembers err",
				zap.Int64("fromUserId", fromUserId),
				zap.Int64("toUserId", toUserId),
				zap.String("topicname", topicName),
				zap.String("topic", v),
			)
			continue
		}
		ret = append(ret, topic)
	}
	return ret, nil
}

func DelUsedUserTopic(fromUserId, toUserId int64, topicName string) error {
	if RedisClient == nil {
		appzaplog.Error("DelUsedUserTopic RedisClient nil",
			zap.Int64("fromUserId", fromUserId),
			zap.Int64("toUserId", toUserId),
			zap.String("topicname", topicName),
		)
		return fmt.Errorf("RedisClient ni")
	}
	return RedisClient.Del(usedUserTopicKey(fromUserId, toUserId, topicName)).Err()
}
