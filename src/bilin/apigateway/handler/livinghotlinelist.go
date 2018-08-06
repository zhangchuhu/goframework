package handler

import (
	"bilin/apigateway/cache"
	"bilin/apigateway/services/sensitive"
	"bilin/apigateway/services/userinfo"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/kissgo/httpmetrics"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type LivingListBody struct {
	UserId      int64 `json:"userId"`
	CategroyID  int64 `json:"categoryId"`  // 品类id
	Count       int64 `json:"count"`       // 每页大小，默认20
	Launchtimes int64 `json:"launchtimes"` // 累计启动次数，默认18
	Pagenum     int64 `json:"page"`        // 希望拉取的页码，默认为0
}

func (p *LivingListBody) parseForm(r *http.Request) error {
	var err error
	if err = r.ParseForm(); err != nil {
		return err
	}

	if UserIdstr := r.PostFormValue("userId"); UserIdstr != "" {
		if p.UserId, err = strconv.ParseInt(UserIdstr, 10, 64); err != nil {
			return err
		}
	}

	p.CategroyID = 1000 //默认是热门
	if CategroyIDstr := r.PostFormValue("categoryId"); CategroyIDstr != "" {
		if p.CategroyID, err = strconv.ParseInt(CategroyIDstr, 10, 64); err != nil {
			return err
		}
	}

	p.Count = 20
	if countstr := r.PostFormValue("count"); countstr != "" {
		if p.Count, err = strconv.ParseInt(countstr, 10, 64); err != nil {
			return err
		}
	}
	// 最多50
	if p.Count > 50 || p.Count < 0 {
		appzaplog.Warn("Count unexpected", zap.Int64("count", p.Count))
		p.Count = 50
	}

	p.Launchtimes = OLD_USER_LOGIN_TIMES + 1
	if launchtimesstr := r.PostFormValue("launchtimes"); launchtimesstr != "" {
		if p.Launchtimes, err = strconv.ParseInt(launchtimesstr, 10, 64); err != nil {
			return err
		}
	}

	if pagestr := r.PostFormValue("page"); pagestr != "" {
		if p.Pagenum, err = strconv.ParseInt(pagestr, 10, 64); err != nil {
			return err
		}
	}
	if p.Pagenum < 0 {
		appzaplog.Warn("Pagenum unexpected", zap.Int64("Pagenum", p.Pagenum))
		p.Pagenum = 0
	}

	return nil
}

func HandleGetLivingHotLineList(rw http.ResponseWriter, r *http.Request) *HttpError {
	appzaplog.Debug("[+]HandleGetLivingHotLineList enter")
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
		appzaplog.Debug("[+]HandleGetLivingHotLineList enter", zap.Any("livingboyd", livingboyd), zap.Any("reqparam", reqparam))
		body, err := cache.GetRecLivingBody(livingboyd.CategroyID, olduser(livingboyd.Launchtimes))
		if err != nil {
			appzaplog.Error("failed to GetRecLivingBody", zap.Error(err), zap.Int64("CategroyID", livingboyd.CategroyID))
			fmt.Fprintf(rw, "%s", cache.ErrNotFoundCache)
			return noCacheHttpErr
		}

		realbody := cache.RecommandLivingBody{
			LastPage: "false",
		}

		endindex := (livingboyd.Pagenum + 1) * livingboyd.Count
		startindex := livingboyd.Pagenum * livingboyd.Count
		reclen := int64(len(body.RecommandLivingList))
		if reclen <= endindex {
			realbody.LastPage = "true"
			if startindex < reclen {
				realbody.RecommandLivingList = body.RecommandLivingList[startindex:]
			}
		} else {
			realbody.RecommandLivingList = body.RecommandLivingList[startindex:endindex]
		}

		//苹果审核帐号特殊处理
		if userinfo.IsAppleCheckUser(uint64(livingboyd.UserId), reqparam.Version, reqparam.ClientType, remoteip) {
			httpmetrics.CounterMetric("AppleCheckUser", 1)
			var applereclivinglist []*cache.RecommandLivingInfo
			//审核
			for _, v := range realbody.RecommandLivingList {
				if sensitive.SensiTiveWord(v.Title) {
					appzaplog.Info("HandleGetLivingHotLineList Sensitive world detect",
						zap.String("title", v.Title),
						zap.Int64("uid", livingboyd.UserId),
					)
					continue
				}
				applereclivinglist = append(applereclivinglist, v)
			}
			realbody.RecommandLivingList = applereclivinglist
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

func olduser(launchtimes int64) bool {
	return launchtimes > OLD_USER_LOGIN_TIMES
}

func realclientip(r *http.Request) string {
	retip := ""
	forwardips := r.Header.Get("X-Forwarded-For")
	if forwardips != "" {
		if ips := strings.Split(forwardips, ","); len(ips) > 0 {
			retip = ips[0]
		}
	}
	return retip
}
