package handler

import (
	"bilin/chattag/dao"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"github.com/gin-gonic/gin"
	"strings"
)

type PUATopic struct {
	Topics []TopicElement `json:"topics"`
}

type TopicElement struct {
	TopicId uint     `json:"topic_id"`
	Lines   []string `json:"lines"`
}

const puatopic = "puatopic"

func HandlePUATopic(c *gin.Context) *HttpError {
	appzaplog.Debug("[+]HandlePUATopic")
	var (
		body *PUATopic = &PUATopic{}
		ret            = successHttp
	)

	for {
		//uidStr := c.PostForm("to_userid")
		//toUid, err := strconv.ParseInt(uidStr, 10, 64)
		//if err != nil {
		//	appzaplog.Error("HandlePUATopic ParseInt err", zap.Error(err), zap.String("uid", uidStr))
		//	ret = parseFormHttpErr
		//	break
		//}
		//usedtopic, err := dao.GetUsedUserTopic(0, toUid, puatopic)
		//if err != nil {
		//	appzaplog.Error("HandlePUATopic ParseInt err", zap.Error(err), zap.String("uid", uidStr))
		//	ret = daoGetHttpErr
		//	break
		//}
		//info, err := dao.GetTopicNotIn(usedtopic, 1)
		//if err != nil {
		//	ret = daoGetHttpErr
		//	appzaplog.Error("HandlePUATopic GetTopicNotIn err", zap.Error(err), zap.String("uid", uidStr))
		//	break
		//}
		//
		//// recycle
		//if (info == nil || len(info) == 0) && len(usedtopic) > 0 {
		//
		//	if err := dao.DelUsedUserTopic(0, toUid, puatopic); err != nil {
		//		appzaplog.Error("HandlePUATopic DelUsedUserTopic err", zap.Error(err), zap.Int64("uid", toUid))
		//		ret = delUsedTopicHttpErr
		//		break
		//	}
		//	usedtopic = []int64{}
		//	info, err = dao.GetTopicNotIn(usedtopic, 1)
		//	if err != nil {
		//		ret = daoGetHttpErr
		//		appzaplog.Error("HandlePUATopic GetTopicNotIn err", zap.Error(err))
		//		break
		//	}
		//}
		//
		//// add it to used
		//var toAddTopicId []int64
		//for _, v := range info {
		//	toAddTopicId = append(toAddTopicId, int64(v.ID))
		//}
		//if err := dao.AddUsedUserTopic(0, toUid, puatopic, toAddTopicId); err != nil {
		//	ret = daoPutHttpErr
		//	appzaplog.Error("HandlePUATopic AddUsedUserTopic err", zap.Error(err), zap.String("uid", uidStr))
		//	break
		//}
		info, err := dao.RandPuaTopic(1)
		if err != nil {
			ret = daoGetHttpErr
			appzaplog.Error("HandlePUATopic RandPuaTopic err", zap.Error(err))
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
	appzaplog.Debug("[-]HandleRCPlugin", zap.Any("resp", httpret))
	return ret
}
