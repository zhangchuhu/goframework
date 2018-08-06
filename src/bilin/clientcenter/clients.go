/*
client统一放置的地方
*/
package clientcenter

import (
	"bilin/protocol"
	"bilin/protocol/userinfocenter"
	"context"

	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/tars"
)

var (
	confCenterClient    bilin.ConfInfoServantClient
	roomcCenterClient   bilin.RoomInfoServantClient
	userInfoClient      userinfocenter.UserInfoCenterObjClient
	chatTagClient       bilin.ChatTagTarsClient
	bizRoomCenterClient bilin.BizRoomCenterServantClient
	guildClient         bilin.GuildTarsClient
)

func init() {
	comm := tars.NewCommunicator()
	//comm.SetProperty("locator", "tars.tarsregistry.QueryObj@tcp -h 183.36.111.89 -p 17890") /*测试环境*/
	//comm.SetProperty("locator", "tars.tarsregistry.QueryObj@tcp -h 183.36.111.61 -p 17890") /*正式环境，谨慎用！*/
	confCenterClient = bilin.NewConfInfoServantClient("bilin.confinfocenter.ConfInfoServantObj", comm)
	roomcCenterClient = bilin.NewRoomInfoServantClient("bilin.roominfocenter.RoomInfoCenterServantObj", comm)
	userInfoClient = userinfocenter.NewUserInfoCenterObjClient("bilin.userinfocenter.UserInfoCenterObj", comm)
	chatTagClient = bilin.NewChatTagTarsClient("bilin.chattagtars.ChatTagTarsObj", comm)
	bizRoomCenterClient = bilin.NewBizRoomCenterServantClient("bilin.bizroomcenter.BizRoomCenterPbObj", comm)
	guildClient = bilin.NewGuildTarsClient("bilin.guildtars.GuildTarsObj", comm)
}

func ConfClient() bilin.ConfInfoServantClient {
	return confCenterClient
}

func RoomCenterClient() bilin.RoomInfoServantClient {
	return roomcCenterClient
}

func UserInfoClient() userinfocenter.UserInfoCenterObjClient {
	return userInfoClient
}

func ChatTagClient() bilin.ChatTagTarsClient {
	return chatTagClient
}

func BizRoomCenterClient() bilin.BizRoomCenterServantClient {
	return bizRoomCenterClient
}

func GuildTarsClient() bilin.GuildTarsClient {
	return guildClient
}

const maxuserinfonum = 20

func TakeUserInfo(userids []uint64) (map[uint64]*userinfocenter.UserInfo, error) {
	useridlen := len(userids)
	if useridlen == 0 {
		return nil, nil
	}

	var ret = make(map[uint64]*userinfocenter.UserInfo, useridlen)
	round := useridlen % maxuserinfonum
	yu := useridlen / maxuserinfonum

	for i := 0; i < yu; i++ {
		userinfos, err := UserInfoClient().GetUserInfo(context.TODO(), &userinfocenter.GetUserInfoReq{
			Uids: userids[i*maxuserinfonum : (i+1)*maxuserinfonum],
		})
		if err != nil {
			log.Error("GetUserInfo failed", zap.Error(err))
			continue
		}
		for k, v := range userinfos.Users {
			ret[k] = v
		}
	}
	if round > 0 {
		userinfos, err := UserInfoClient().GetUserInfo(context.TODO(), &userinfocenter.GetUserInfoReq{
			Uids: userids[yu*maxuserinfonum:],
		})
		if err != nil {
			log.Error("GetUserInfo failed", zap.Error(err))
		} else {
			for k, v := range userinfos.Users {
				ret[k] = v
			}
		}
	}
	return ret, nil
}

func TakeRoomInfo(roomids []uint64) (ret map[uint64]*bilin.RoomInfo, err error) {
	if len(roomids) == 0 {
		return
	}
	rsp, err := RoomCenterClient().LivingRoomsInfo(context.TODO(), &bilin.LivingRoomsInfoReq{})
	if err != nil {
		log.Error("GetRoomInfo failed", zap.Error(err))
		return
	}
	ret = make(map[uint64]*bilin.RoomInfo, len(roomids))
	for _, rid := range roomids {
		if info, ok := rsp.Livingrooms[rid]; ok {
			ret[rid] = info
		}
	}
	return
}

func BatchGetBizRoomInfo(roomids []uint64) (ret map[uint64]*bilin.BizRoomInfo, err error) {
	if len(roomids) == 0 {
		return
	}
	rsp, err := BizRoomCenterClient().BatchGetBizRoomInfo(context.TODO(), &bilin.BatchGetBizRoomInfoReq{Roomids: roomids})
	if err != nil {
		log.Error("GetRoomInfo failed", zap.Error(err))
		return
	}

	return rsp.Bizroominfos, err
}
