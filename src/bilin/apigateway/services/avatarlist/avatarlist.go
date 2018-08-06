package avatarlist

import (
	"bilin/clientcenter"
	"bilin/common/cacheprocessor"
	"bilin/protocol"
	"bilin/protocol/userinfocenter"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"context"
	"errors"
	"sync/atomic"
	"time"
)

const (
	FEMALE = 0
	MALE   = 1
)

type AvatarList struct {
	MaleAvatarURL   []string
	FemaleAvatarURL []string
}

var SafeAvatarList atomic.Value

const (
	FEMALE_AVATAR1 = "https://vipweb.bs2cdn.yy.com/vipinter_83ac9c823dd24ed2adfa173ae5480c1b.jpg"
	FEMALE_AVATAR2 = "https://vipweb.bs2cdn.yy.com/vipinter_9c461d5f5122495aa3c4f644c7fef643.jpg"
	FEMALE_AVATAR3 = "https://vipweb.bs2cdn.yy.com/vipinter_269b2f1a915c4cc6a1099c1567b72455.jpg"
	MALE_AVATAR1   = "https://vipweb.bs2cdn.yy.com/vipinter_9588e2f6be93493eb5610ae1498acc87.jpg"
	MALE_AVATAR2   = "https://vipweb.bs2cdn.yy.com/vipinter_f689577159904b63a960bb74e40cdcde.jpg"
	MALE_AVATAR3   = "https://vipweb.bs2cdn.yy.com/vipinter_3650f0815b9f4f36b72d7f77f95fe8ce.jpg"
)

func init() {
	if err := cacheprocessor.CacheProcessor("LoadAvatar", 300*time.Second, LoadAvatar); err != nil {
		panic("LoadCarouselIndex failed")
	}
}

func GetAvatarList(sex int64) []string {
	avatar_list := SafeAvatarList.Load().(*AvatarList)
	var use_list []string
	if sex == MALE {
		use_list = avatar_list.FemaleAvatarURL
	} else {
		use_list = avatar_list.MaleAvatarURL
	}

	list_len := len(use_list)
	if list_len == 0 {
		return use_list //这里要不要增加默认值 还是客户端使用默认
	}

	rondom := time.Now().Unix()
	ret_list := make([]string, 0, 3) //只有小于三个时候才会重复
	for i := 0; i < 3; i++ {
		index := uint64(rondom) % uint64(list_len)
		ret_list = append(ret_list, use_list[index])
		rondom++
	}
	return ret_list
}

func LoadAvatar() error {
	var avatar_list AvatarList
	avatar_list.MaleAvatarURL = []string{MALE_AVATAR1, MALE_AVATAR2, FEMALE_AVATAR3}
	avatar_list.FemaleAvatarURL = []string{FEMALE_AVATAR1, FEMALE_AVATAR2, FEMALE_AVATAR3}
	SafeAvatarList.Store(&avatar_list)
	//return nil
	//原来取正在直播的用户的头像展示 这个暂时不这么做了
	//保留自动拉取和更新的代码
	room_info_resp, err := clientcenter.RoomCenterClient().LivingRoomsInfo(context.TODO(), &bilin.LivingRoomsInfoReq{})
	if err != nil {
		appzaplog.Error("get LivingRoomsInfo fail", zap.Error(err))
		SafeAvatarList.Store(&avatar_list)
		return nil
	}

	if room_info_resp.Livingrooms == nil {
		appzaplog.Info("get LivingRoomsInfo resp nil")
		SafeAvatarList.Store(&avatar_list)
		return nil
	}

	//appzaplog.Debug("get LivingRoomsInfo success", zap.Any("room_info_resp.Livingrooms", room_info_resp.Livingrooms))

	var i int = 0
	uids := []uint64{}
	for _, room := range room_info_resp.Livingrooms {
		//每次查询20个用户  最多查询10次
		if room.Owner == 0 {
			continue
		}

		uids = append(uids, room.Owner)
		if len(uids) < 20 { //满20个才差
			continue
		}

		err := avatar_list.AddAvatarList(uids)
		if err != nil {
			return nil //更新失败则使用默认的
		}

		uids = []uint64{} //清空
		i++
		if i == 10 {
			break //不需要查太多 200个足够了
		}

	}

	if len(uids) > 0 { //有数据才查
		//最后一次不足20个的还没有请求
		avatar_list.AddAvatarList(uids) //这里可能不足20个uid
	}

	appzaplog.Debug("LoadAvatar success, reload data")
	SafeAvatarList.Store(&avatar_list)
	return nil
}

func (a *AvatarList) AddAvatarList(uids []uint64) error {

	resp, err := clientcenter.UserInfoClient().GetUserInfo(context.TODO(), &userinfocenter.GetUserInfoReq{uids})
	if err != nil {
		appzaplog.Error("AddAvatarList GetUserInfo err", zap.Error(err))
		return err
	}

	if resp.Ret == nil || resp.Ret.Code != userinfocenter.Result_SUCCESS {
		appzaplog.Error("AddAvatarList GetUserInfo err")
		return errors.New("AddAvatarList GetUserInfo err")
	}

	if resp.Users == nil {
		appzaplog.Warn("AddAvatarList, no users", zap.Any("uids", uids))
		return nil
	}

	//appzaplog.Info("AddAvatarList GetUserInfo", zap.Any("resp", resp))

	for _, uid := range uids {
		if v, ok := resp.Users[uid]; ok {
			if v.Avatar == "" {
				continue
			}
			if v.Sex == FEMALE {
				a.FemaleAvatarURL = append(a.FemaleAvatarURL, v.Avatar)
			} else {
				a.MaleAvatarURL = append(a.MaleAvatarURL, v.Avatar)
			}
		}
	}
	return nil
}
