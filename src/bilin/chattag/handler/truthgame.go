package handler

import (
	"bilin/chattag/dao"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"github.com/gin-gonic/gin"
	"strings"
)

type TruthGame struct {
	Topics []TopicElement `json:"topics"`
}

const truthtopic = "truthtopic"

func HandleTruthGame(c *gin.Context) *HttpError {
	appzaplog.Debug("[+]HandleTruthGame")
	var (
		body *TruthGame = &TruthGame{}
		ret             = successHttp
	)

	for {
		//toUidStr := c.PostForm("to_userid")
		//toUid, err := strconv.ParseInt(toUidStr, 10, 64)
		//if err != nil {
		//	appzaplog.Error("HandleTruthGame ParseInt err", zap.Error(err), zap.String("uid", toUidStr))
		//	ret = parseFormHttpErr
		//	break
		//}
		//fromUid := c.GetInt64(uidContextKey)
		//
		//usedtopic, err := dao.GetUsedUserTopic(fromUid, toUid, truthtopic)
		//if err != nil {
		//	appzaplog.Error("HandleTruthGame ParseInt err", zap.Error(err), zap.Int64("uid", toUid))
		//	ret = daoGetHttpErr
		//	break
		//}
		//
		//info, err := dao.GetTruthTopicNotIn(usedtopic, 1)
		//if err != nil {
		//	ret = daoGetHttpErr
		//	appzaplog.Error("HandleTruthGame GetTruthTopicNotIn err", zap.Error(err))
		//	break
		//}
		//
		////recycle
		//if (info == nil || len(info) == 0) && len(usedtopic) > 0 {
		//
		//	if err := dao.DelUsedUserTopic(fromUid, toUid, truthtopic); err != nil {
		//		appzaplog.Error("HandleTruthGame DelUsedUserTopic err", zap.Error(err), zap.Int64("uid", toUid))
		//		ret = delUsedTopicHttpErr
		//		break
		//	}
		//	usedtopic = []int64{}
		//	info, err = dao.GetTruthTopicNotIn(usedtopic, 1)
		//	if err != nil {
		//		ret = daoGetHttpErr
		//		appzaplog.Error("HandleTruthGame GetTruthTopicNotIn err", zap.Error(err))
		//		break
		//	}
		//}
		//
		//// add it to used
		//var toAddTopicId []int64
		//for _, v := range info {
		//	toAddTopicId = append(toAddTopicId, int64(v.ID))
		//}
		//if err := dao.AddUsedUserTopic(fromUid, toUid, truthtopic, toAddTopicId); err != nil {
		//	ret = daoPutHttpErr
		//	appzaplog.Error("HandlePUATopic AddUsedUserTopic err", zap.Error(err),
		//		zap.Int64("fromUid", fromUid),
		//		zap.Int64("toUid", toUid),
		//	)
		//	break
		//}
		info, err := dao.RandTruthTopic(1)
		if err != nil {
			ret = daoGetHttpErr
			appzaplog.Error("HandleTruthGame RandTruthTopic err", zap.Error(err))
			break
		}
		for _, v := range info {
			elem := TopicElement{
				TopicId: v.ID,
				Lines:   strings.Split(v.Topic, "\n"),
			}
			body.Topics = append(body.Topics, elem)
		}
		break
	}

	httpret := HttpRetComm{
		IsEncrypt: "false",
		Data: HttpRetDataComm{
			Code: ret.code,
			Msg:  ret.desc,
			Body: body,
		},
	}
	c.JSON(200, httpret)
	appzaplog.Debug("[-]HandleTruthGame", zap.Any("resp", httpret))
	return ret
}
