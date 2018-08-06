package handler

import (
	"context"
	"bilin/protocol"
	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"github.com/golang/protobuf/proto"
	"math/rand"
	"bilin/common/onlinepush"
	"bilin/bigexpression/config"
	"errors"
)
type BigExpressionObjObj struct {
}
func RangeNum(min, max int) int {
	randNum := rand.Intn(max - min)
	randNum = randNum + min
	return randNum
}

func NewBigExpressionObjObj() *BigExpressionObjObj {
	return &BigExpressionObjObj{}
}

func (p *BigExpressionObjObj) GetEmotionConfig(ctx context.Context, r *bilin.GetEmotionConfigReq) (resp *bilin.GetEmotionConfigRes,err error) {
	return config.GetAppConfig(),nil
}

func (p *BigExpressionObjObj) SendEmotion(ctx context.Context, req *bilin.SendEmotionReq) (res *bilin.SendEmotionRes,err error) {
	var (
		result_index_start uint32
		result_index_end uint32
		result_count uint64
		emotion_type  bilin.EmotionType
		exist = false
		BigExpressionBroadCast bilin.BigExpressionBroadcast
	)
	EmotionConfigRes := config.GetAppConfig()
	res =&bilin.SendEmotionRes{&bilin.Emotion{}}
	res.Emotion.Id = req.Emotion.Id
	BigExpressionBroadCast.Id = req.Emotion.Id
	BigExpressionBroadCast.FromUid = req.GetHeader().GetUserid()
	for _,emotion_config := range EmotionConfigRes.GetEmotionConfig(){
		if emotion_config.Id == req.Emotion.Id{
			result_index_start = emotion_config.ResultIndexStart
			result_index_end = emotion_config.ResultIndexEnd
			emotion_type = emotion_config.Type
			result_count = emotion_config.ResultCount
			exist = true
		}
	}
	if !exist {
		log.Error("invalid request", zap.Any("id:",req.GetEmotion().Id))
		return res, errors.New("wrong emotion id")
	}
	if emotion_type == 1 { //心情表情
		res.Emotion.ResultIndex = result_index_start
		BigExpressionBroadCast.ResultIndex=res.Emotion.ResultIndex
	} else if emotion_type == 2{ //随机事件表情
		if result_count == 1 {
			if result_index_start == result_index_end {
				res.Emotion.ResultIndex = result_index_start
			} else {
				res.Emotion.ResultIndex=(uint32)(RangeNum(int(result_index_start),int(result_index_end)))
			}
			BigExpressionBroadCast.ResultIndex=res.Emotion.ResultIndex
		} else if result_count > 1{
			for i:=0;i< (int)(result_count);i++ {
				result :=(uint32)(RangeNum(int(result_index_start),int(result_index_end)))
				res.Emotion.ResultIndexes=append(res.Emotion.ResultIndexes, result)
			}
			BigExpressionBroadCast.ResultIndexes=res.Emotion.ResultIndexes
		} else {
			res.Emotion.ResultIndex=0
			log.Error("config wrong", zap.Any("emotionConfig:",EmotionConfigRes))
			return res, errors.New("config wrong")
		}
	} else {
		res.Emotion.ResultIndex=0
		log.Error("config wrong", zap.Any("emotionConfig:",EmotionConfigRes))
		return res, errors.New("config wrong")
	}
	var body bilin.BcMessageBody
	body.Type = int32(bilin.MinType_BC_NotifyBigExpression)
	body.Data,err = proto.Marshal(&BigExpressionBroadCast)
	if err != nil {
		log.Error("proto.Marshal failed", zap.Any("err", err))
		return res, err
	}
	pushData,err := proto.Marshal(&body)
	if err != nil {
		log.Error("proto.Marshal failed", zap.Any("err", err))
		return res, err
	}
	push := bilin.ServerPush{
		MessageType: int32(bilin.MaxType_BC_MSG),
		PushBuffer:  pushData,
		MessageDesc: "麦上大表情",
		ServiceName: "bigexpression",
		MethodName:  "pushBroadcastToRoom",
	}
	log.Error("PushToRoom", zap.Any("roomid", req.GetHeader().GetRoomid()))
	err = onlinepush.PushToRoom(int64(req.GetHeader().GetRoomid()), push)
	if err != nil {
		log.Error("push to room failed", zap.Any("err", err))
		return res, err
	}
	return res,nil
}