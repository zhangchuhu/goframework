package handler

import (
	"bilin/apigateway/services/carousel"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

const (
	OLD_USER_LOGIN_TIMES = 15
)

type QueryCarouselBody struct {
	UserId      uint64 `json:"userId"`
	Launchtimes int64  `json:"launchtimes"` // 登录次数
}

func (p *QueryCarouselBody) parseForm(r *http.Request) error {

	var err error
	if err = r.ParseForm(); err != nil {
		return err
	}

	if UserIdstr := r.PostFormValue("userId"); UserIdstr != "" {
		if p.UserId, err = strconv.ParseUint(UserIdstr, 10, 64); err != nil {
			return err
		}
	}

	p.Launchtimes = OLD_USER_LOGIN_TIMES + 1
	if launchtimesstr := r.PostFormValue("launchtimes"); launchtimesstr != "" {
		if p.Launchtimes, err = strconv.ParseInt(launchtimesstr, 10, 64); err != nil {
			return err
		}
	}
	return nil
}

func HandleQueryCarousel(rw http.ResponseWriter, r *http.Request) *HttpError {

	appzaplog.Debug("HandleQueryCarousel enter")

	if r.Method == "POST" {
		var reqparam ReqParam
		if err := reqparam.ParseURL(r); err != nil {
			appzaplog.Error("ParseURL failed", zap.Any("reqparam", reqparam), zap.Error(err))
		}

		var carousel_body QueryCarouselBody
		if err := carousel_body.parseForm(r); err != nil {
			appzaplog.Error("ParseFrom failed", zap.Error(err))
			//fmt.Fprintf(rw, "%s", cache.ErrBadRequest)
			return parseFormHttpErr
		}
		appzaplog.Debug("HandleQueryCarousel enter", zap.Any("carousel_body", carousel_body), zap.Any("reqparam", reqparam))

		var user_type int = 1
		if carousel_body.Launchtimes > OLD_USER_LOGIN_TIMES {
			user_type = 2
		}

		resp := &HttpRetComm{
			IsEncrypt: "false",
			Data: HttpRetDataComm{
				Code: 0,
				Msg:  "success",
				Body: carousel.CarouselListBody{
					CarouselList: carousel.MatchUserCarousel(reqparam.Platform, reqparam.Version, user_type),
				},
			},
		}

		byte, err := json.Marshal(resp)
		if err != nil {
			appzaplog.Error("failed to marshal", zap.Error(err))
			return jsonMarshalHttpErr
		}
		fmt.Fprintf(rw, "%s", byte)
		return successHttp
	}
	return notSupportMethodHttpErr
}
