/*
 * Copyright (c) 2018-07-16.
 * Author: kordenlu
 * 功能描述:${<VARIABLE_NAME>}
 */

package handler

import (
	"bilin/clientcenter"
	"bilin/protocol"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"context"
)

func fillUserChatTag(bc *bilin.OptionalMatchingResult) error {
	if bc == nil {
		return nil
	}
	var toptaguser []*bilin.TopNUser
	for _, v := range bc.Attendees {
		toptaguser = append(toptaguser, &bilin.TopNUser{
			Touserid: int64(v.Uid),
			Topn:     2,
		})
	}
	taginfo, err := clientcenter.ChatTagClient().BatchRTopNUserChatTagSummary(context.TODO(), &bilin.BatchRTopNUserChatTagSummaryReq{
		Topnuser: toptaguser,
	})
	if err != nil {
		appzaplog.Error("fillUserChatTag BatchRTopNUserChatTagSummary err", zap.Error(err))
		return err
	}

	for _, v := range bc.Attendees {
		if info, ok := taginfo.Summary[int64(v.Uid)]; ok {
			v.UserChatTag = info
		}
	}
	return nil
}
