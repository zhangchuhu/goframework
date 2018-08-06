package cache

import (
	"bilin/common/cacheprocessor"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var caches sync.Map

type Cacher interface {
	refresh() error
}

var (
	//roomcenterclient     bilin.RoomInfoServantClient
	//confinfocenterclient bilin.ConfInfoServantClient
	//userinfoclient       userinfocenter.UserInfoCenterObjClient
	ErrNotFoundCache []byte
	ErrBadRequest    []byte
)

const (
	hotLineListCacheKey   = "hotlinelist"    // 新用户热门缓存
	oldUserHotRecCacheKey = "olduserhotrec"  // 老用户热门缓存
	commRecInfoCacheKey   = "comrecinfolist" // 通用信息缓存
	exceptHotCacheKey     = "expecthotlist"  // 品类列表缓存

	categroyCacheKey          = "category"              // 品类分类缓存key,包括了热门
	startLivecategroyCacheKey = "startlivecategorylist" // 开播页的品类分类缓存key
)

func InitCache() error {
	ErrNotFoundCache, _ = getErrorCache()
	ErrBadRequest, _ = getBadRequestCache()

	//开启品类分类缓存
	if err := cacheprocessor.CacheProcessor("refreshCategory", 30*time.Second, refreshCategory); err != nil {
		return err
	}

	// 通用信息缓存
	if err := cacheprocessor.CacheProcessor("refreshCommonRecInfo", 5*time.Second, refreshCommonRecInfo); err != nil {
		return err
	}

	// 新热门品类列表
	if err := cacheprocessor.CacheProcessor("refreshHotList", 6*time.Second, refreshHotList); err != nil {
		return err
	}

	// 老用户热门品类列表缓存
	if err := cacheprocessor.CacheProcessor("refreshOldUserHotList", 6*time.Second, refreshOldUserHotList); err != nil {
		return err
	}

	// 其他品类列表
	if err := cacheprocessor.CacheProcessor("refreshOtherCategorylist", 8*time.Second, refreshOtherCategorylist); err != nil {
		return err
	}
	return nil
}

func GetRecLivingBody(categoryid int64, olduser bool) (*RecommandLivingBody, error) {
	return getRecLivingBodyFromCache(takeCacheKey(categoryid, olduser))
}

func getRecLivingBodyFromCache(key string) (*RecommandLivingBody, error) {
	if v, ok := caches.Load(key); ok {
		if ret, yes := v.(*RecommandLivingBody); yes {
			return ret, nil
		}
	}
	return nil, errors.New("not found")
}

func getErrorCache() ([]byte, error) {
	homelist := &HttpRetComm{
		IsEncrypt: "false",
		Data: HttpRetDataComm{
			Code: 1,
			Msg:  "查找缓存失败",
		},
	}
	return json.Marshal(homelist)
}

func getBadRequestCache() ([]byte, error) {
	homelist := &HttpRetComm{
		IsEncrypt: "false",
		Data: HttpRetDataComm{
			Code: http.StatusBadRequest,
			Msg:  http.StatusText(http.StatusBadRequest),
		},
	}
	return json.Marshal(homelist)
}

func takeCacheKey(categoryid int64, olduser bool) string {
	switch categoryid {
	case hotcategoryid:
		if olduser {
			return oldUserHotRecCacheKey
		}
		return hotLineListCacheKey
	default:
		return exceptHotCacheKey + strconv.FormatInt(categoryid, 10)
	}
}
