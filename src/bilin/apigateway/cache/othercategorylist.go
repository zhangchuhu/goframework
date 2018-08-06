/*
除热门外其他品类的展示推荐逻辑
*/
package cache

import (
	"bilin/apigateway/services/rank"
	"bilin/clientcenter"
	"bilin/protocol"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"context"
	"sort"
	"strconv"
)

func hostRec(typeid uint64) ([]*bilin.CategoryHostRecInfo, error) {
	hostrec, err := clientcenter.ConfClient().CategoryHostRec(context.TODO(), &bilin.CategoryHostRecReq{})
	if err != nil {
		appzaplog.Error("LivingCategorys failed", zap.Error(err))
		return nil, err
	}
	var ret []*bilin.CategoryHostRecInfo
	for _, v := range hostrec.Cateogryinfos {
		if v.Typeid == typeid {
			ret = append(ret, v)
		}
	}
	return ret, nil
}

func refreshOtherCategorylist() error {

	resp, err := clientcenter.ConfClient().LivingCategorys(context.TODO(), &bilin.LivingCategorysReq{})
	if err != nil {
		appzaplog.Error("LivingCategorys failed", zap.Error(err))
		return err
	}

	last5mrank := rank.GetChannelBillRank()

	if resp != nil {
		for _, v := range resp.Livingcategorys {
			if v.Typeid == hotcategoryid {
				continue
			}
			if v.Typeid == multiplayercategoryid {
				if err := specialRec(last5mrank, v); err != nil {
					appzaplog.Error("getCommonRecInfo failed", zap.Error(err))
				}
				continue
			}
			comrecinfo, err := getCommonRecInfo(uint64(v.Typeid))
			if err != nil {
				appzaplog.Error("getCommonRecInfo failed", zap.Error(err))
				return err
			}
			body := &RecommandLivingBody{}
			stickpost, err := OntimeStickyPostList(v.Typeid)
			if err != nil {
				appzaplog.Error("OntimeStickyPostList failed", zap.Int64("typeid", v.Typeid), zap.Error(err))
				return err
			}
			//栏目主播私档
			hostrecs, err := hostRec(uint64(v.Typeid))
			if err != nil {
				appzaplog.Error("hostRec failed", zap.Int64("typeid", v.Typeid), zap.Error(err))
				return err
			}

			for _, post := range stickpost {
				if roominfo, ok := comrecinfo[post.Roomid]; ok {
					body.RecommandLivingList = append(body.RecommandLivingList, roominfo)
					delete(comrecinfo, post.Roomid)
				}
			}

			//私档
			var privatehost []*RecommandLivingInfo
			var privatehostmap map[uint64]*RecommandLivingInfo = make(map[uint64]*RecommandLivingInfo)
			for _, hostrec := range hostrecs {
				for _, rec := range comrecinfo {
					if rec.LivingHostInfo.UserID == hostrec.Hostid {
						if rank, ok := (*last5mrank)[hostrec.Hostid]; ok {
							rec.Last5MinRankValue = rank.Value
						}
						privatehost = append(privatehost, rec)
						privatehostmap[rec.LiveID] = rec
						break
					}
				}
			}
			sort.Sort(Latest5MinRankSlice(privatehost))
			top3 := topn(privatehost, 3)
			var privatemap = make(map[uint64]*RecommandLivingInfo)
			for _, private := range top3 {
				if private.Last5MinRankValue > 0 {
					privatemap[private.LiveID] = private
				}
			}

			sort.Sort(OnlineUserSlice(privatehost))
			top3 = topn(privatehost, 3)
			for _, private := range top3 {
				privatemap[private.LiveID] = private
			}

			showtop3 := 3
			for k, private := range privatemap {
				if showtop3 > 0 {
					body.RecommandLivingList = append(body.RecommandLivingList, private)
					showtop3--
					delete(comrecinfo, k)
					delete(privatehostmap, k)
				} else {
					break
				}
			}

			// 剩余未开播的私档
			var leftprivatehost []*RecommandLivingInfo
			for _, v := range privatehostmap {
				leftprivatehost = append(leftprivatehost, v)
			}
			sort.Sort(RecommandLivingInfoSlice(leftprivatehost))
			for _, v := range leftprivatehost {
				body.RecommandLivingList = append(body.RecommandLivingList, v)
				delete(comrecinfo, v.LiveID)
			}

			// 新开播
			newliving := latest5MinLivingRandomN(comrecinfo, 2)
			for _, v := range newliving {
				body.RecommandLivingList = append(body.RecommandLivingList, v)
				delete(comrecinfo, v.LiveID)
			}

			var normalhost []*RecommandLivingInfo
			for _, rec := range comrecinfo {
				normalhost = append(normalhost, rec)
			}
			if len(normalhost) > 0 {
				sort.Sort(RecommandLivingInfoSlice(normalhost))
				body.RecommandLivingList = append(body.RecommandLivingList, normalhost...)
			}

			key := exceptHotCacheKey + strconv.FormatInt(v.Typeid, 10)
			//appzaplog.Debug("refreshOtherCategorylist", zap.String("key", key), zap.Any("body", body))
			caches.Store(key, body)
		}
	}
	return nil
}

func specialRec(last5mrank *map[uint64]rank.ChanlBillRank, v *bilin.LivingCategoryInfo) error {
	comrecinfo, err := getCommonRecInfo(uint64(v.Typeid))
	if err != nil {
		appzaplog.Error("getCommonRecInfo failed", zap.Error(err))
		return err
	}
	body := &RecommandLivingBody{}
	stickpost, err := OntimeStickyPostList(v.Typeid)
	if err != nil {
		appzaplog.Error("OntimeStickyPostList failed", zap.Int64("typeid", v.Typeid), zap.Error(err))
		return err
	}

	resp, err := clientcenter.ConfClient().CategoryGuildRec(context.TODO(), &bilin.CategoryGuildRecReq{})
	if err != nil {
		appzaplog.Error("CategoryGuildRec failed", zap.Error(err))
		return err
	}

	for _, post := range stickpost {
		if roominfo, ok := comrecinfo[post.Roomid]; ok {
			body.RecommandLivingList = append(body.RecommandLivingList, roominfo)
			delete(comrecinfo, post.Roomid)
		}
	}

	var guildrec []*RecommandLivingInfo
	for _, v := range resp.Cateogryguildinfos {
		if roominfo, ok := comrecinfo[v.Roomid]; ok {
			if rank, ok := (*last5mrank)[roominfo.LivingHostInfo.UserID]; ok {
				roominfo.Last5MinRankValue = rank.Value
			}
			guildrec = append(guildrec, roominfo)
			delete(comrecinfo, v.Roomid)
		}
	}
	sort.Sort(Latest5MinRankSlice(guildrec))

	body.RecommandLivingList = append(body.RecommandLivingList, guildrec...)

	// 剩余用户排行
	//var normalhost []*RecommandLivingInfo
	//for _, rec := range comrecinfo {
	//	normalhost = append(normalhost, rec)
	//}
	//if len(normalhost) > 0 {
	//	sort.Sort(RecommandLivingInfoSlice(normalhost))
	//	body.RecommandLivingList = append(body.RecommandLivingList, normalhost...)
	//}
	key := exceptHotCacheKey + strconv.FormatInt(v.Typeid, 10)

	caches.Store(key, body)
	return nil
}

// 最近开播排行
type LatestLivingListSlice []*RecommandLivingInfo

func (a LatestLivingListSlice) Len() int { // 重写 Len() 方法
	return len(a)
}
func (a LatestLivingListSlice) Swap(i, j int) { // 重写 Swap() 方法
	a[i], a[j] = a[j], a[i]
}
func (a LatestLivingListSlice) Less(i, j int) bool { // 重写 Less() 方法， 从大到小排序
	return a[j].StartTime > a[i].StartTime
}

// 5分钟礼物排行
type Latest5MinRankSlice []*RecommandLivingInfo

func (a Latest5MinRankSlice) Len() int { // 重写 Len() 方法
	return len(a)
}
func (a Latest5MinRankSlice) Swap(i, j int) { // 重写 Swap() 方法
	a[i], a[j] = a[j], a[i]
}
func (a Latest5MinRankSlice) Less(i, j int) bool { // 重写 Less() 方法， 从大到小排序
	return a[j].Last5MinRankValue < a[i].Last5MinRankValue
}

// 同时在线排行
type OnlineUserSlice []*RecommandLivingInfo

func (a OnlineUserSlice) Len() int { // 重写 Len() 方法
	return len(a)
}
func (a OnlineUserSlice) Swap(i, j int) { // 重写 Swap() 方法
	a[i], a[j] = a[j], a[i]
}
func (a OnlineUserSlice) Less(i, j int) bool { // 重写 Less() 方法， 从大到小排序
	return a[j].UserCount < a[i].UserCount
}

// 开播时长排行
type LivingTimeSlice []*RecommandLivingInfo

func (a LivingTimeSlice) Len() int { // 重写 Len() 方法
	return len(a)
}
func (a LivingTimeSlice) Swap(i, j int) { // 重写 Swap() 方法
	a[i], a[j] = a[j], a[i]
}
func (a LivingTimeSlice) Less(i, j int) bool { // 重写 Less() 方法， 从大到小排序
	return a[j].StartTime > a[i].StartTime
}
