package handler

import (
	"bilin/clientcenter"
	"bilin/protocol"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"context"
	"github.com/gin-gonic/gin"
	"strconv"
)

type SetChatTagParam struct {
	ToUserId   int64
	TagIDs     string
	TalkSecond int64
}

func (p *SetChatTagParam) Unmarshal(c *gin.Context) error {
	var (
		err error
	)
	if userIdStr, ok := c.GetPostForm("to_userid"); ok {
		if p.ToUserId, err = strconv.ParseInt(userIdStr, 10, 64); err != nil {
			appzaplog.Error("SetChatTagParam.Unmarshal ParseInt err", zap.Error(err))
			return err
		}
	}
	if tagIds, ok := c.GetPostForm("tag_ids"); ok {
		p.TagIDs = tagIds
	}

	if userIdStr, ok := c.GetPostForm("talk_second"); ok {
		if p.TalkSecond, err = strconv.ParseInt(userIdStr, 10, 64); err != nil {
			appzaplog.Error("TagStatusParam.Unmarshal ParseInt err", zap.Error(err))
			return err
		}
	}
	appzaplog.Debug("[-]SetChatTagParam.Unmarshal", zap.Any("resp", p))
	return nil
}

func calStatus(toUserId, fromUserId, talkSecond int64) int64 {
	appzaplog.Debug("[+]calStatus",
		zap.Int64("touid", toUserId),
		zap.Int64("fromuid", fromUserId),
		zap.Int64("talkSecond", talkSecond))
	tagStatus := TagStatusClosed
	switch {
	case talkSecond < tagOverTalkSecond:
		tagStatus = TagStatusClosed
	default:
		resp, err := clientcenter.TakeUserInfo([]uint64{uint64(toUserId), uint64(fromUserId)})
		if err != nil {
			break
		}
		var (
			fromsex = uint32(1)
			tosex   = uint32(0)
		)
		if fromuinfo, ok := resp[uint64(fromUserId)]; ok {
			fromsex = fromuinfo.Showsex
		}
		if touinfo, ok := resp[uint64(toUserId)]; ok {
			tosex = touinfo.Showsex
		}
		if fromsex == 0 && tosex == 1 {
			tagStatus = TagStatusOpen
		}
	}
	appzaplog.Debug("[+]calStatus",
		zap.Int64("touid", toUserId),
		zap.Int64("fromuid", fromUserId),
		zap.Int64("talkSecond", talkSecond),
		zap.Int("tagStatus", tagStatus))
	return int64(tagStatus)
}

func SetChatTag(c *gin.Context) *HttpError {
	appzaplog.Debug("[+]SetChatTag")
	var (
		param SetChatTagParam
		ret   = successHttp
		body  *ChatTagList
		tags  string
	)
	for {
		if err := param.Unmarshal(c); err != nil {
			appzaplog.Error("GetChatTagsList param.Unmarshal err", zap.Error(err))
			ret = parseFormHttpErr
			break
		}

		record, err := userTagRecord(param.ToUserId, c.GetInt64(uidContextKey))
		if err != nil {
			appzaplog.Error("GetChatTagsList userTagRecord err", zap.Error(err))
			ret = daoGetHttpErr
			break
		}
		tagStatus := int64(TagStatusClosed)
		if record == nil {
			//caculate status
			tagStatus = calStatus(param.ToUserId, c.GetInt64(uidContextKey), param.TalkSecond)
		} else {
			tagStatus = record.Tagstatus
		}

		switch tagStatus {
		case TagStatusOpen:
			if record == nil {
				//should not be here
				appzaplog.Warn("should not be here")
				_, err = clientcenter.ChatTagClient().CUserChatTag(context.TODO(), &bilin.CUserChatTagReq{
					Info: &bilin.UserChatTag{
						Fromuserid:  c.GetInt64(uidContextKey),
						Touserid:    param.ToUserId,
						Chattags:    param.TagIDs,
						Updatetimes: 1,
						Tagstatus:   TagStatusEditable,
						Talksecond:  param.TalkSecond,
					},
				})
				if err != nil {
					appzaplog.Error("GetChatTagsList CUserChatTag Create err", zap.Error(err), zap.Any("record", record))
					ret = daoPutHttpErr
					break
				}
			} else { //后续补充设值的
				_, err := clientcenter.ChatTagClient().UUserChatTag(context.TODO(), &bilin.UUserChatTagReq{
					Info: &bilin.UserChatTag{
						Id:          int64(record.Id),
						Touserid:    record.Touserid,
						Updatetimes: 1,
						Chattags:    param.TagIDs,
						Tagstatus:   TagStatusEditable,
						Talksecond:  Max(param.TalkSecond, record.Talksecond),
					},
				})
				if err != nil {
					appzaplog.Error("GetChatTagsList UserChatTag Create err", zap.Error(err), zap.Any("record", record))
					ret = daoUpdateHttpErr
					break
				}
			}
			tags = param.TagIDs
		case TagStatusEditable:
			req := &bilin.UUserChatTagReq{
				Info: &bilin.UserChatTag{
					Id:          int64(record.Id),
					Touserid:    param.ToUserId,
					Chattags:    param.TagIDs,
					Updatetimes: record.Updatetimes + 1,
				},
			}
			if req.Info.Updatetimes >= editAbleTimes {
				req.Info.Tagstatus = TagStatusClosed
			}
			if _, err = clientcenter.ChatTagClient().UUserChatTag(context.TODO(), req); err != nil {
				appzaplog.Error("GetChatTagsList UserChatTag Create err", zap.Error(err), zap.Any("record", record))
				ret = daoUpdateHttpErr
				break
			}
			tags = param.TagIDs
		default:
			if record != nil {
				tags = record.Chattags
			}
			ret = TagNotAllowedErr
		}

		body = &ChatTagList{
			ChatTags: make([]ChatTag, 0),
		}
		if tags != "" {
			chattags, err := FillChatTag(convertTagId(tags))
			if err != nil {
				appzaplog.Error("FillChatTag err", zap.Error(err))
				break
			}
			body.ChatTags = append(body.ChatTags, chattags...)
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
	appzaplog.Debug("[-]SetChatTag", zap.Any("resp", httpret))
	return ret
}

const tagOverTalkSecond = 60 // 超出这个时间可以打标签

func Max(first, second int64) int64 {
	if first > second {
		return first
	}
	return second
}

const editAbleTimes = 2

func userTagRecord(toUserId, fromUserId int64) (*bilin.UserChatTag, error) {
	resp, err := clientcenter.ChatTagClient().RUserChatTag(context.TODO(), &bilin.RUserChatTagReq{
		Info: &bilin.UserChatTag{
			Fromuserid: fromUserId,
			Touserid:   toUserId,
		},
	})
	if err != nil {
		appzaplog.Error("userTagRecord RUserChatTag err", zap.Error(err), zap.Int64("fromuid", fromUserId),
			zap.Int64("touid", toUserId))
		return nil, err
	}

	if len(resp.Info) <= 0 {
		return nil, nil
	}
	// todo
	return resp.Info[0], nil
}
