package handler

import (
	"bilin/bcserver/domain/service"
	"bilin/operationManagement/entity"
	myservice "bilin/operationManagement/service"
	"bilin/operationManagement/token_sdk/go/src/vod_token"
	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Headgear_REQUESTURI string

const (
	Headgear_ADDITEM    Headgear_REQUESTURI = "add"
	Headgear_DELITEM    Headgear_REQUESTURI = "del"
	Headgear_UPDATEITEM Headgear_REQUESTURI = "update"
	Headgear_SEARCHITEM Headgear_REQUESTURI = "search"

	ResultHttpSuccess      int32 = 0
	ResultHttpFailed       int32 = 1
	ResultUserAlreadyExist int32 = 2

	//bs2 固定配置，不可变
	AppId           = 1375992228
	AppSecret       = "04d8a918_66f1_"
	UploadBaseUrl   = "https://bilinoperationmanagement.bs2ul.yy.com"
	DownloadBaseUrl = "https://bilinoperationmanagement.bs2dl.yy.com"
)

// retWrite marshal the result and write to client(get).
func retWrite(w http.ResponseWriter, r *http.Request, res map[string]interface{}, start time.Time) {
	data, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		log.Error(r.URL.String(), zap.Any("err", err), zap.Any("res", res))
		return
	}
	dataStr := string(data) + "\n"
	if _, err := w.Write([]byte(dataStr)); err != nil {
		log.Error(r.URL.String(), zap.Any("err", err), zap.Any("data", dataStr))
		return
	}
	log.Debug(r.URL.String(), zap.Any("ip", r.RemoteAddr), zap.Any("time", time.Now().Sub(start).Seconds()))
}

type HeadgearRequest struct {
	Uri  Headgear_REQUESTURI  `json:"uri"`
	Data *entity.HeadgearInfo `json:"data"`
}

type HeadgearResp struct {
	Result    int32                `json:"result"`
	ErrorDesc string               `json:"errordesc"`
	Uri       Headgear_REQUESTURI  `json:"uri"`
	Data      *entity.HeadgearInfo `json:"data"`
}

type AllHeadgearsResp struct {
	Result    int32                  `json:"result"`
	ErrorDesc string                 `json:"errordesc"`
	Data      []*entity.HeadgearInfo `json:"items"`
}

type Bs2Context struct {
	Bucket   string `json:"bucket"`
	FileName string `json:"filename"`
}

type Bs2TokenResp struct {
	Result      int32  `json:"result"`
	ErrorDesc   string `json:"errordesc"`
	Token       string `json:"token"`
	UploadUrl   string `json:"uploadurl"`
	DownloadUrl string `json:"downloadurl"`
}

type OptManagementHttpObj struct {
}

func stringTimeToUnixTimestamp(tobeChange string) int64 {
	timeLayout := "2006-01-02 15:04:05"                             //转化所需模板
	loc, _ := time.LoadLocation("Local")                            //重要：获取时区
	theTime, _ := time.ParseInLocation(timeLayout, tobeChange, loc) //使用模板在对应时区转化为time.time类型
	return theTime.Unix()                                           //转化为时间戳 类型是int64
}

func TimerReloadHeadgears(obj *OptManagementHttpObj, interval time.Duration) {
	const prefix = "TimerHandlerManager "

	for {
		log.Debug(prefix + "begin")

		infos, _ := myservice.MysqlGetAllVipUsers()

		now := time.Now().Unix()
		for _, item := range infos {
			effectTime := stringTimeToUnixTimestamp(item.EffectTime)
			expireTime := stringTimeToUnixTimestamp(item.ExpireTime)

			if expireTime <= effectTime || now < effectTime || now > expireTime {
				service.RedisDelUserHeadgear(uint64(item.Uid))
				continue
			}

			service.RedisSetUserHeadgear(uint64(item.Uid), item.Headgear)
		}

		time.Sleep(interval)
	}
}

func NewOptManagementHttpObj() (o *OptManagementHttpObj) {
	service.RedisInit()
	myservice.MysqlInit()

	//load data to redis
	_, err := myservice.MysqlGetAllVipUsers()
	if err != nil {
		panic("MysqlGetAllVipUsers error!")
	}

	o = &OptManagementHttpObj{}

	go TimerReloadHeadgears(o, 10*time.Second)
	return
}

func (o *OptManagementHttpObj) Hello(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var (
		res = make(map[string]interface{})
	)
	defer retWrite(w, r, res, time.Now())
	res["Hello"] = "success"
}

func (o *OptManagementHttpObj) HeadgearOperation(resp http.ResponseWriter, req *http.Request) {
	//便于测试，允许所有跨域访问
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Set("Access-Control-Allow-Credentials", "true")
	resp.Header().Set("Access-Control-Allow-Methods", "*")
	resp.Header().Set("Access-Control-Allow-Headers", "Content-Type,Access-Token")
	resp.Header().Set("Access-Control-Expose-Headers", "*")

	const prefix = "HeadgearOperation "
	if req.Method == "OPTIONS" {
		fmt.Fprintf(resp, "%d", http.StatusOK)
		return
	}

	if req.Method != "POST" {
		http.Error(resp, "Method Not Allowed", http.StatusMethodNotAllowed)
		log.Error(prefix+"Method Not Allowed", zap.Any("req.Method", req.Method))
		return
	}

	headgearResp := &HeadgearResp{Result: ResultHttpSuccess}
	headRequest := &HeadgearRequest{}
	var err error

	body, _ := ioutil.ReadAll(req.Body)
	if err = json.Unmarshal(body, headRequest); err != nil {
		log.Error(prefix+"json.Unmarshal failed", zap.Error(err), zap.Any("body", string(body)))
		headgearResp.Result = ResultHttpFailed
		goto RETURN
	}

	log.Debug("[+]"+prefix+"begin", zap.Any("headRequest", headRequest))

	if headRequest.Data.Uid == 0 {
		log.Error(prefix + "headRequest.Data.Uid == 0")
		err = fmt.Errorf("Data.Uid == 0")
		goto RETURN
	}

	switch headRequest.Uri {
	case Headgear_ADDITEM:
		err = myservice.MysqlAddVipUser(headRequest.Data)
		if err != nil {
			headgearResp.Result = ResultUserAlreadyExist
		}
	case Headgear_UPDATEITEM:
		err = myservice.MysqlUpdateVipUser(headRequest.Data)
		if err != nil {
			headgearResp.Result = ResultHttpFailed
		}
	case Headgear_DELITEM:
		err = myservice.MysqlDelVipUser(headRequest.Data.Uid)
		if err == nil {
			service.RedisDelUserHeadgear(uint64(headRequest.Data.Uid))
		} else {
			headgearResp.Result = ResultHttpFailed
		}
	case Headgear_SEARCHITEM:
		headgearResp.Data, _ = myservice.MysqlGetVipUser(headRequest.Data.Uid)
	default:
		err = fmt.Errorf("unknown uri %s", headRequest.Uri)
	}

RETURN:
	headgearResp.Uri = headRequest.Uri
	if err != nil {
		headgearResp.ErrorDesc = err.Error()
	}

	ret, _ := json.Marshal(headgearResp)
	fmt.Fprintf(resp, string(ret))

	log.Debug("[+]"+prefix+"end", zap.Any("headRequest", headRequest), zap.Any("headgearResp", headgearResp))

}

func (o *OptManagementHttpObj) GetAllHeadgears(resp http.ResponseWriter, req *http.Request) {

	resp.Header().Set("Access-Control-Allow-Origin", "*")

	const prefix = "GetAllHeadgears "
	if req.Method != "GET" {
		http.Error(resp, "Method Not Allowed", http.StatusMethodNotAllowed)
		log.Error(prefix+"Method Not Allowed", zap.Any("req.Method", req.Method))
		return
	}

	log.Debug("[+]" + prefix + "begin")

	headgearResp := &AllHeadgearsResp{Result: ResultHttpSuccess}
	headgearResp.Data, _ = myservice.MysqlGetAllVipUsers()

	ret, _ := json.Marshal(headgearResp)
	fmt.Fprintf(resp, string(ret))

	log.Debug("[+]"+prefix+"end", zap.Any("headgearResp", headgearResp))

}

func (o *OptManagementHttpObj) GenBS2Token(resp http.ResponseWriter, req *http.Request) {

	resp.Header().Set("Access-Control-Allow-Origin", "*")

	const prefix = "GenBS2Token "

	req.ParseForm()
	log.Debug("[+]"+prefix+"begin", zap.Any("filename", req.Form["filename"][0]))

	//生成一个filename，用于上传bs2
	strNow := strconv.Itoa(int(time.Now().Unix()) + rand.Int())
	bs2FileName := fmt.Sprintf("%x", md5.Sum([]byte(strNow))) + "." + strings.Split(req.Form["filename"][0], ".")[1]

	ctx := &Bs2Context{Bucket: "bilinoperationmanagement"}
	ctx.FileName = bs2FileName
	bs2ctx, _ := json.Marshal(ctx)
	info := &vod_token.TokenInfo{
		Appid:   AppId,
		Uid:     0,
		Ttl:     uint64(time.Now().Unix()) + 300, // token有效期5分钟
		Auth:    0,
		Context: string(bs2ctx),
	}
	token, ret := vod_token.GenToken(1, AppSecret, info, vod_token.TYPE_VOD)
	//str := base64.StdEncoding.EncodeToString([]byte(token))
	log.Debug("[+]"+prefix+"GenToken", zap.Any("token", token), zap.Any("ret", ret))

	ret = vod_token.ValidateToken(token, AppSecret)
	log.Debug("[+]"+prefix+"validate token:", zap.Any("token", token), zap.Any("ret", ret))

	tokenResp := &Bs2TokenResp{Result: ResultHttpSuccess, Token: token}

	tokenResp.UploadUrl = UploadBaseUrl + "/" + ctx.FileName
	tokenResp.DownloadUrl = DownloadBaseUrl + "/" + ctx.FileName
	strRsp, _ := json.Marshal(tokenResp)
	fmt.Fprintf(resp, string(strRsp))

	ret, ninfo := vod_token.GetProperty(token)
	log.Debug("[+]"+prefix+"end", zap.Any("ret", ret), zap.Any("ninfo", ninfo), zap.Any("strRsp", strRsp))
}

//主播个人流水
func (o *OptManagementHttpObj) GetAllLivingRecord(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Access-Control-Allow-Origin", "*")

	const prefix = "GetAllLivingRecord "
	if req.Method != "GET" {
		http.Error(resp, "Method Not Allowed", http.StatusMethodNotAllowed)
		log.Error(prefix+"Method Not Allowed", zap.Any("req.Method", req.Method))
		return
	}

	log.Debug("[+]" + prefix + "begin")

	req.ParseForm()
	log.Debug(req.URL.String(), zap.Any("Query", req.Form))

	var queryDate string
	if req.Form["date"] == nil || len(req.Form["date"][0]) == 0 {
		queryDate = strings.Replace(time.Now().String()[0:10], "-", "", -1)
	} else {
		queryDate = req.Form["date"][0]
	}

	result := &entity.AllLivingRecordInfoResp{Result: ResultHttpSuccess}
	//默认取当天的开播流水信息
	result.Data, _ = myservice.MysqlGetAllLivingRecord(queryDate)

	ret, _ := json.Marshal(result)
	fmt.Fprintf(resp, string(ret))

	log.Debug("[+]"+prefix+"end", zap.Any("result", result))
}
