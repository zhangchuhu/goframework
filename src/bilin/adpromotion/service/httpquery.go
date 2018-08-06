package service

import (
	"bilin/adpromotion/entity"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

const (
	retSuccess = 0
)

var (
	httpClient *http.Client
)

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

type Qihu360Result struct {
	errno int    `json:"errno"`
	error string `json:"error"`
}

func ReportQihu360(clickInfo *entity.ClickInfo) (resp string, err error) {
	var (
		body     []byte
		request  *http.Request
		response *http.Response
		result   Qihu360Result
	)
	if request, err = http.NewRequest("GET", clickInfo.Callback_url, nil); err != nil {
		return
	}
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
	if result.errno != retSuccess {
		err = fmt.Errorf("query failed, ret code = %d, error :%s", result.errno, result.error)
	}

	return string(body), err
}
