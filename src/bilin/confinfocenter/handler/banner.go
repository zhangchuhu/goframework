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

func (this *ConfInfoServantObj) Banner(ctx context.Context, r *bilin.BannerReq) (*bilin.BannerResp, error) {
	code := int64(0)
	defer func(now time.Time) {
		httpmetrics.DefReport("Banner", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())

	infos, err := dao.Banner(r.Typid)
	if err != nil {
		code = GetBannerFailed
		appzaplog.Error("Banner failed", zap.Error(err))
		return nil, err
	}

	if infos != nil {
		sort.Sort(CarouselSlice(infos))
	}

	resp := &bilin.BannerResp{}
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
			Width:        v.Width,
			Height:       v.Height,
			Position:     v.Position,
			Hotlinetype:  v.HotLineType,
		})
	}

	return resp, nil
}
