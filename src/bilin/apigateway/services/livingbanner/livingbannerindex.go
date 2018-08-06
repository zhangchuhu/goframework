package livingbanner

import (
	"bilin/apigateway/common/index"
	"bilin/common/cacheprocessor"
	"sync/atomic"

	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

type LivingBanner struct {
	BannerID    uint64
	BGURL       string `json:"bgurl"`       //背景图片
	TargetType  uint32 `json:"target_type"` //类型，有直播间，主持，H5,功能模块
	TargetURL   string `json:"target_url"`  // 目标地址，H5为链接，功能模块为约定string，直播间和主播为房间id
	Channel     string `json:"-"`
	Version     string `json:"-"`
	ForUserType int32  `json:"-"`
	Height      int32
	Position    int32
	HotLineType string
}

const (
	CHANNEL_MATCH_ALL  = ""  //web页面配置空表示匹配所有渠道
	VERSION_MATCH_ALL  = ""  //web页面配置空表示匹配所有渠道
	FOR_USER_MATCH_ALL = "0" //web页面配置0表示匹配所有渠道
)

type LivingBannerIndex struct {
	Channel          index.Index     //根据CarouselList生成索引
	Version          index.Index     //根据CarouselList生成索引
	ForUser          index.Index     //根据CarouselList生成索引
	LivingBannerList []*LivingBanner //已经配置的轮播信息列表
}

func (c *LivingBannerIndex) InitLivingBannerIndex() {
	c.Channel.InitCache()
	c.Version.InitCache()
	c.ForUser.InitCache()
	c.LivingBannerList = make([]*LivingBanner, 0)
}

var SafeLivingBannerIndex atomic.Value

func init() {
	if err := cacheprocessor.CacheProcessor("LoadLivingBannerIndex", 60*time.Second, LoadLivingBannerIndex); err != nil {
		panic("LoadLivingBannerIndex failed")
	}
}

func GetCarouselIndex() *LivingBannerIndex {
	return SafeLivingBannerIndex.Load().(*LivingBannerIndex)
}

func LoadLivingBannerIndex() error {

	var carousel_index LivingBannerIndex

	carousel_index.InitLivingBannerIndex()

	carous_list, err := GetEffectiveLivingBanner()
	if err != nil {
		return err
	}

	if carous_list == nil {
		carousel_index.LivingBannerList = make([]*LivingBanner, 0)
	} else {
		carousel_index.LivingBannerList = carous_list
	}

	for _, v := range carousel_index.LivingBannerList {
		carousel_index.Channel.Add(v.BannerID, v.Channel, CHANNEL_MATCH_ALL)
		carousel_index.Version.Add(v.BannerID, v.Version, VERSION_MATCH_ALL)
		for_user := fmt.Sprintf("%d", v.ForUserType)
		carousel_index.ForUser.Add(v.BannerID, for_user, FOR_USER_MATCH_ALL)
	}

	SafeLivingBannerIndex.Store(&carousel_index)
	return nil
}

//根据用户的渠道 版本 类型, 匹配不同的信息返回
func MatchLivingBanner(channel, version string, user_type int, categroyid int64) []*LivingBanner {

	carousel_index := GetCarouselIndex()
	list := make([]*LivingBanner, 0, len(carousel_index.LivingBannerList))

	for _, v := range carousel_index.LivingBannerList {
		if !haveId(v.HotLineType, categroyid) {
			continue
		}
		if ok := carousel_index.Channel.Match(v.BannerID, channel); !ok {
			continue
		}

		if ok := carousel_index.Version.Match(v.BannerID, version); !ok {
			continue
		}

		for_user := fmt.Sprintf("%d", user_type)
		if ok := carousel_index.ForUser.Match(v.BannerID, for_user); !ok {
			continue
		}

		list = append(list, v)
	}
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].Position < list[j].Position
	})
	return list
}

func haveId(hotLineType string, categoryId int64) bool {
	// 配置全部的时候是个空值
	if hotLineType == "" {
		return true
	}
	categoryIds := strings.Split(hotLineType, ",")
	cidStr := strconv.FormatInt(categoryId, 10)
	for _, v := range categoryIds {
		if v == cidStr {
			return true
		}
	}
	return false
}
