/*
缓存通用的用户信息
*/
package cache

import (
	"bilin/clientcenter"
	"bilin/protocol"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"context"
	"errors"
	"time"
)

func refreshCommonRecInfo() error {
	resp, err := clientcenter.RoomCenterClient().LivingRoomsInfo(context.TODO(), &bilin.LivingRoomsInfoReq{})
	if err != nil {
		appzaplog.Error("LivingRoomsInfo faied", zap.Error(err))
		return err
	}
	var (
		userids []uint64
	)
	for _, v := range resp.Livingrooms {
		userids = append(userids, v.Owner)
	}
	userinfos, err := clientcenter.TakeUserInfo(userids)
	if err != nil {
		appzaplog.Error("refresh GetUserInfo faied", zap.Error(err))
		return err
	}

	badges, err := BadgeInfo()
	if err != nil {
		appzaplog.Error("refresh BadgeInfo faied", zap.Error(err))
		//return err
	}
	nowunix := time.Now().Unix()
	commonrec := make(map[uint64]*RecommandLivingInfo, len(resp.Livingrooms))
	for k, roominfo := range resp.Livingrooms {
		if roominfo.LockStatus == 0 {
			if userinfo, uok := userinfos[roominfo.Owner]; uok {
				commonrec[k] = &RecommandLivingInfo{
					TypeID:     1,
					LiveID:     roominfo.Roomid,
					Title:      roominfo.Title,
					UserCount:  roominfo.Usernumber,
					CategoryId: uint64(roominfo.RoomcategoryID),
					StartTime:  roominfo.Starttime,
					SortWeight: calcSortWeight(takelivingtime(nowunix, int64(roominfo.Starttime)), int64(roominfo.Usernumber)),
					LivingHostInfo: LivingHostInfo{
						UserID:   roominfo.Owner,
						NickName: userinfo.NickName,
						City:     userinfo.City,
						SmallURL: tosmall(userinfo.Avatar),
					},
				}
				if badges != nil {
					for _, badge := range badges {
						if badge.Userid == commonrec[k].LivingHostInfo.UserID {
							commonrec[k].LivingHostInfo.TagUrl = append(commonrec[k].LivingHostInfo.TagUrl, badge.Url)
						}
					}
				}
			}
		}

	}
	caches.Store(commRecInfoCacheKey, commonrec)
	return nil
}

// 根据分类获取推荐的原始数据
func getCommonRecInfo(typeid uint64) (map[uint64]*RecommandLivingInfo, error) {
	if v, ok := caches.Load(commRecInfoCacheKey); ok {
		if ret, yes := v.(map[uint64]*RecommandLivingInfo); yes {
			myret := make(map[uint64]*RecommandLivingInfo)
			for k, v := range ret {
				if typeid == hotcategoryid || v.CategoryId == typeid {
					myret[k] = v
				}
			}
			return myret, nil
		}
	}
	return nil, errors.New("not found")
}
