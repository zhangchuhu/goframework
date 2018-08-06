package onlinepush

import (
	"bilin/protocol"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
	pb "github.com/golang/protobuf/proto"
	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
)

const (
	retSuccess = 1
)

var (
	// URL should be set before usage.
	URL        = "http://test-goim.yy.com:7172/1"
	httpClient *http.Client
)

type pushResult struct {
	RetCode int     `json:"ret"`
	Offline []int64 `json:"offline,omitempty"`
}

func init() {
	httpDialer := &net.Dialer{
		Timeout: 5 * time.Second,
	}
	httpTransport := &http.Transport{
		DialContext:       httpDialer.DialContext,
		DisableKeepAlives: false,
	}
	httpClient = &http.Client{
		Transport: httpTransport,
		Timeout:   5 * time.Second,
	}
}
func SetUrl(pushUrl string)(){
	URL = pushUrl
}
// PushToUser tells which users are offline.
// Offline users can not receive push messages. Server should retry push later.
func PushToUser(mpush bilin.MultiPush) (offline []int64, err error) {
	var (
		body       []byte
		bodyReader io.Reader
		request    *http.Request
		response   *http.Response
		result     pushResult
		url        = URL + "/pushs"
	)
	if body, err = pb.Marshal(&mpush); err != nil {
		return
	}
	bodyReader = bytes.NewBuffer(body)
	if request, err = http.NewRequest("POST", url, bodyReader); err != nil {
		return
	}
	request.Header.Set("Content-Type", "application/x-protobuf")
	if response, err = httpClient.Do(request); err != nil {
		return
	}
	defer response.Body.Close()
	if body, err = ioutil.ReadAll(response.Body); err != nil {
		return
	}
	if err = json.Unmarshal(body, &result); err != nil {
		return
	}
	if result.RetCode != retSuccess {
		err = fmt.Errorf("push failed, ret code = %d", result.RetCode)
		return
	}
	offline = result.Offline
	return
}

func PushToRoom(roomid int64, push bilin.ServerPush) (err error) {
	var (
		body       []byte
		bodyReader io.Reader
		request    *http.Request
		response   *http.Response
		result     pushResult
		url        = fmt.Sprintf(URL+"/push/room?rid=%d", roomid)
	)
	if body, err = pb.Marshal(&push); err != nil {
		return
	}
	log.Error("PushToRoom:",zap.Any("url:",url))
	bodyReader = bytes.NewBuffer(body)
	if request, err = http.NewRequest("POST", url, bodyReader); err != nil {
		return
	}
	request.Header.Set("Content-Type", "application/x-protobuf")
	if response, err = httpClient.Do(request); err != nil {
		return
	}
	defer response.Body.Close()
	if body, err = ioutil.ReadAll(response.Body); err != nil {
		return
	}
	if err = json.Unmarshal(body, &result); err != nil {
		return
	}
	if result.RetCode != retSuccess {
		err = fmt.Errorf("push failed, ret code = %d", result.RetCode)
		return
	}
	return
}
