package cache

import (
	"bilin/apigateway/services/rank"
	"bilin/clientcenter"
	"bilin/protocol"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"context"
	"sort"
)

func refreshOldUserHotList() error {

	body := &RecommandLivingBody{
		LastPage: "false",
	}

	comrecinfo, err := getCommonRecInfo(hotcategoryid)
	if err != nil {
		appzaplog.Error("getCommonRecInfo failed", zap.Error(err))
		return err
	}

	// 运营置顶区
	stickyinfo, err := OntimeStickyPostList(hotcategoryid)
	if err != nil {
		appzaplog.Error("OntimeStickyPostList refresh failed", zap.Error(err))
		//推荐区出错不影响返回，
	} else {
		if stickyinfo != nil && comrecinfo != nil {
			for roomid, _ := range stickyinfo {
				if roominfo, ok := comrecinfo[roomid]; ok {
					body.RecommandLivingList = append(body.RecommandLivingList, roominfo)
					// 从全局列表去掉
					delete(comrecinfo, roomid)
				}
			}
		}
	}

	//第一品类推荐区，多人语玩，第二品类推荐，多人语玩
	retrec, err := guildRec4OldUser(comrecinfo)
	if err != nil {
		appzaplog.Error("guildRec failed", zap.Error(err))
		return err
	}
	if retrec != nil {
		body.RecommandLivingList = append(body.RecommandLivingList, retrec...)
		for _, v := range retrec {
			delete(comrecinfo, v.LiveID)
		}
	}

	// 普通主播
	var normalhost []*RecommandLivingInfo
	for _, v := range comrecinfo {
		normalhost = append(normalhost, v)
	}
	sort.Sort(RecommandLivingInfoSlice(normalhost))
	body.RecommandLivingList = append(body.RecommandLivingList, normalhost...)
	caches.Store(oldUserHotRecCacheKey, body)
	return nil
}

func guildRec4OldUser(comrecinfo map[uint64]*RecommandLivingInfo) ([]*RecommandLivingInfo, error) {
	resp, err := clientcenter.ConfClient().CategoryGuildRec(context.TODO(), &bilin.CategoryGuildRecReq{})
	if err != nil {
		appzaplog.Error("CategoryGuildRec failed", zap.Error(err))
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}
	var (
		friendrec, musicrec, emotionrec, multiplayerrec, othersrec []*RecommandLivingInfo
		retrec                                                     []*RecommandLivingInfo
		leftrec                                                    map[uint64]*RecommandLivingInfo = make(map[uint64]*RecommandLivingInfo)
	)

	last5mrank := rank.GetChannelBillRank()
	for _, v := range resp.Cateogryguildinfos {
		if roominfo, ok := comrecinfo[v.Roomid]; ok {
			if rank, rankok := (*last5mrank)[roominfo.LivingHostInfo.UserID]; rankok {
				roominfo.Last5MinRankValue = rank.Value
			}
			leftrec[v.Roomid] = roominfo
			switch v.Typeid {
			case friendcatogryid:
				friendrec = append(friendrec, roominfo)
			case musiccatogryid:
				musicrec = append(musicrec, roominfo)
			case emotioncatogryid, readpiaid:
				emotionrec = append(emotionrec, roominfo)
			case multiplayercategoryid:
				multiplayerrec = append(multiplayerrec, roominfo)
			}
		}
	}

	// 交友
	first_friend, second_friend := getRec(friendrec)

	// 音乐
	first_music, second_music := getRec(musicrec)

	// 情感
	first_emotion, second_emotion := getRec(emotionrec)

	//第一推荐区
	for _, v := range first_friend {
		retrec = append(retrec, v)
		delete(second_friend, v.LiveID)
		break
	}

	for _, v := range first_music {
		retrec = append(retrec, v)
		delete(second_music, v.LiveID)
		break
	}

	for _, v := range first_emotion {
		retrec = append(retrec, v)
		delete(second_emotion, v.LiveID)
		break
	}

	//多人语玩,可能有重复
	top3_multiplay_rec := Latest5MinRankTopNMap(multiplayerrec, 3)
	for k, v := range top3_multiplay_rec {
		retrec = append(retrec, v)
		delete(top3_multiplay_rec, k)
		break
	}

	//第二推荐区
	for _, v := range second_friend {
		retrec = append(retrec, v)
		break
	}
	for _, v := range second_music {
		retrec = append(retrec, v)
		break
	}
	for _, v := range second_emotion {
		retrec = append(retrec, v)
		break
	}

	//多人语玩二
	for _, v := range top3_multiplay_rec {
		retrec = append(retrec, v)
		break
	}

	//去掉已经推荐的频道
	for _, v := range retrec {
		delete(comrecinfo, v.LiveID)
	}
	// 新开播
	newliving := latest5MinLivingRandomN(comrecinfo, 2)
	retrec = append(retrec, newliving...)

	//剩余栏目展示
	for _, v := range retrec {
		delete(leftrec, v.LiveID)
	}
	for _, v := range leftrec {
		othersrec = append(othersrec, v)
	}
	sort.Sort(RecommandLivingInfoSlice(othersrec))
	retrec = append(retrec, othersrec...)

	return retrec, nil
}

func Latest5MinRankTopNMap(rec []*RecommandLivingInfo, n int) map[uint64]*RecommandLivingInfo {
	ret := make(map[uint64]*RecommandLivingInfo)
	sort.Sort(Latest5MinRankSlice(rec))
	topn := topn(rec, n)
	for _, v := range topn {
		if v.Last5MinRankValue > 0 {
			ret[v.LiveID] = v
		}
	}
	return ret
}
