package cache

import (
	"bilin/clientcenter"
	"bilin/protocol"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"context"
	"encoding/json"
	"errors"
)

type HttpRetDataComm struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Body interface{} `json:"body"`
}

type HttpRetComm struct {
	IsEncrypt string          `json:"isEncrypt"`
	Data      HttpRetDataComm `json:"data"`
}

//LivingCategory 直播间品类详细信息
type LivingCategory struct {
	CategoryID      int    `json:"typeId"`
	CategoryName    string `json:"typeName"`
	FontColor       string `json:"fontColor"`
	BackgroundImage string `json:"backgroundImage"`
}

type LivingCategoryBody struct {
	LivingCategoryList []LivingCategory `json:"hotlineDirectTypeList"`
}

var hotCategory = LivingCategory{
	CategoryID:      hotcategoryid,
	CategoryName:    "热门",
	FontColor:       "#FFFFFF",
	BackgroundImage: "http://img.onbilin.com/0/15209943423653008.jpg",
}

func refreshCategory() error {
	resp, err := clientcenter.ConfClient().LivingCategorys(context.TODO(), &bilin.LivingCategorysReq{})
	if err != nil {
		appzaplog.Error("LivingCategorys failed", zap.Error(err))
		return err
	}
	if err = storeCategoryCache(resp); err != nil {
		appzaplog.Error("storeCategoryCache failed", zap.Error(err))
		return err
	}
	return storeStartLiveCategoryCache(resp)
}

func storeCategoryCache(resp *bilin.LivingCategorysResp) error {
	livecategory := []LivingCategory{}
	//add 热门
	livecategory = append(livecategory, hotCategory)

	for _, v := range resp.Livingcategorys {
		if v.Typeid != hotcategoryid {
			livecategory = append(livecategory, LivingCategory{
				CategoryID:      int(v.Typeid),
				CategoryName:    v.Name,
				FontColor:       v.Color,
				BackgroundImage: v.Backgroudimage,
			})
		}
	}

	homelist := &HttpRetComm{
		IsEncrypt: "false",
		Data: HttpRetDataComm{
			Code: 0,
			Msg:  "success",
			Body: LivingCategoryBody{
				LivingCategoryList: livecategory,
			},
		},
	}
	byte, err := json.Marshal(homelist)
	if err != nil {
		appzaplog.Error("failed to marshal", zap.Error(err))
		return err
	}
	caches.Store(categroyCacheKey, byte)
	return nil
}

func storeStartLiveCategoryCache(resp *bilin.LivingCategorysResp) error {
	livecategory := []LivingCategory{}
	for _, v := range resp.Livingcategorys {
		livecategory = append(livecategory, LivingCategory{
			CategoryID:      int(v.Typeid),
			CategoryName:    v.Name,
			FontColor:       v.Color,
			BackgroundImage: v.Backgroudimage,
		})
	}

	homelist := &HttpRetComm{
		IsEncrypt: "false",
		Data: HttpRetDataComm{
			Code: 0,
			Msg:  "success",
			Body: LivingCategoryBody{
				LivingCategoryList: livecategory,
			},
		},
	}
	byte, err := json.Marshal(homelist)
	if err != nil {
		appzaplog.Error("failed to marshal", zap.Error(err))
		return err
	}
	caches.Store(startLivecategroyCacheKey, byte)
	return nil
}

func GetStaticCategoryCache() ([]byte, error) {
	return takeBinaryCache(categroyCacheKey)
}

func GetStartLiveStaticCategoryCache() ([]byte, error) {
	return takeBinaryCache(startLivecategroyCacheKey)
}

func takeBinaryCache(key string) ([]byte, error) {
	if v, ok := caches.Load(key); ok {
		if ret, yes := v.([]byte); yes {
			return ret, nil
		}
	}
	return nil, errors.New("not found")
}
