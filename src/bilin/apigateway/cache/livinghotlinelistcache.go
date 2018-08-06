package cache

import (
	"bilin/apigateway/config"
	"bilin/apigateway/services/rank"
	"bilin/clientcenter"
	"bilin/protocol"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"context"
	"sort"
	"strings"
	"time"
)

type LivingHostInfo struct {
	UserID   uint64   `json:"user_id"`
	NickName string   `json:"nick_name"`
	City     string   `json:"city"`
	SmallURL string   `json:"small_url"`
	TagUrl   []string `json:"tag_url,omitempty"`
}

type RecommandLivingInfo struct {
	TypeID            int64          `json:"type_id"`
	LiveID            uint64         `json:"live_id"`
	Title             string         `json:"title"`
	UserCount         uint64         `json:"user_count"`
	LivingHostInfo    LivingHostInfo `json:"living_host_info"`
	SortWeight        int64          `json:"-"` // 排序
	CategoryId        uint64         `json:"category_id"`
	StartTime         uint64         `json:"_"` // 排序
	Last5MinRankValue int64          `json:"_"`
}

type RecommandLivingBody struct {
	LastPage            string                 `json:"last_page"`
	RecommandLivingList []*RecommandLivingInfo `json:"recommand_living_list"`
}

type RecommandLivingInfoSlice []*RecommandLivingInfo

func (a RecommandLivingInfoSlice) Len() int { // 重写 Len() 方法
	return len(a)
}
func (a RecommandLivingInfoSlice) Swap(i, j int) { // 重写 Swap() 方法
	a[i], a[j] = a[j], a[i]
}
func (a RecommandLivingInfoSlice) Less(i, j int) bool { // 重写 Less() 方法， 从大到小排序
	return a[j].SortWeight < a[i].SortWeight
}

// 当前的运营置顶区
func OntimeStickyPostList(typeid int64) (map[uint64]*bilin.CategoryStickieInfo, error) {
	resp, err := clientcenter.ConfClient().CategoryStickie(context.TODO(), &bilin.CategoryStickieReq{})
	if err != nil {
		appzaplog.Error("CategoryStickie failed", zap.Error(err))
		return nil, err
	}
	now := time.Now().Unix()
	for k, v := range resp.Categoryinfo {
		//删除未开始和已经结束的
		if v.Starttime > now || v.Endtime < now || v.Typeid != typeid {
			delete(resp.Categoryinfo, k)
		}
	}
	return resp.Categoryinfo, nil
}

func BadgeInfo() ([]*bilin.UserBabgeInfo, error) {
	resp, err := clientcenter.ConfClient().BatchUserBabge(context.TODO(), &bilin.UserBabgeReq{})
	if err != nil {
		appzaplog.Error("CategoryStickie failed", zap.Error(err))
		return nil, err
	}
	return resp.Userbabgeinfo, nil
}

const (
	onlinecountweight   = 12                     // 直播间人数
	livingtimeweight    = 3                      //开播时长,相当每4分钟顶一个人进频道
	maxlivingtimeweight = 15 * onlinecountweight // 时间最多增加15个人的排位

	//品类id
	friendcatogryid       = 5    // 交友
	musiccatogryid        = 2    // 音乐
	emotioncatogryid      = 1    // 情感
	readpiaid             = 3    // 读文PIA戏
	multiplayercategoryid = 6    //多人语玩
	hotcategoryid         = 1000 //热门
)

func getRec(info []*RecommandLivingInfo) (map[uint64]*RecommandLivingInfo, map[uint64]*RecommandLivingInfo) {
	var (
		first_rec  = make(map[uint64]*RecommandLivingInfo)
		second_rec = make(map[uint64]*RecommandLivingInfo)
	)
	// 取在线前三
	sort.Sort(OnlineUserSlice(info))
	onlinetop3 := topn(info, 3)
	for _, v := range onlinetop3 {
		first_rec[v.LiveID] = v
		second_rec[v.LiveID] = v
	}

	//取最近5分钟礼物前2
	sort.Sort(Latest5MinRankSlice(info))
	proptop2 := topn(info, 2)
	for _, v := range proptop2 {
		if v.Last5MinRankValue > 0 {
			first_rec[v.LiveID] = v
		}
	}

	//取开播时长前三
	sort.Sort(LivingTimeSlice(info))
	livingtimetop3 := topn(info, 3)
	for _, v := range livingtimetop3 {
		second_rec[v.LiveID] = v
	}
	appzaplog.Debug("getRec", zap.Any("first_rec", first_rec), zap.Any("second_rec", second_rec))
	return first_rec, second_rec
}

// 最近5分钟新开播的随机N个
func latest5MinLivingRandomN(comrecinfo map[uint64]*RecommandLivingInfo, n int) []*RecommandLivingInfo {
	var (
		ret []*RecommandLivingInfo
		now = uint64(time.Now().Unix())
	)
	for _, rec := range comrecinfo {
		if n <= 0 {
			break
		}
		if rec.StartTime+5*60 > now {
			ret = append(ret, rec)
			n--
		}
	}
	return ret
}

func guildRec(comrecinfo map[uint64]*RecommandLivingInfo) ([]*RecommandLivingInfo, error) {
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
			if rankinfo, rankok := (*last5mrank)[roominfo.LivingHostInfo.UserID]; rankok {
				roominfo.Last5MinRankValue = rankinfo.Value
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
	first_multiplayer_rec := Latest5MinRankTopNMap(multiplayerrec, 3)
	for _, v := range first_multiplayer_rec {
		retrec = append(retrec, v)
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

func topn(in []*RecommandLivingInfo, num int) []*RecommandLivingInfo {
	var ret []*RecommandLivingInfo
	top2 := min(num, len(in))
	if top2 > 0 {
		ret = append(ret, in[0:top2]...)
	}
	return ret
}
func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func refreshHotList() error {

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
	retrec, err := guildRec(comrecinfo)
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
	caches.Store(hotLineListCacheKey, body)
	return nil
}

func takelivingtime(now, start int64) int64 {
	if now > start {
		return now - start
	}
	return 0
}

func calcSortWeight(livingtime int64, onlinenum int64) int64 {
	wLivingTime := int64(livingtimeweight)
	wMaxLivingTime := int64(maxlivingtimeweight)
	wOnLineCount := int64(onlinecountweight)
	if conf := config.GetAppConfig(); conf != nil {
		wLivingTime = conf.LivingTimeWeight
		wMaxLivingTime = conf.MaxLivingTimeWeight
		wOnLineCount = conf.OnLineCountWeight
	}
	return timeDimenWeight(livingtime, wLivingTime, wMaxLivingTime) + wOnLineCount*onlinenum
}

func timeDimenWeight(livingtime, weight, maxweight int64) int64 {
	calweight := (livingtime / 60) * weight
	if calweight > maxweight {
		return maxweight
	}
	return calweight
}

func tosmall(url string) string {
	pos := strings.LastIndex(url, "-")
	if pos >= 0 {
		return url[0:pos] + "-small"
	}
	return url
}

//func tobig(url string) string {
//	pos := strings.LastIndex(url, "-")
//	if pos >= 0 {
//		return url[0:pos] + "-big"
//	}
//	return url
//}
