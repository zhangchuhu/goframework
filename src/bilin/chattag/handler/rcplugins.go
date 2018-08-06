package handler

import (
	"bilin/clientcenter"
	"bilin/protocol/userinfocenter"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"context"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

const (
	puaname     = "聊妹套话"
	truthname   = "真心话"
	maleTitle   = "想要Get撩妹技能，快戳下方~~~"
	femaleTitle = "不知道怎么聊天？不要慌，快戳下方~~~"
)

type RCPlugin struct {
	PluginID   int64  `json:"plugin_id"`
	PluginName string `json:"plugin_name"`
}

type RCPluginList struct {
	Title     string     `json:"title"` // 文案标题
	RCPlugins []RCPlugin `json:"rc_plugins"`
}

func HandleRCPlugin(c *gin.Context) *HttpError {
	appzaplog.Debug("[+]HandleRCPlugin")
	var (
		body *RCPluginList
		ret  = successHttp
	)

	for {
		uid := c.GetInt64(uidContextKey)
		resp, err := clientcenter.UserInfoClient().GetSingleUserInfo(context.TODO(), &userinfocenter.GetSingleUserInfoReq{
			Uid: uint64(uid),
		})
		if err != nil || resp.Uinfo == nil {
			appzaplog.Error("HandleRCPlugin GetSingleUserInfo err", zap.Error(err), zap.Int64("uid", uid))
			ret = getUserInfoHttpErr
			break
		}
		switch resp.Uinfo.Showsex {
		case 1:
			body = &RCPluginList{
				Title: maleTitle,
				RCPlugins: []RCPlugin{
					{PluginID: 1, PluginName: puaname},
					{PluginID: 2, PluginName: truthname},
				},
			}
		case 0:
			body = &RCPluginList{
				Title: femaleTitle,
				RCPlugins: []RCPlugin{
					{PluginID: 2, PluginName: truthname},
				},
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
	appzaplog.Debug("[-]HandleRCPlugin", zap.Any("resp", httpret))
	return ret
}

type UserDetailInCall struct {
	RCPluginList RCPluginList  `json:"rc_plugin_list"`
	Dynamic      Dynamic       `json:"dynamic"`
	Like         string        `json:"like"`
	Musics       []LikeElement `json:"musics"`
	Movies       []LikeElement `json:"movies"`
	Books        []LikeElement `json:"books"`
}

type Dynamic struct {
	CreatedTime int64        `json:"created_time"`
	Content     string       `json:"content"`
	ImgList     []DynamicImg `json:"img_list"`
}
type DynamicImg struct {
	BigUrl   string `json:"big_url"`
	SmallUrl string `json:"small_url"`
	Size     string `json:"size"`
}

type LikeElement struct {
	Title   string `json:"title"`
	BilinID int64  `json:"bilin_id"`
	OtherID int64  `json:"other_id"`
	Image   string `json:"image"`
}

func HandleUserDetailInCall(c *gin.Context) *HttpError {
	appzaplog.Debug("[+]HandleUserDetailInCall")
	var (
		body *UserDetailInCall = &UserDetailInCall{}
		ret                    = successHttp
	)

	for {
		uid := c.GetInt64(uidContextKey)
		resp, err := clientcenter.UserInfoClient().GetSingleUserInfo(context.TODO(), &userinfocenter.GetSingleUserInfoReq{
			Uid: uint64(uid),
		})
		if err != nil || resp.Uinfo == nil {
			appzaplog.Error("HandleUserDetailInCall GetSingleUserInfo err", zap.Error(err), zap.Int64("uid", uid))
			ret = getUserInfoHttpErr
			break
		}
		switch resp.Uinfo.Showsex {
		case 1:
			body.RCPluginList = RCPluginList{
				Title: maleTitle,
				RCPlugins: []RCPlugin{
					{PluginID: 1, PluginName: puaname},
					{PluginID: 2, PluginName: truthname},
				},
			}
		case 0:
			body.RCPluginList = RCPluginList{
				Title: femaleTitle,
				RCPlugins: []RCPlugin{
					{PluginID: 2, PluginName: truthname},
				},
			}
		}
		toUerIdStr := c.PostForm("to_userid")
		if toUerIdStr == "" {
			appzaplog.Warn("HandleUserDetailInCall to_userid not exist")
			ret = parseFormHttpErr
			break
		}
		toUserId, err := strconv.ParseInt(toUerIdStr, 10, 64)
		if err != nil {
			appzaplog.Warn("HandleUserDetailInCall ParseInt to_userid err", zap.Error(err), zap.Int64("uid", uid), zap.String("touid", toUerIdStr))
			ret = parseFormHttpErr
			break
		}
		dinfo, err := clientcenter.QueryDynamicAndDynamicAttrListByUserByPage(uid, toUserId, time.Now().Unix()*1000, 1)
		if err != nil {
			appzaplog.Warn("HandleUserDetailInCall QueryDynamicAndDynamicAttrListByUserByPage err", zap.Error(err), zap.Int64("uid", uid), zap.String("touid", toUerIdStr))
			ret = daoGetHttpErr
			break
		}
		for _, v := range dinfo {
			if v.Dynamic.IsDelete == 0 {
				body.Dynamic = Dynamic{
					CreatedTime: v.Dynamic.CreateOn,
					Content:     "分享图片",
					ImgList: []DynamicImg{
						{
							BigUrl:   v.Dynamic.FirstImgBigUrl,
							SmallUrl: v.Dynamic.FirstImgSmallUrl,
							Size:     v.Dynamic.Size1,
						},
					},
				}
				if v.Dynamic.Content != "" {
					body.Dynamic.Content = v.Dynamic.Content
				}
				break
			}
		}

		minfo, err := clientcenter.GetUserHobbies(toUserId)
		if err != nil {
			appzaplog.Warn("HandleUserDetailInCall QueryDynamicAndDynamicAttrListByUserByPage err", zap.Error(err), zap.Int64("uid", uid), zap.String("touid", toUerIdStr))
			ret = daoGetHttpErr
			break
		}
		if minfo != nil {
			for _, v := range minfo.Movie {
				body.Movies = append(body.Movies, LikeElement{
					Title:   v.Title,
					BilinID: v.BilinID,
					OtherID: v.OtherID,
					Image:   v.Image,
				})
			}
			for _, v := range minfo.Book {
				body.Books = append(body.Books, LikeElement{
					Title:   v.Title,
					BilinID: v.BilinID,
					OtherID: v.OtherID,
					Image:   v.Image,
				})
			}
			for _, v := range minfo.Music {
				body.Musics = append(body.Musics, LikeElement{
					Title:   v.Title,
					BilinID: v.BilinID,
					OtherID: v.OtherID,
					Image:   v.Image,
				})
			}
		}
		if uinfo, err := clientcenter.GetUserInfoByUserId(toUserId); err == nil {
			body.Like = uinfo.Like
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
	appzaplog.Debug("[-]HandleUserDetailInCall", zap.Any("resp", httpret))
	return ret
}
