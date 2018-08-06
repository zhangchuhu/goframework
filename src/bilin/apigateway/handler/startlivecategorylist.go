package handler

import (
	"bilin/apigateway/cache"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"fmt"
	"net/http"
)

func HandleStartLiveCategoryList(rw http.ResponseWriter, r *http.Request) *HttpError {
	if r.Method == "POST" {
		ret := successHttp
		byte, err := cache.GetStartLiveStaticCategoryCache()
		if err != nil {
			appzaplog.Error("GetStaticCategoryCache for Category failed", zap.Error(err))
			byte = cache.ErrNotFoundCache
			ret = noCacheHttpErr
		}
		fmt.Fprintf(rw, "%s", byte)
		return ret
	}
	return notSupportMethodHttpErr
}
