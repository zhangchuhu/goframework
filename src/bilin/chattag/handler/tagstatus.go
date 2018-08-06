package handler

import (
	"bilin/chattag/cache"
	"bilin/clientcenter"
	"bilin/protocol"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"context"
	"github.com/gin-gonic/gin"
	"strconv"
)

const (
	TagStatusOpen = iota
	TagStatusEditable
	TagStatusClosed
)

type TagStatus struct {
	Status  int64     `json:"status"`
	SetTags []ChatTag `json:"set_tags,omitempty"` //请求的userid已打的标签信息
}

type TagStatusParam struct {
	ToUserId   int64
	TalkSecond int64
}

func (p *TagStatusParam) Unmarshal(c *gin.Context) error {
	var (
		err error
	)
	if userIdStr, ok := c.GetPostForm("to_userid"); ok {
		if p.ToUserId, err = strconv.ParseInt(userIdStr, 10, 64); err != nil {
			appzaplog.Error("TagStatusParam.Unmarshal ParseInt err", zap.Error(err))
			return err
		}
	}

	if userIdStr, ok := c.GetPostForm("talk_second"); ok {
		if p.TalkSecond, err = strconv.ParseInt(userIdStr, 10, 64); err != nil {
			appzaplog.Error("TagStatusParam.Unmarshal ParseInt err", zap.Error(err))
			return err
		}
	}
	appzaplog.Debug("[-]TagStatusParam.Unmarshal", zap.Any("resp", p))
	return nil
}

func HandleTagStatus(c *gin.Context) *HttpError {
	appzaplog.Debug("[+]HandleTagStatus")
	var (
		param TagStatusParam
		body  *TagStatus
		ret   = successHttp
	)
	for {
		if err := param.Unmarshal(c); err != nil {
			appzaplog.Error("HandleTagStatus param.Unmarshal err", zap.Error(err))
			ret = parseFormHttpErr
			break
		}

		record, err := userTagRecord(param.ToUserId, c.GetInt64(uidContextKey))
		if err != nil {
			appzaplog.Error("HandleTagStatus userTagRecord err", zap.Error(err))
			ret = daoGetHttpErr
			break
		}
		body = &TagStatus{
			Status:  TagStatusClosed,
			SetTags: make([]ChatTag, 0),
		}
		if record == nil {
			body.Status = calStatus(param.ToUserId, c.GetInt64(uidContextKey), param.TalkSecond)
			if body.Status == TagStatusOpen {
				//created here
				if _, err := clientcenter.ChatTagClient().CUserChatTag(context.TODO(), &bilin.CUserChatTagReq{
					Info: &bilin.UserChatTag{
						Fromuserid: c.GetInt64(uidContextKey),
						Touserid:   param.ToUserId,
						Tagstatus:  TagStatusOpen,
						Talksecond: param.TalkSecond,
					},
				}); err != nil {
					appzaplog.Error("GetChatTagsList UserChatTag Create err", zap.Error(err), zap.Any("record", record))
					ret = daoPutHttpErr
					break
				}
			}
		} else {
			body.Status = record.Tagstatus
		}

		switch body.Status {
		case TagStatusEditable, TagStatusClosed:
			if record != nil {
				taginfo, err := FillChatTag(convertTagId(record.Chattags))
				if err != nil {
					appzaplog.Error("HandleTagStatus FillChatTag err", zap.Error(err))
					break
				}
				body.SetTags = append(body.SetTags, taginfo...)
			}
		}
		break
	}

	httpret := HttpRetComm{
		IsEncrypt: "false",
		Data: HttpRetDataComm{
			Code: ret.code,
			Msg:  ret.desc,
			Body: body,
		},
	}
	c.JSON(200, httpret)
	appzaplog.Debug("[-]HandleTagStatus", zap.Any("resp", httpret))
	return ret
}

func FillChatTag(tagids []int64) ([]ChatTag, error) {
	ret := make([]ChatTag, 0)
	if len(tagids) == 0 {
		return ret, nil
	}
	cachetagsmap := cache.TakeChatTagCache()
	if cachetagsmap == nil {
		appzaplog.Warn("TagList cache.TakeChatTagCache empty")
		return ret, nil
	}

	for _, v := range tagids {
		if taginfo, ok := cachetagsmap[v]; ok && taginfo != nil {
			ret = append(ret, ChatTag{
				TagId:   v,
				TagName: taginfo.TagName,
				Color:   taginfo.TagColor,
			})
		} else {
			appzaplog.Warn("chattag not exist or null")
		}
	}
	return ret, nil
}
