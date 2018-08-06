package handler

import (
	"bilin/chattag/cache"
	"bilin/clientcenter"
	"bilin/protocol"
	"bilin/protocol/userinfocenter"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"context"
	"github.com/gin-gonic/gin"
	"sort"
	"strconv"
	"strings"
)

type ChatTag struct {
	TagId       int64  `json:"tag_id"`
	TagName     string `json:"tag_name"`
	Color       string `json:"color,omitempty"`
	TotalTagNum int64  `json:"total_tag_num,omitempty"`
}

type ChatTagList struct {
	ChatTags []ChatTag `json:"chat_tags"`
}

type GetChatTagsParam struct {
	UserId int64
	Top    int64
}

func (p *GetChatTagsParam) Unmarshal(c *gin.Context) error {
	var (
		err error
	)
	if userIdStr, ok := c.GetPostForm("user_id"); ok {
		if p.UserId, err = strconv.ParseInt(userIdStr, 10, 64); err != nil {
			appzaplog.Error("GetChatTagsParam.Unmarshal ParseInt err", zap.Error(err))
			return err
		}
	}
	if topStr, ok := c.GetPostForm("top"); ok {
		if p.Top, err = strconv.ParseInt(topStr, 10, 64); err != nil {
			appzaplog.Error("GetChatTagsParam.Unmarshal ParseInt err", zap.Error(err))
			return err
		}
	}
	appzaplog.Debug("[-]GetChatTagsParam.Unmarshal", zap.Any("resp", p))
	return nil
}

func GetChatTagsList(c *gin.Context) *HttpError {
	appzaplog.Debug("[+]GetChatTagsList")
	var (
		params GetChatTagsParam
		body   *ChatTagList
		err    error
		ret    = successHttp
	)

	for {
		if err = params.Unmarshal(c); err != nil {
			appzaplog.Error("GetChatTagsList param.Unmarshal err", zap.Error(err))
			ret = parseFormHttpErr
			break
		}
		switch {
		case params.UserId > 0 && params.Top > 0:
			body = TagListByUser(params.UserId, params.Top)
		case params.UserId > 0:
			// todo, use 1000 to return all, since no more than 1000 tags now
			body = TagListByUser(params.UserId, 1000)
		default:
			body = TagList()
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
	appzaplog.Debug("[-]GetChatTagsList", zap.Any("resp", httpret))
	return ret
}

func TagList() *ChatTagList {
	tags := cache.TakeChatTagCache()
	if tags == nil {
		appzaplog.Warn("TagList cache.TakeChatTagCache empty")
		return nil
	}
	data := &ChatTagList{
		ChatTags: make([]ChatTag, 0),
	}
	for _, v := range tags {
		if v != nil {
			data.ChatTags = append(data.ChatTags, ChatTag{
				TagName: v.TagName,
				TagId:   int64(v.Id),
				Color:   v.TagColor,
			})
		} else {
			appzaplog.Warn("TagList nil chattag")
		}
	}
	return data
}

func TagListByUser(userId int64, topn int64) *ChatTagList {
	body := &ChatTagList{ChatTags: make([]ChatTag, 0)}
	if !isMale(userId) {
		return body
	}

	utags, err := clientcenter.ChatTagClient().RTopNUserChatTagSummary(context.TODO(), &bilin.RTopNUserChatTagSummaryReq{
		Topuser: &bilin.TopNUser{
			Touserid: userId,
			Topn:     topn,
		},
	})
	if err != nil {
		appzaplog.Error("TagListByUser RTopNUserChatTagSummary err", zap.Error(err), zap.Int64("uid", userId))
		return body
	}

	if utags.Summary == nil {
		return body
	}
	for _, v := range utags.Summary.Summary {
		body.ChatTags = append(body.ChatTags, ChatTag{
			TagId:       v.Tagid,
			TagName:     v.Tagname,
			Color:       v.Tagcolor,
			TotalTagNum: v.Totaltagnum,
		})
	}
	return body
}

func convertTagId(tagid string) (tagids []int64) {
	if tagid == "" {
		return []int64{}
	}
	tagids_ := strings.Split(tagid, ",")
	for _, v := range tagids_ {
		if tagIdInt, err := strconv.ParseInt(v, 10, 64); err == nil {
			tagids = append(tagids, tagIdInt)
		}
	}
	return
}

func topNTagList(in *ChatTagList, topn int64) *ChatTagList {
	ret := &ChatTagList{
		ChatTags: make([]ChatTag, 0),
	}
	if in == nil || len(in.ChatTags) == 0 {
		return ret
	}
	sort.SliceStable(in.ChatTags, func(i, j int) bool {
		return in.ChatTags[i].TotalTagNum > in.ChatTags[j].TotalTagNum
	})
	if topn > int64(len(in.ChatTags)) {
		topn = int64(len(in.ChatTags))
	}
	in.ChatTags = in.ChatTags[0:topn]
	return in
}

func isMale(userid int64) bool {
	uids := []uint64{uint64(userid)}
	uinfo, err := clientcenter.UserInfoClient().GetUserInfo(context.TODO(), &userinfocenter.GetUserInfoReq{uids})
	if err != nil {
		appzaplog.Error("isMale.GetUserInfo err", zap.Error(err), zap.Int64("uid", userid))
		//出错时返回，继续后续逻辑
		return true
	}
	if info, ok := uinfo.Users[uint64(userid)]; ok {
		return info.Showsex == 1
	}
	return true
}
