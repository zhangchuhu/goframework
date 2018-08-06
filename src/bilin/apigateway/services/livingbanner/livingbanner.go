package livingbanner

import (
	"bilin/clientcenter"
	"bilin/protocol"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"context"
	"fmt"
	"strconv"
)

const (
	livingbannertypeid = 7
	TARGET_TYPE_ANCHOR = 19
)

func GetEffectiveLivingBanner() ([]*LivingBanner, error) {
	resp, err := clientcenter.ConfClient().Banner(context.TODO(), &bilin.BannerReq{
		Typid: livingbannertypeid,
	})
	if err != nil {
		appzaplog.Error("GetEffectiveLivingBanner err", zap.Error(err))
		return nil, err
	}

	appzaplog.Debug("GetEffectiveLivingBanner success", zap.Any("carousel_resp", resp))

	livingbanner_list := []*LivingBanner{}
	if resp == nil || len(resp.Carousel) == 0 {
		return livingbanner_list, nil
	}

	//获取所有主播
	uids := []uint64{}
	for _, v := range resp.Carousel {
		if v.TargetType == TARGET_TYPE_ANCHOR {
			uid, err := strconv.ParseUint(v.TargetURL, 10, 64)
			if err != nil {
				appzaplog.Error("GetEffectiveLivingBanner ParseUint of anchor id err", zap.Error(err), zap.String("TargetURL", v.TargetURL))
			} else {
				uids = append(uids, uid)
			}
		}
	}

	var (
		roominfo *bilin.BatchLivingRoomsInfoByHostsResp
	)
	if len(uids) > 0 {
		roominfo, err = clientcenter.RoomCenterClient().BatchLivingRoomsInfoByHosts(context.TODO(), &bilin.BatchLivingRoomsInfoByHostsReq{Hosts: uids})
		if err != nil {
			appzaplog.Error("BatchLivingRoomsInfoByHosts err", zap.Error(err), zap.Uint64s("uids", uids))
		}
	}

	for _, v := range resp.Carousel {
		switch v.TargetType {
		case TARGET_TYPE_ANCHOR:
			if roominfo == nil || roominfo.Livingrooms == nil {
				continue
			}
			uid, err := strconv.ParseUint(v.TargetURL, 10, 64)
			if err != nil {
				appzaplog.Error("ParseUint err", zap.Error(err))
				continue
			}
			if info, ok := roominfo.Livingrooms[uid]; !ok || info.LockStatus != 0 {
				continue
			} else {
				v.TargetURL = fmt.Sprintf("inbilin://live/hotline?hotlineId=%d", info.Roomid)
			}
		}
		livingbanner_list = append(livingbanner_list, &LivingBanner{
			BannerID:    uint64(v.Id),
			BGURL:       v.BackgroudURL,
			TargetType:  uint32(v.TargetType),
			TargetURL:   v.TargetURL,
			Version:     v.Version,
			Channel:     v.Channel,
			ForUserType: v.ForUserType,
			Height:      v.Height,
			Position:    v.Position,
			HotLineType: v.Hotlinetype,
		})
	}
	appzaplog.Debug("GetEffectiveCarousel success", zap.Any("livingbanner_list", livingbanner_list))
	return livingbanner_list, nil
}
