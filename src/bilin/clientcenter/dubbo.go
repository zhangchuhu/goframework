package clientcenter

import (
	"bilin/protocol"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/tars"
)

var (
	dubboCommonClient bilin.DubboProxyClient
)

func init() {
	comm := tars.NewCommunicator()
	dubboCommonClient = bilin.NewDubboProxyClient("bilin.dubboproxy.CommonProxy", comm)
}

func DubboCommonClient() bilin.DubboProxyClient {
	return dubboCommonClient
}

// VerifyAccessToken with userid, returns err!=nil if and only if dubbo service is unavailable.
func VerifyAccessToken(accesstoken string, userid string) (ok bool, err error) {
	rsp, err := dubboCommonClient.Invoke(context.Background(), &bilin.DPInvokeReq{
		Service: "com.bilin.user.account.service.IUserLoginService",
		Method:  "verifyUserAccessToken",
		Args: []*bilin.DPInvokeArg{
			{
				Type:  "long",
				Value: userid,
			},
			{
				Type:  "java.lang.String",
				Value: accesstoken,
			},
		},
	})
	const errmsg = "fail to verify access token"
	if err != nil {
		err = fmt.Errorf("error calling dubbo proxy: %v", err)
		log.Error(errmsg, zap.Error(err))
		return
	}
	if rsp.ThrewException {
		err = fmt.Errorf("error invoking dubbo service: %v", rsp.Type)
		log.Error(errmsg, zap.Error(err), zap.String("detail", rsp.Value))
		return
	}
	if rsp.Type == "java.lang.Boolean" && rsp.Value == "true" {
		// success
		ok = true
	} else {
		log.Warn(errmsg, zap.String("token", accesstoken), zap.String("uid", userid),
			zap.String("rsp.Type", rsp.Type), zap.String("rsp.Value", rsp.Value))
		return
	}
	return
}

func AddFlower(flowerNum, userID int64) {
	rsp, err := dubboCommonClient.Invoke(context.Background(), &bilin.DPInvokeReq{
		Service: "com.bilin.user.center.service.IUserCenterService",
		Method:  "increaseGlamourValueByFlower",
		Args: []*bilin.DPInvokeArg{
			{
				Type:  "long",
				Value: strconv.FormatInt(userID, 10),
			},
			{
				Type:  "long",
				Value: strconv.FormatInt(flowerNum, 10),
			},
		},
	})
	const errmsg = "fail to add flower"
	if err != nil {
		err = fmt.Errorf("error calling dubbo proxy: %v", err)
		log.Error(errmsg, zap.Error(err))
		return
	}
	if rsp.ThrewException {
		err = fmt.Errorf("error invoking dubbo service: %v", rsp.Type)
		log.Error(errmsg, zap.Error(err), zap.String("detail", rsp.Value))
		return
	}
}

type BilinKtv struct {
	Name     string `json:"name"`
	Artist   string `json:"artist"`
	Duration int64  `json:"duration"`
	UploadBy string `json:"uploadBy"`
	Pkg      string `json:"pkg"`
	PkgLen   int64  `json:"pkgLen"`
	PkgMd5   string `json:"pkgMd5"`
}

func GetBilinKtvById(id int64) (result BilinKtv, found bool) {
	rsp, err := dubboCommonClient.Invoke(context.Background(), &bilin.DPInvokeReq{
		Service: "com.bilin.call.service.IBilinKtvService",
		Method:  "getBilinKtvById",
		Args: []*bilin.DPInvokeArg{
			{
				Type:  "java.lang.Long",
				Value: strconv.FormatInt(id, 10),
			},
		},
	})
	const errmsg = "fail to getBilinKtvById"
	if err != nil {
		err = fmt.Errorf("error calling dubbo proxy: %v", err)
		log.Error(errmsg, zap.Error(err))
		return
	}
	if rsp.ThrewException {
		err = fmt.Errorf("error invoking dubbo service: %v", rsp.Type)
		log.Error(errmsg, zap.Error(err), zap.String("detail", rsp.Value))
		return
	}
	if rsp.Type == "" && rsp.Value == "null" {
		found = false
	} else {
		found = true
	}
	if err = json.Unmarshal([]byte(rsp.Value), &result); err != nil {
		err = fmt.Errorf("error deserializing java bean: %v", rsp.Type)
		log.Error(errmsg, zap.Error(err), zap.String("detail", rsp.Value))
		return
	}
	return
}

type PersonalUserInfo struct {
	Like             string `json:"like"`
	NotLike          string `json:"notLike"`
	IntroMe          string `json:"introMe"`
	City             int64  `json:"city"`
	CityName         string `json:"cityName"`
	EvalutionLike    int64  `json:"evalutionLike"`    //个人评价-喜欢
	EvalutionNotLike int64  `json:"evalutionNotLike"` //个人评价-不喜欢
}

func GetUserInfoByUserId(userID int64) (PersonalUserInfo, error) {
	var ret PersonalUserInfo
	rsp, err := dubboCommonClient.Invoke(context.Background(), &bilin.DPInvokeReq{
		Service: "com.bilin.user.center.service.IUserCenterService",
		Method:  "getUserInfoByUserId",
		Args: []*bilin.DPInvokeArg{
			{
				Type:  "long",
				Value: strconv.FormatInt(userID, 10),
			},
		},
	})
	const errmsg = "fail to GetUserInfoByUserId"
	if err != nil {
		err = fmt.Errorf("error calling dubbo proxy: %v", err)
		log.Error(errmsg, zap.Error(err))
		return ret, err
	}
	if rsp == nil || rsp.ThrewException || rsp.Type != "com.bilin.user.center.bean.UserInfo" {
		err = fmt.Errorf("error invoking dubbo service: %v", rsp.Type)
		log.Error(errmsg, zap.Error(err), zap.String("detail", rsp.Value))
		return ret, err
	}
	err = json.Unmarshal([]byte(rsp.Value), &ret)
	if err != nil {
		log.Error(errmsg, zap.Error(err), zap.String("detail", rsp.Value))
		return ret, err
	}
	log.Debug("[-]GetUserInfoByUserId", zap.Any("resp", rsp))
	return ret, nil
}

type LikeElement struct {
	Title     string `json:"title"`
	BilinID   int64  `json:"bilinId"`
	OtherID   int64  `json:"otherId"`
	Image     string `json:"image"`
	CreatedOn int64  `json:"createOn"` // ms
}

type Hobbies struct {
	Movie []LikeElement `json:"MOVIE"`
	Music []LikeElement `json:"MUSIC"`
	Book  []LikeElement `json:"BOOK"`
}

func GetUserHobbies(userID int64) (*Hobbies, error) {
	rsp, err := dubboCommonClient.Invoke(context.Background(), &bilin.DPInvokeReq{
		Service: "com.bilin.user.center.service.IUserCenterService",
		Method:  "getUserHobbies",
		Args: []*bilin.DPInvokeArg{
			{
				Type:  "long",
				Value: strconv.FormatInt(userID, 10),
			},
		},
	})

	const errmsg = "fail to GetUserHobbies"
	if err != nil {
		err = fmt.Errorf("error calling dubbo proxy: %v", err)
		log.Error(errmsg, zap.Error(err))
		return nil, err
	}
	if rsp == nil || rsp.ThrewException || rsp.Type != "java.util.HashMap" {
		err = fmt.Errorf("error invoking dubbo service: %v", rsp.Type)
		log.Error(errmsg, zap.Error(err), zap.String("detail", rsp.Value))
		return nil, err
	}
	log.Debug("[-]GetUserInfoByUserId", zap.Any("resp", rsp))
	var hobbie Hobbies
	if err := json.Unmarshal([]byte(rsp.Value), &hobbie); err != nil {
		log.Error(errmsg, zap.Error(err))
		return nil, err
	}
	return &hobbie, nil
}

type Dynamic struct {
	FirstImgSmallUrl string `json:"firstImgSmallUrl"`
	FirstImgBigUrl   string `json:"firstImgBigUrl"`
	SmallUrl1        string `json:"smallUrl1"`
	SmallUrl2        string `json:"smallUrl2"`
	Content          string `json:"content"`
	FirstImgUrl      string `json:"firstImgUrl"`
	Url1             string `json:"url1"`
	Size1            string `json:"size1"`
	Url2             string `json:"url1"`
	Size2            string `json:"size2"`
	Url3             string `json:"url5"`
	Size3            string `json:"size3"`
	Url4             string `json:"url5"`
	Size4            string `json:"size4"`
	Url5             string `json:"url5"`
	Size5            string `json:"size5"`
	Url6             string `json:"url6"`
	Size6            string `json:"size6"`
	Url7             string `json:"url7"`
	Size7            string `json:"size7"`
	Url8             string `json:"url8"`
	Size8            string `json:"size8"`
	Url9             string `json:"url9"`
	Size9            string `json:"size9"`
	AllPhotoForbid   bool   `json:"allPhotoForbid"`
	ID               int64  `json:"id"`

	BigUrl1  string `json:"bigUrl1"`
	BigUrl2  string `json:"bigUrl2"`
	BigUrl3  string `json:"bigUrl3"`
	BigUrl4  string `json:"bigUrl4"`
	BigUrl5  string `json:"bigUrl5"`
	BigUrl6  string `json:"bigUrl6"`
	BigUrl7  string `json:"bigUrl7"`
	BigUrl8  string `json:"bigUrl8"`
	BigUrl9  string `json:"bigUrl9"`
	UserId   int64  `json:"userId"`
	IsHidden int64  `json:"isHidden"`
	IsDelete int64  `json:"isDelete"`
	CreateOn int64  `json:"createOn"`
}

type DynamicAttr struct {
	TotalPraiseNum  int64  `json:"totalPraiseNum"`
	TotalCommentNum int64  `json:"totalCommentNum"`
	IsComment       int64  `json:"isComment"`
	IsPraise        int64  `json:"isPraise"`
	JumpStr         string `json:"jumpStr"`
	TypeId          int64  `json:"type"`
}

type DynamicValue struct {
	DynamicAttr DynamicAttr `json:"dynamicAttr"`
	Dynamic     Dynamic     `json:"dynamic"`
}

func QueryDynamicAndDynamicAttrListByUserByPage(userID, dynamicUserId, timestamp, pageSize int64) ([]DynamicValue, error) {
	rsp, err := dubboCommonClient.Invoke(context.Background(), &bilin.DPInvokeReq{
		Service: "com.bilin.user.dynamic.service.IUserDynamicService",
		Method:  "queryDynamicAndDynamicAttrListByUserByPage",
		Args: []*bilin.DPInvokeArg{
			{
				Type:  "long",
				Value: strconv.FormatInt(userID, 10),
			},
			{
				Type:  "long",
				Value: strconv.FormatInt(dynamicUserId, 10),
			},
			{
				Type:  "long",
				Value: strconv.FormatInt(timestamp, 10),
			},
			{
				Type:  "int",
				Value: strconv.FormatInt(pageSize, 10),
			},
		},
	})

	const errmsg = "fail to QueryDynamicAndDynamicAttrListByUserByPage"
	if err != nil {
		err = fmt.Errorf("error calling dubbo proxy: %v", err)
		log.Error(errmsg, zap.Error(err))
		return nil, err
	}
	if rsp == nil || rsp.ThrewException {
		err = fmt.Errorf("error invoking dubbo service: %v", rsp.Type)
		log.Error(errmsg, zap.Error(err), zap.String("detail", rsp.Value))
		return nil, err
	}
	log.Debug("[-]QueryDynamicAndDynamicAttrListByUserByPage", zap.Any("resp", rsp))
	var dynamic []DynamicValue
	if err := json.Unmarshal([]byte(rsp.Value), &dynamic); err != nil {
		log.Error(errmsg, zap.Error(err), zap.Error(err), zap.String("detail", rsp.Value))
		return nil, err
	}
	return dynamic, nil
}

type BilinFriend struct {
	UserID int64 `json:"userId"`
}

func QueryFriendList(id int64) (result []BilinFriend) {
	rsp, err := dubboCommonClient.Invoke(context.Background(), &bilin.DPInvokeReq{
		Service: "com.bilin.user.relation.service.IUserRelationService",
		Method:  "queryFriendList",
		Args: []*bilin.DPInvokeArg{
			{
				Type:  "long",
				Value: strconv.FormatInt(id, 10),
			},
		},
	})
	const errmsg = "fail to queryFriendList"
	if err != nil {
		err = fmt.Errorf("error calling dubbo proxy: %v", err)
		log.Error(errmsg, zap.Error(err))
		return
	}
	if rsp.ThrewException {
		err = fmt.Errorf("error invoking dubbo service: %v", rsp.Type)
		log.Error(errmsg, zap.Error(err), zap.String("detail", rsp.Value))
		return
	}
	if err = json.Unmarshal([]byte(rsp.Value), &result); err != nil {
		err = fmt.Errorf("error deserializing java bean: %v", rsp.Type)
		log.Error(errmsg, zap.Error(err), zap.String("detail", rsp.Value))
		return
	}
	return
}

type BilinUser struct {
	UserID int64 `json:"id"`
}

// 通过比邻号搜索用户
func GetUserByBLId(id int64) (result BilinUser, found bool) {
	rsp, err := dubboCommonClient.Invoke(context.Background(), &bilin.DPInvokeReq{
		Service: "com.bilin.user.center.service.IUserCenterService",
		Method:  "getUserByBLId",
		Args: []*bilin.DPInvokeArg{
			{
				Type:  "long",
				Value: strconv.FormatInt(id, 10),
			},
		},
	})
	const errmsg = "fail to getUserByBLId"
	if err != nil {
		err = fmt.Errorf("error calling dubbo proxy: %v", err)
		log.Error(errmsg, zap.Error(err))
		return
	}
	if rsp.ThrewException {
		err = fmt.Errorf("error invoking dubbo service: %v", rsp.Type)
		log.Error(errmsg, zap.Error(err), zap.String("detail", rsp.Value))
		return
	}
	if rsp.Type == "" && rsp.Value == "null" {
		found = false
	} else {
		found = true
	}
	if err = json.Unmarshal([]byte(rsp.Value), &result); err != nil {
		err = fmt.Errorf("error deserializing java bean: %v", rsp.Type)
		log.Error(errmsg, zap.Error(err), zap.String("detail", rsp.Value))
		return
	}
	return
}

type BilinAttention struct {
	UserID int64 `json:"userId"`
}
/**
   add result limit
 */
func QueryAttentionList(id int64) (result []BilinAttention) {
	rsp, err := dubboCommonClient.Invoke(context.Background(), &bilin.DPInvokeReq{
		Service: "com.bilin.attention.center.service.IAttentionCenterService",
		Method:  "queryMyAttentionUserList",
		Args: []*bilin.DPInvokeArg{
			{
				Type:  "long",
				Value: strconv.FormatInt(id, 10),
			},
			{
				Type:  "long",
				Value: "0",
			},
			{
				Type:  "int",
				Value: "1000",
			},
		},
	})
	const errmsg = "fail to queryMyAttentionUserList"
	if err != nil {
		err = fmt.Errorf("error calling dubbo proxy: %v", err)
		log.Error(errmsg, zap.Error(err))
		return
	}
	if rsp.ThrewException {
		err = fmt.Errorf("error invoking dubbo service: %v", rsp.Type)
		log.Error(errmsg, zap.Error(err), zap.String("detail", rsp.Value))
		return
	}
	if err = json.Unmarshal([]byte(rsp.Value), &result); err != nil {
		err = fmt.Errorf("error deserializing java bean: %v", rsp.Type)
		log.Error(errmsg, zap.Error(err), zap.String("detail", rsp.Value))
		return
	}
	return
}
