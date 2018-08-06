package controller

import (
	"bilin/clientcenter"
	"bilin/protocol/userinfocenter"
	"bilin/relationlist/entity"
	"bilin/relationlist/service"
	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
)

//获取用户信息，稍微封装一下
func BatchGetUserInfo(uids []uint64) (ret map[uint64]*userinfocenter.UserInfo, err error) {
	const prefix = "BatchGetUserInfo "
	var queryUserids []uint64

	ret = make(map[uint64]*userinfocenter.UserInfo)
	for _, uid := range uids {
		//先从redis查，查不到才去userinfocenter获取数据,本地缓存
		user, _ := service.RedisGetUserInfo(uid)
		if user != nil {
			ret[uid] = user
			continue
		}

		//not found ,query userinfo center
		queryUserids = append(queryUserids, uid)
	}

	usersInfo, err := clientcenter.TakeUserInfo(queryUserids)
	if err != nil {
		log.Error(prefix+"TakeUserInfo failed", zap.Any("queryUserids", queryUserids))
		return
	}

	for key, value := range usersInfo {
		//写redis
		service.RedisSetUserInfo(value)
		ret[key] = value
	}

	log.Info(prefix+"end", zap.Any("ret", ret))
	return
}

func fillContent(owner uint64, result *entity.RelationStatistics) (err error) {
	var uids []uint64
	uids = append(uids, owner)
	for _, item := range result.RelationList {
		uids = append(uids, item.UserID)
	}

	usersInfo, _ := BatchGetUserInfo(uids)
	user, exist := usersInfo[result.AnchorInfo.UserID]
	if !exist {
		log.Error("userinfo not find", zap.Any("uid", result.AnchorInfo.UserID))
	}
	result.AnchorInfo.Avatar = user.Avatar
	result.AnchorInfo.Nick = user.NickName

	for _, item := range result.RelationList {
		user, exist := usersInfo[item.UserID]
		if !exist {
			log.Error("userinfo not find", zap.Any("uid", item.UserID))
			continue
		}

		item.Avatar = user.Avatar
		item.Nick = user.NickName
	}

	return
}

func GetDailyRelationList(owner uint64, start int, rows int) (result *entity.RelationStatistics) {
	const prefix = "GetDailyRelationList "
	//result, _ = service.RedisGetDailyStatisticsRelationList(owner)

	//直接从数据库中获取，后面再优化
	result, _ = service.MysqlGetDailyStatisticsRelationList(owner)

	fillContent(owner, result)

	result.Start = start
	result.NumFound = len(result.RelationList)
	log.Info(prefix+"end", zap.Any("owner", owner), zap.Any("result", result))
	return
}

func GetWeeklyRelationList(owner uint64, start int, rows int) (result *entity.RelationStatistics) {
	const prefix = "GetWeeklyRelationList "
	//result, _ = service.RedisGetWeeklyStatisticsRelationList(owner)
	//if result != nil { //缓存中有数据，直接返回
	//	log.Info(prefix+"end", zap.Any("owner", owner), zap.Any("result", result))
	//	return
	//}

	//缓存中没有查到数据，需要从db中获取
	result, _ = service.MysqlGetWeeklyStatisticsRelationList(owner)
	//if result == nil { //db中没有查找到
	//	//set一个只有主播的数据进去，防止频繁查询数据库
	//	result = &entity.RelationStatistics{}
	//	result.AnchorInfo = &entity.UserRelationInfo{UserID: owner}
	//}

	fillContent(owner, result)

	//周榜用户需要根据亲密度分配勋章
	cfg := GetMedalConfig()
	for index, item := range result.RelationList {
		if index < 3 && item.RelationVal >= 10 {
			item.MedalUrl = cfg[int32(index+1)].MedalUrl
			item.MedalText = cfg[int32(index+1)].MedalName
		}
	}

	//写redis，防止频繁查询db
	service.RedisSetWeeklyStatisticsRelationList(owner, result)

	result.Start = start
	result.NumFound = len(result.RelationList)
	log.Info(prefix+"end", zap.Any("owner", owner), zap.Any("result", result))
	return
}

func GetTotalRelationList(owner uint64, start int, rows int) (result *entity.RelationStatistics) {
	const prefix = "GetTotalRelationList "

	//result, _ = service.RedisGetTotalStatisticsRelationList(owner)

	//直接从数据库中获取，后面再优化
	result, _ = service.MysqlGetTotalStatisticsRelationList(owner)
	fillContent(owner, result)

	result.Start = start
	result.NumFound = len(result.RelationList)
	log.Info(prefix+"end", zap.Any("owner", owner), zap.Any("result", result))
	return
}
