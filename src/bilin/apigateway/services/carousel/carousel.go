package carousel

import (
	"bilin/protocol"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	//"code.yy.com/yytars/goframework/tars"
	"bilin/clientcenter"
	"context"
	"strconv"
)

const (
	TARGET_TYPE_ANCHOR   = 19
	CAROUSEL_TYPE_BANNER = 1
	CAROUSEL_TYPE_ANCHOR = 2
)

func GetEffectiveCarousel() ([]*Carousel, error) {
	resp, err := clientcenter.ConfClient().GetCarousel(context.TODO(), &bilin.CarouselReq{})
	if err != nil {
		appzaplog.Error("GetCarousel err", zap.Error(err))
		return nil, err
	}

	appzaplog.Debug("GetCarousel success", zap.Any("carousel_resp", resp))

	carousel_list := []*Carousel{}
	if resp.Carousel == nil {
		return carousel_list, nil
	}

	//获取所有主播
	uids := []uint64{}
	for _, v := range resp.Carousel {
		if v.TargetType == TARGET_TYPE_ANCHOR {
			uid, err := strconv.ParseUint(v.TargetURL, 10, 64)
			if err != nil {
				appzaplog.Error("ParseUint of anchor id err", zap.Error(err), zap.String("TargetURL", v.TargetURL))
			} else {
				uids = append(uids, uid)
			}
		}
	}

	var room_info_resp *bilin.BatchLivingRoomsInfoByHostsResp
	//没有主播配置就不要rpc查询了
	var no_anchor bool
	if len(uids) == 0 {
		appzaplog.Info("no config anchor type")
		no_anchor = true
	} else {
		room_info_resp, err = clientcenter.RoomCenterClient().BatchLivingRoomsInfoByHosts(context.TODO(), &bilin.BatchLivingRoomsInfoByHostsReq{Hosts: uids})
		if err != nil {
			appzaplog.Error("BatchLivingRoomsInfoByHosts err", zap.Error(err))
			no_anchor = true
		} else if room_info_resp != nil && room_info_resp.Livingrooms == nil {
			appzaplog.Info("BatchLivingRoomsInfoByHosts livingrooms nil, no anchor online")
			no_anchor = true
		} else {
			appzaplog.Info("BatchLivingRoomsInfoByHosts", zap.Any("room_info_resp", room_info_resp))
			no_anchor = false
		}
	}

	for _, v := range resp.Carousel {

		if v.TargetType == TARGET_TYPE_ANCHOR { //主播
			if no_anchor {
				continue
			}
			//如果在线
			uid, _ := strconv.ParseUint(v.TargetURL, 10, 64)
			if room_info, ok := room_info_resp.Livingrooms[uid]; ok && room_info.LockStatus == 0 {
				carousel_list = append(carousel_list, &Carousel{
					CarouselId:   uint64(v.Id),
					CarouselType: CAROUSEL_TYPE_ANCHOR,
					BackgroudURL: v.BackgroudURL,
					TargetType:   uint32(v.TargetType),
					TargetURL:    v.TargetURL,
					Version:      v.Version,
					Channel:      v.Channel,
					ForUserType:  v.ForUserType,
					Room: LivingRoom{
						RoomID: room_info.Roomid,
						UserID: uid,
					},
				})
			}
		} else { //广告
			carousel_list = append(carousel_list, &Carousel{
				CarouselId:   uint64(v.Id),
				CarouselType: CAROUSEL_TYPE_BANNER,
				BackgroudURL: v.BackgroudURL,
				TargetType:   uint32(v.TargetType),
				TargetURL:    v.TargetURL,
				Version:      v.Version,
				Channel:      v.Channel,
				ForUserType:  v.ForUserType,
			})
		}

	}
	appzaplog.Debug("GetEffectiveCarousel success", zap.Any("carousel_list", carousel_list))
	return carousel_list, nil
}
