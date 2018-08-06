package onlinequery

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

const (
	retSuccess = 1
)

var (
	// URL should be set before usage.
	URL        = "http://test-goim.yy.com:7172/1"
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

type userSession struct {
	RetCode int `json:"ret"`
	Session struct {
		UserId  int64
		Count   int32
		Seq     int32
		Servers map[int32]struct {
			Comet     int32
			Birth     string
			Heartbeat string
		}
		Rooms map[int64]map[int32]int32 // roomid:seq:server
	} `json:"session"`
}

type roomUsers struct {
	RetCode int             `json:"ret"`
	Users   map[int64]int32 `json:"users"`
}

// GetUserRoom returns rid
//   -2  user offline
//   -1  user online, but not in any room
//   >=0 user online, and in room rid
func GetUserRoom(uid int64) (rid int64, err error) {
	var (
		body     []byte
		request  *http.Request
		response *http.Response
		result   userSession
		url      = fmt.Sprintf(URL+"/session?uid=%d", uid)
	)
	if request, err = http.NewRequest("GET", url, nil); err != nil {
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
	if result.RetCode != retSuccess {
		err = fmt.Errorf("query failed, ret code = %d", result.RetCode)
		return
	}
	if result.Session.Count <= 0 { // user offline
		rid = -2
		return
	}
	if len(result.Session.Rooms) == 0 { // not in any room
		rid = -1
		return
	}
	var seqmax int32
	for roomid, val := range result.Session.Rooms { // get room with maximum seq
		for seq := range val {
			if seq >= seqmax {
				rid = roomid
				seqmax = seq
			}
		}
	}
	return
}

func GetRoomUser(rid int64) (users map[int64]int32, err error) {
	var (
		body     []byte
		request  *http.Request
		response *http.Response
		result   roomUsers
		url      = fmt.Sprintf(URL+"/room?rid=%d", rid)
	)
	if request, err = http.NewRequest("GET", url, nil); err != nil {
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
	if result.RetCode != retSuccess {
		err = fmt.Errorf("query failed, ret code = %d", result.RetCode)
		return
	}
	users = result.Users
	return
}

type userCount struct {
	RetCode int `json:"ret"`
	Data    []struct {
		Server int
		Count  int
	} `json:"data"`
}

func UserCount() (count int, err error) {
	var (
		body     []byte
		request  *http.Request
		response *http.Response
		result   userCount
		url      = URL + "/count?type=server"
	)
	if request, err = http.NewRequest("GET", url, nil); err != nil {
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
	if result.RetCode != retSuccess {
		err = fmt.Errorf("query failed, ret code = %d", result.RetCode)
		return
	}
	for i := range result.Data {
		count += result.Data[i].Count
	}
	return
}
