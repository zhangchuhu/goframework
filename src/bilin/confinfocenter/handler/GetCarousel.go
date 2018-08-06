package handler

import (
	"bilin/confinfocenter/dao"
	"bilin/protocol"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/kissgo/httpmetrics"
	"context"
	"sort"
	"time"
)

type CarouselSlice []*dao.Carousel

func (a CarouselSlice) Len() int { // 重写 Len() 方法
	return len(a)
}
func (a CarouselSlice) Swap(i, j int) { // 重写 Swap() 方法
	a[i], a[j] = a[j], a[i]
}
func (a CarouselSlice) Less(i, j int) bool { // 重写 Less() 方法， 从大到小排序
	return a[j].Sort < a[i].Sort
}

// 获取直播间分类信息配置
func (this *ConfInfoServantObj) GetCarousel(ctx context.Context, r *bilin.CarouselReq) (*bilin.CarouselResp, error) {
	code := int64(0)
	defer func(now time.Time) {
		httpmetrics.DefReport("GetCarousel", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())

	infos, err := dao.GetCarousel()
	if err != nil {
		code = GetCarouselFailed
		appzaplog.Error("GetCarousel failed", zap.Error(err))
		return nil, err
	}

	if infos != nil {
		sort.Sort(CarouselSlice(infos))
	}

	resp := &bilin.CarouselResp{}
	for _, v := range infos {
		resp.Carousel = append(resp.Carousel, &bilin.CarouselInfo{
			Id:           v.ID,
			BackgroudURL: v.BackgroudURL,
			TargetType:   v.TargetType,
			TargetURL:    v.TargetURL,
			StartTime:    v.StartTime,
			EndTime:      v.EndTime,
			Channel:      v.Channel,
			Version:      v.Version,
			ForUserType:  v.ForUserType,
		})
	}

	return resp, nil
}
