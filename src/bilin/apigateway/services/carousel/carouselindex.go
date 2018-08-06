package carousel

import (
	"bilin/apigateway/common/index"
	"bilin/common/cacheprocessor"
	"sync/atomic"

	"fmt"
	"time"
	//"encoding/json"
)

type LivingRoom struct {
	RoomID uint64 `json:"room_id"`
	UserID uint64 `json:"user_id"`
	//NickName string `json:"nick_name"`
	//Count    uint64 `json:"count"`
}

type Carousel struct {
	CarouselId   uint64     `json:"carousel_id"`
	CarouselType uint32     `json:"carousel_type"`
	BackgroudURL string     `json:"backgroud_url"`
	TargetType   uint32     `json:"target_type"`
	TargetURL    string     `json:"target_url"`
	Room         LivingRoom `json:"living_room"`
	Channel      string     `json:"-"`
	Version      string     `json:"-"`
	ForUserType  int32      `json:"-"`
}

type CarouselListBody struct {
	CarouselList []*Carousel `json:"carousel_list"`
}

const (
	CHANNEL_MATCH_ALL  = ""  //web页面配置空表示匹配所有渠道
	VERSION_MATCH_ALL  = ""  //web页面配置空表示匹配所有渠道
	FOR_USER_MATCH_ALL = "0" //web页面配置0表示匹配所有渠道
)

type CarouselIndex struct {
	Channel      index.Index //根据CarouselList生成索引
	Version      index.Index //根据CarouselList生成索引
	ForUser      index.Index //根据CarouselList生成索引
	CarouselList []*Carousel //已经配置的轮播信息列表
}

func (c *CarouselIndex) InitCarouselIndex() {
	c.Channel.InitCache()
	c.Version.InitCache()
	c.ForUser.InitCache()
	c.CarouselList = make([]*Carousel, 0)
}

var SafeCarouselIndex atomic.Value

func init() {
	if err := cacheprocessor.CacheProcessor("LoadCarouselIndex", 3*time.Second, LoadCarouselIndex); err != nil {
		panic("LoadCarouselIndex failed")
	}
}

func GetCarouselIndex() *CarouselIndex {
	return SafeCarouselIndex.Load().(*CarouselIndex)
}

func LoadCarouselIndex() error {

	var carousel_index CarouselIndex

	carousel_index.InitCarouselIndex()

	carous_list, err := GetEffectiveCarousel()
	if err != nil {
		return err
	}

	if carous_list == nil {
		carousel_index.CarouselList = make([]*Carousel, 0)
	} else {
		carousel_index.CarouselList = carous_list
	}

	for _, v := range carousel_index.CarouselList {
		carousel_index.Channel.Add(v.CarouselId, v.Channel, CHANNEL_MATCH_ALL)
		carousel_index.Version.Add(v.CarouselId, v.Version, VERSION_MATCH_ALL)
		for_user := fmt.Sprintf("%d", v.ForUserType)
		carousel_index.ForUser.Add(v.CarouselId, for_user, FOR_USER_MATCH_ALL)
	}

	SafeCarouselIndex.Store(&carousel_index)
	return nil
}

//根据用户的渠道 版本 类型, 匹配不同的信息返回
func MatchUserCarousel(channel, version string, user_type int) []*Carousel {

	carousel_index := GetCarouselIndex()
	list := make([]*Carousel, 0, len(carousel_index.CarouselList))

	for _, v := range carousel_index.CarouselList {
		if ok := carousel_index.Channel.Match(v.CarouselId, channel); !ok {
			continue
		}

		if ok := carousel_index.Version.Match(v.CarouselId, version); !ok {
			continue
		}

		for_user := fmt.Sprintf("%d", user_type)
		if ok := carousel_index.ForUser.Match(v.CarouselId, for_user); !ok {
			continue
		}

		list = append(list, v)
	}

	return list
}

func GetAllCarousel() []*Carousel {
	carousel_index := GetCarouselIndex()
	return carousel_index.CarouselList
}
