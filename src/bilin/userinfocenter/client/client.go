package main

import (
	"bilin/protocol/userinfocenter"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/tars"
	"context"
	"fmt"
	"time"
)

func main() {
	comm := tars.NewCommunicator()
	//objName := fmt.Sprintf("bilin.userinfocenter.UserInfoCenterObj@tcp -h 183.36.111.89 -t 60000 -p 12002")
	objName := fmt.Sprintf("bilin.userinfocenter.UserInfoCenterObj@tcp -h 127.0.0.1 -t 60000 -p 20008")
	client := userinfocenter.NewUserInfoCenterObjClient(objName, comm)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	// list, err := client.GetAvatarUserInfo(ctx, &userinfocenter.GetAvatarUserInfoReq{})
	// if err != nil {
	// 	appzaplog.Error("GetAvatarUserInfo err", zap.Error(err))
	// 	return
	// }
	// for _, v := range list.Users {
	// 	fmt.Println(v.Avatar)
	// }

	resp, err := client.GetUserInfo(ctx, &userinfocenter.GetUserInfoReq{Uids: []uint64{699, 17795069}})
	if err != nil {
		appzaplog.Error("GetUserInfo err", zap.Error(err))
		return
	}
	appzaplog.Info("GetUserInfo resp msg", zap.Any("resp", resp))

	check_resp, err := client.IsAppleCheckUser(ctx, &userinfocenter.IsAppleCheckUserReq{Uid:17795069})
	if err != nil {
		appzaplog.Error("IsAppleCheckUser err", zap.Error(err))
		return
	}
	appzaplog.Info("IsAppleCheckUser resp msg", zap.Any("check_resp", check_resp))
}
	



	
