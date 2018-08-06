package handler

import (
	"bilin/protocol/userinfocenter"
	// "bilin/userinfocenter/dao"
	// log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	// "code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"context"
	//"strings"
	//"time"
)

const (
	AVATAR_USER_NUM = 60
)

func (this *UserInfoCenterObj) GetAvatarUserInfo(ctx context.Context, req *userinfocenter.GetAvatarUserInfoReq) (*userinfocenter.GetAvatarUserInfoResp, error) {

	//首页的头像已经不再从数据库加载 这个接口不使用了
	return nil, nil

	// log.Debug("GetRondomAvatarUsers")

	// var (
	// 	resp = &userinfocenter.GetAvatarUserInfoResp{
	// 		Users: make([]*userinfocenter.UserInfo, 0, AVATAR_USER_NUM),
	// 	}
	// )

	// try_times := 3
	// for i := 0; i < try_times; i++ {
	// 	rand := time.Now().Unix() % 100
	// 	avatar_list, err := dao.GetUserAvatatrInfos(uint64(rand), AVATAR_USER_NUM)
	// 	if err != nil {
	// 		log.Error("GetUserAvatatrInfos failed", zap.Error(err))
	// 		return nil, err
	// 	}

	// 	if avatar_list != nil {
	// 		for _, v := range avatar_list {
	// 			var user userinfocenter.UserInfo
	// 			user.Avatar = v.GetAvatar()
	// 			user.Uid = v.UserId
	// 			resp.Users = append(resp.Users, &user)
	// 		}
	// 	}

	// 	if len(resp.Users) >= AVATAR_USER_NUM {
	// 		break
	// 	}
	// }

	// log.Debug("GetRondomAvatarUsers success", zap.Any("resp", *resp))
	// return resp, nil
}
