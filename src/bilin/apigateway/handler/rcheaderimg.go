package handler

import (
	"bilin/apigateway/cache"
	"bilin/apigateway/services/avatarlist"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type RCHeaderImgBody struct {
	SmallImgUrls []string `json:"small_img_urls"`
}

type QueryHeaderImgBody struct {
	UserId uint64 `json:"userId"`
	Sex    int64  `json:"sex"`
}

func (p *QueryHeaderImgBody) parseForm(r *http.Request) error {
	var err error
	if err = r.ParseForm(); err != nil {
		return err
	}

	if UserIdstr := r.PostFormValue("userId"); UserIdstr != "" {
		if p.UserId, err = strconv.ParseUint(UserIdstr, 10, 64); err != nil {
			return err
		}
	}

	Sexstr := r.PostFormValue("sex")
	appzaplog.Debug("HandleRCHeaderImg enter", zap.String("Sexstr", Sexstr))
	if Sexstr != "" {
		if p.Sex, err = strconv.ParseInt(Sexstr, 10, 64); err != nil {
			return err
		}
	} else {
		p.Sex = 1
	}

	return nil
}

func HandleRCHeaderImg(rw http.ResponseWriter, r *http.Request) *HttpError {

	var head_img_body QueryHeaderImgBody
	if err := head_img_body.parseForm(r); err != nil {
		appzaplog.Error("ParseFrom failed", zap.Error(err))
		//fmt.Fprintf(rw, "%s", cache.ErrBadRequest)
		return parseFormHttpErr
	}
	appzaplog.Debug("HandleRCHeaderImg enter", zap.Any("head_img_body", head_img_body))

	ret := cache.HttpRetComm{
		IsEncrypt: "false",
		Data: cache.HttpRetDataComm{
			Msg: "success",
		},
	}
	ret.Data.Body = &RCHeaderImgBody{
		SmallImgUrls: avatarlist.GetAvatarList(head_img_body.Sex),
	}
	byte, err := json.Marshal(ret)
	if err != nil {
		appzaplog.Error("json marshal failed", zap.Error(err))
		return jsonMarshalHttpErr
	}
	fmt.Fprintf(rw, "%s", byte)
	return successHttp
}
