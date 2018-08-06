package handler

import (
	"bilin/apigateway/cache"
	"bilin/apigateway/services/livingbanner"
	"bilin/apigateway/services/sensitive"
	"bilin/apigateway/services/userinfo"
	"bilin/clientcenter"
	"bilin/protocol"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/kissgo/httpmetrics"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type HotLineListV2 struct {
	LastPage            string        `json:"last_page"`
	RecommandLivingList []interface{} `json:"recommand_living_list"`
}

type RecBanner struct {
	TypeID     int64  `json:"type_id"`
	BGURL      string `json:"bgurl"`       //背景图片
	Height     int32  `json:"height"`      // 图片高度
	TargetType uint32 `json:"target_type"` //类型，有直播间，主持，H5,功能模块
	TargetURL  string `json:"target_url"`  // 目标地址，H5为链接，功能模块为约定string，直播间和主播为房间id
}

func HandleGetLivingHotLineListV2(rw http.ResponseWriter, r *http.Request) *HttpError {
	if r.Method == "POST" {
		var reqparam ReqParam
		if err := reqparam.ParseURL(r); err != nil {
			appzaplog.Error("ParseURL failed", zap.Any("reqparam", reqparam), zap.Error(err))
		}

		remoteip := realclientip(r)

		var livingboyd LivingListBody
		if err := livingboyd.parseForm(r); err != nil {
			appzaplog.Error("ParseFrom failed", zap.Error(err))
			fmt.Fprintf(rw, "%s", cache.ErrBadRequest)
			return parseFormHttpErr
		}
		appzaplog.Debug("[+]HandleGetLivingHotLineListV2 enter", zap.Any("livingboyd", livingboyd), zap.Any("reqparam", reqparam))
		body, err := cache.GetRecLivingBody(livingboyd.CategroyID, olduser(livingboyd.Launchtimes))
		if err != nil {
			appzaplog.Error("failed to GetRecLivingBody", zap.Error(err), zap.Int64("CategroyID", livingboyd.CategroyID))
			fmt.Fprintf(rw, "%s", cache.ErrNotFoundCache)
			return noCacheHttpErr
		}

		realbody := &HotLineListV2{
			LastPage: "false",
		}

		//苹果审核帐号特殊处理
		if userinfo.IsAppleCheckUser(uint64(livingboyd.UserId), reqparam.Version, reqparam.ClientType, remoteip) {
			httpmetrics.CounterMetric("AppleCheckUser", 1)
			//审核
			for _, v := range body.RecommandLivingList {
				if sensitive.SensiTiveWord(v.Title) {
					appzaplog.Info("HandleGetLivingHotLineList Sensitive world detect",
						zap.String("title", v.Title),
						zap.Int64("uid", livingboyd.UserId),
					)
					continue
				}
				realbody.RecommandLivingList = append(realbody.RecommandLivingList, v)
			}
		} else {
			for _, v := range body.RecommandLivingList {
				realbody.RecommandLivingList = append(realbody.RecommandLivingList, v)
			}
		}

		// 第一页加入banner
		if livingboyd.Pagenum == 0 {
			var user_type int = 1
			if livingboyd.Launchtimes > OLD_USER_LOGIN_TIMES {
				user_type = 2
			}
			banner := livingbanner.MatchLivingBanner(reqparam.Platform, reqparam.Version, user_type, livingboyd.CategroyID)
			for _, v := range banner {
				if strings.Contains(v.TargetURL, roomIdTargetUrlFilterKey) {
					roomid := RoomId(v.TargetURL)
					if roomid <= 0 {
						appzaplog.Warn("MatchLivingBanner roomid parse illegal")
						continue
					} else if !onLivingRoom(roomid) {
						continue
					}
				}
				newelement := &RecBanner{
					TypeID:     2,
					BGURL:      v.BGURL,
					Height:     v.Height,
					TargetType: v.TargetType,
					TargetURL:  v.TargetURL,
				}
				if int(v.Position) > len(realbody.RecommandLivingList) {
					realbody.RecommandLivingList = append(realbody.RecommandLivingList, newelement)
				} else {
					// insert into
					pos := v.Position - 1
					realbody.RecommandLivingList = append(realbody.RecommandLivingList, 0)
					copy(realbody.RecommandLivingList[pos+1:], realbody.RecommandLivingList[pos:])
					realbody.RecommandLivingList[pos] = newelement
				}
			}
		}

		endindex := (livingboyd.Pagenum + 1) * livingboyd.Count
		startindex := livingboyd.Pagenum * livingboyd.Count
		reclen := int64(len(realbody.RecommandLivingList))
		if reclen <= endindex {
			realbody.LastPage = "true"
			if startindex < reclen {
				realbody.RecommandLivingList = realbody.RecommandLivingList[startindex:]
			}
		} else {
			realbody.RecommandLivingList = realbody.RecommandLivingList[startindex:endindex]
		}

		homelist := &cache.HttpRetComm{
			IsEncrypt: "false",
			Data: cache.HttpRetDataComm{
				Msg:  "success",
				Body: realbody,
			},
		}
		jsonbin, err := json.Marshal(homelist)
		if err != nil {
			appzaplog.Error("Marshal json failed", zap.Error(err))
			return jsonMarshalHttpErr
		}
		fmt.Fprintf(rw, "%s", jsonbin)
		return successHttp
	}
	return notSupportMethodHttpErr
}

const roomIdTargetUrlFilterKey = "hotlineId="

func RoomId(targetUrl string) int64 {
	if targetUrl == "" {
		return 0
	}
	info := strings.SplitAfter(targetUrl, roomIdTargetUrlFilterKey)
	if len(info) > 1 {
		roomid, err := strconv.ParseInt(info[1], 10, 64)
		if err != nil {
			appzaplog.Error("RoomId ParseInt err", zap.Error(err), zap.String("roomid", info[1]))
			return 0
		}
		return roomid
	}
	return 0
}

func onLivingRoom(roomid int64) bool {
	resp, err := clientcenter.RoomCenterClient().IsLiving(context.TODO(), &bilin.IsLivingReq{
		Roomid: roomid,
	})
	if err != nil {
		appzaplog.Error("onLivingRoom IsLiving err", zap.Error(err))
		// fake to true
		return true
	}
	return resp.Isliving
}
