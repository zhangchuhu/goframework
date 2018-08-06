package main

import (
	"context"
	"bilin/protocol"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"fmt"
	"code.yy.com/yytars/goframework/tars"
	"os"
)

func main() {
	var ip, port string
	if len(os.Args) == 3 {
		ip = os.Args[1]
		port = os.Args[2]
	} else {
		fmt.Errorf("input error:")
		fmt.Println("./client ip port")
		return
	}
	comm := tars.NewCommunicator()
	objName := fmt.Sprintf("bilin.bigexpression.BigExpressionObjObj@tcp -h %s -p %s -t 6000", ip, port)
	client := bilin.NewBigExpressionObjClient(objName,comm)
	res_config, err := client.GetEmotionConfig(context.TODO(), &bilin.GetEmotionConfigReq{})
	if err != nil {
		appzaplog.Error("GetEmotionConfig err",zap.Error(err))
		return
	} else {
		appzaplog.Info("GetEmotionConfig",zap.Any("result:",res_config))
	}
	//普通表情
	req :=  &bilin.SendEmotionReq{
		&bilin.Header{
			Userid: uint64(10000),
			Roomid: 400000367,
		},
		&bilin.Emotion{
			Id: uint32(1),
		},
	}
	res_send, err := client.SendEmotion(context.TODO(), req)
	if err != nil {
		appzaplog.Error("SendEmotion err,id:1",zap.Error(err))
		return
	} else {
		appzaplog.Info("SendEmotion,id:1",zap.Any("result:",res_send))
	}
	//
	req.Emotion.Id=2
	res_send, err = client.SendEmotion(context.TODO(), req)
	if err != nil {
		appzaplog.Error("SendEmotion err,id:2",zap.Error(err))
		return
	} else {
		appzaplog.Info("SendEmotion,id:2",zap.Any("result:",res_send))
	}
	//随机表情,result_count=1
	req.Emotion.Id=18  //骰子
	res_send, err = client.SendEmotion(context.TODO(), req)
	if err != nil {
		appzaplog.Error("SendEmotion err,id:18",zap.Error(err))
		return
	} else {
		appzaplog.Info("SendEmotion,id:18",zap.Any("result:",res_send))
	}
	req.Emotion.Id=17
	res_send, err = client.SendEmotion(context.TODO(), req)
	if err != nil {
		appzaplog.Error("SendEmotion err,id:17",zap.Error(err))
		return
	} else {
		appzaplog.Info("SendEmotion,id:17",zap.Any("result:",res_send))
	}
	//随机表情,result_count>1
	req.Emotion.Id=23  //扑克牌
	res_send, err = client.SendEmotion(context.TODO(), req)
	if err != nil {
		appzaplog.Error("SendEmotion err,id:23",zap.Error(err))
		return
	} else {
		appzaplog.Info("SendEmotion,id:23",zap.Any("result:",res_send))
	}
}
