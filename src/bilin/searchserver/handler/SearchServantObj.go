package handler

import (
	bc "bilin/bcserver/domain/service"
	"bilin/clientcenter"
	"bilin/protocol"
	"bilin/protocol/userinfocenter"
	"bilin/searchserver/config"
	e "bilin/searchserver/entity"
	u "bilin/searchserver/updater"
	"context"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
)

const (
	officailRoom = bc.OFFICAIL_ROOM
)

var (
	emptyUserInfo = &userinfocenter.UserInfo{
		Birthday: time.Now().Unix() - 18*365*24*60*60,
	}
	emptyRoomInfo = &bilin.RoomInfo{
		Starttime: uint64(time.Now().Unix()),
	}
	_ bilin.SearchServantServer = &searchServantObj{}
)

type searchServantObj struct {
	httpClient *http.Client
}

type InternalSearchRsp struct {
	Head InternalSearchRspHeader              `json:"responseHeader,omitempty"`
	Data map[string]InternalSearchRspDataItem `json:"response,omitempty"`
}

type InternalSearchRspHeader struct {
	Status  int32  `json:"status,omitempty"` // 响应的状态码，0表示正常。
	Qtime   int32  `json:"QTime,omitempty"`  // 查询耗时，单位毫秒。
	ErrDesc string // 状态码不为0时的错误信息。
}

type InternalSearchRspDataItem struct {
	NumFound int32             `json:"numFound,omitempty"`
	Start    int32             `json:"start,omitempty"`
	Docs     []json.RawMessage `json:"docs,omitempty"` // delay parsing until we know the search type
	Error    string            `json:"error,omitempty"`
}

func NewSearchServantObj() *searchServantObj {
	httpDialer := &net.Dialer{
		Timeout: 5 * time.Second,
	}
	httpTransport := &http.Transport{
		DialContext:       httpDialer.DialContext,
		DisableKeepAlives: false,
	}
	httpClient := &http.Client{
		Transport: httpTransport,
		Timeout:   5 * time.Second,
	}
	return &searchServantObj{
		httpClient: httpClient,
	}
}

func (this *searchServantObj) Search(ctx context.Context, req *bilin.SearchReq) (rsp *bilin.SearchRsp, oerr error) {
	const FuncName = "Search: "
	var (
		request  *http.Request
		response *http.Response
		err      error
		body     []byte
	)

	rsp = &bilin.SearchRsp{
		Head: &bilin.SearchRspHeader{},
	}

	if req.Q == "" {
		rsp.Head.Status = -1
		rsp.Head.ErrDesc = "what do you want to search?"
		return
	}

	if req.Typ == 0 {
		rsp.Head.Status = -1
		rsp.Head.ErrDesc = "which type do you want to search?"
		return
	}

	if request, err = http.NewRequest("GET", config.GetAppConfig().SearchURL, nil); err != nil {
		rsp.Head.Status = -1
		rsp.Head.ErrDesc = err.Error()
		log.Error(FuncName+"http NewRequest fail", zap.Error(err))
		return
	}

	query := request.URL.Query()
	query.Add("q", req.Q)
	if req.Rows <= 0 {
		query.Add("rows", "10")
	} else if req.Rows > 100 {
		rsp.Head.Status = -1
		rsp.Head.ErrDesc = "too many rows!"
		return
	} else {
		query.Add("rows", strconv.FormatInt(int64(req.Rows), 10))
	}
	if req.Start <= 0 {
		query.Add("start", "0")
	} else {
		query.Add("start", strconv.FormatInt(int64(req.Start), 10))
	}
	query.Add("typ", strconv.FormatInt(int64(req.Typ), 10))
	query.Add("app", "51")
	query.Add("v", "0")
	if req.Uid == "" {
		req.Uid = "0"
	}
	query.Add("uid", req.Uid)
	request.URL.RawQuery = query.Encode()

	log.Info(FuncName+"http GET", zap.String("url", request.URL.String()))

	if response, err = this.httpClient.Do(request); err != nil {
		rsp.Head.Status = -1
		rsp.Head.ErrDesc = err.Error()
		log.Error(FuncName+"http client do request fail", zap.Error(err))
		return
	}
	defer response.Body.Close()

	if body, err = ioutil.ReadAll(response.Body); err != nil {
		rsp.Head.Status = -1
		rsp.Head.ErrDesc = err.Error()
		log.Error(FuncName+"http read response body fail", zap.Error(err))
		return
	}

	log.Info(FuncName+"http GET response", zap.String("res", string(body)))

	irsp := &InternalSearchRsp{}
	if err := json.Unmarshal(body, irsp); err != nil {
		rsp.Head.Status = -1
		rsp.Head.ErrDesc = err.Error()
		log.Error(FuncName+"json unmarshal response body fail", zap.Error(err))
		return
	}

	// 用户搜索限定在好友范围，不用搜索引擎的结果。
	var (
		userDocs  []string
		userFound int32
	)
	if req.Typ == bilin.SearchType_USER || req.Typ == bilin.SearchType_USER_ROOM {
		userDocs, userFound = searchBilinUser(req.Q, req.Rows, req.Start, req.Uid)
	}

	rsp.Head.Status = irsp.Head.Status
	rsp.Head.Qtime = irsp.Head.Qtime
	rsp.Head.ErrDesc = irsp.Head.ErrDesc
	rsp.Data = make(map[string]*bilin.SearchRspDataItem, len(irsp.Data))
	for k, v := range irsp.Data {
		numFound := v.NumFound
		start := v.Start
		docs := make([]string, len(v.Docs))
		more := req.Rows == int32(len(v.Docs))
		for i := range v.Docs {
			docs[i] = string(v.Docs[i])
		}
		switch k {
		case "user":
			numFound = userFound
			start = req.Start
			docs = userDocs
			more = req.Rows == int32(len(userDocs))
		case "room":
			docs = extendRoom(docs)
		case "song":
			docs = extendSong(docs)
		}
		item := &bilin.SearchRspDataItem{
			NumFound: numFound,
			Start:    start,
			Docs:     docs,
			Error:    v.Error,
			More:     more,
		}
		rsp.Data[k] = item
	}

	return
}

func containsIgnoreCase(s, substr string) bool {
	s, substr = strings.ToUpper(s), strings.ToUpper(substr)
	return strings.Contains(s, substr)
}

func removeDuplicates(elements []string) (result []string) {
	// Use map to record duplicates as we find them.
	encountered := make(map[string]bool)

	for _, v := range elements {
		if encountered[v] {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[v] = true
			// Append to result slice.
			result = append(result, v)
		}
	}
	// Return the new slice.
	return
}

func searchBilinUser(query string, rows, start int32, uid string) (docs []string, numFound int32) {
	// 精确搜索比邻ID
	var (
		bilinID int64
		myID    int64
		err     error
		userids []string
		users   []e.UserE
		matches []e.UserE
		max     = 1000
	)
	if bilinID, err = strconv.ParseInt(query, 10, 64); err == nil {
		if result, found := clientcenter.GetUserByBLId(bilinID); found {
			userids = append(userids, strconv.FormatInt(result.UserID, 10))
		}
	}
	// 模糊搜索好友列表
	if myID, err = strconv.ParseInt(uid, 10, 64); err == nil {
		for _, friend := range clientcenter.QueryAttentionList(myID) {
			userids = append(userids, strconv.FormatInt(friend.UserID, 10))
		}
	}
	// 保护一下，防止好友数量太多
	if len(userids) > max {
		userids = userids[:max]
	}
	// 去重，精确匹配放在最前面
	userids = removeDuplicates(userids)
	// 补全用户信息
	users = extendUserT(userids)
	// 匹配
	for _, user := range users {
		if containsIgnoreCase(user.BilinId, query) || containsIgnoreCase(user.Name, query) {
			matches = append(matches, user)
		}
	}
	numFound = int32(len(matches))
	// 分页
	if rows <= 0 {
		rows = 10
	}
	if start <= 0 {
		start = 0
	}
	if start > numFound {
		start = numFound
	}
	if start+rows > numFound {
		rows = numFound - start
	}
	matches = matches[start : start+rows]
	for _, x := range matches {
		jsonBytes, _ := json.Marshal(x)
		docs = append(docs, string(jsonBytes))
	}
	return
}

func getUserInfo(u *userinfocenter.UserInfo) *userinfocenter.UserInfo {
	if u != nil {
		return u
	} else {
		return emptyUserInfo
	}
}

func getRoomInfo(r *bilin.RoomInfo) *bilin.RoomInfo {
	if r != nil {
		return r
	} else {
		return emptyRoomInfo
	}
}

func extendUserT(users []string) (results []e.UserE) {
	userids := make([]uint64, len(users))
	for i, x := range users {
		id, err := strconv.ParseUint(x, 10, 64)
		if err != nil {
			log.Error("invalid user id", zap.Error(err), zap.String("uid", x))
			id = 0
		}
		userids[i] = id
	}
	blidmap, err := clientcenter.UserInfoClient().BatchUserBiLinId(context.Background(), &userinfocenter.BatchUserBiLinIdReq{
		Uid: userids,
	})
	if err != nil {
		log.Error("call BatchUserBiLinId fail", zap.Error(err), zap.Any("ids", userids))
		return
	}
	userInfos, err := clientcenter.TakeUserInfo(userids)
	if err != nil {
		log.Error("can not get user info by id", zap.Error(err), zap.Any("ids", userids))
		return
	}
	exUsers := make([]e.UserE, 0, len(users))
	for i, x := range users {
		info := getUserInfo(userInfos[userids[i]])
		age := (time.Now().Unix() - info.Birthday) / (365 * 24 * 60 * 60)
		if age < 18 {
			age = 18
		}

		var exUser e.UserE
		exUser.Id = x
		exUser.BilinId = strconv.FormatUint(blidmap.Uid2Bilinid[userids[i]], 10)
		exUser.Name = info.NickName
		exUser.Avatar = info.Avatar
		exUser.Sex = strconv.FormatUint(uint64(info.Showsex), 10)
		exUser.Age = strconv.FormatInt(age, 10)
		exUser.Location = info.City
		exUser.Live = "0"
		exUser.RoomUserNum = "0"
		exUsers = append(exUsers, exUser)
	}
	results = exUsers
	return
}

func extendUser(userDocs []string) (docs []string) {
	docs = make([]string, len(userDocs))
	copy(docs, userDocs)
	users := make([]u.UserU, len(userDocs))
	for i := range userDocs {
		json.Unmarshal([]byte(userDocs[i]), &users[i])
	}
	userids := make([]uint64, len(users))
	for i, x := range users {
		id, err := strconv.ParseUint(x.Id, 10, 64)
		if err != nil {
			log.Error("invalid user id", zap.Error(err), zap.String("uid", x.Id))
			id = 0
		}
		userids[i] = id
	}
	userInfos, err := clientcenter.TakeUserInfo(userids)
	if err != nil {
		log.Error("can not get user info by id", zap.Error(err), zap.Any("ids", userids))
		return
	}
	docs = docs[:0]
	for i, x := range users {
		info := getUserInfo(userInfos[userids[i]])
		age := (time.Now().Unix() - info.Birthday) / (365 * 24 * 60 * 60)
		if age < 18 {
			age = 18
		}

		var exUser e.UserE
		exUser.Id = x.Id
		exUser.BilinId = x.BilinId
		exUser.Name = x.Name
		exUser.Avatar = info.Avatar
		exUser.Sex = strconv.FormatUint(uint64(info.Showsex), 10)
		exUser.Age = strconv.FormatInt(age, 10)
		exUser.Location = info.City
		exUser.Live = "0"
		exUser.RoomUserNum = "0"

		jsonBytes, _ := json.Marshal(exUser)
		docs = append(docs, string(jsonBytes))
	}
	return
}

func extendRoom(roomDocs []string) (docs []string) {
	docs = make([]string, len(roomDocs))
	copy(docs, roomDocs)
	rooms := make([]u.RoomU, len(roomDocs))
	for i := range roomDocs {
		json.Unmarshal([]byte(roomDocs[i]), &rooms[i])
	}
	roomids := make([]uint64, len(rooms))
	for i, x := range rooms {
		id, err := strconv.ParseUint(x.Id, 10, 64)
		if err != nil {
			log.Error("invalid room id", zap.Error(err), zap.String("roomid", x.Id))
			id = 0
		}
		roomids[i] = id
	}
	roomInfos, err := clientcenter.TakeRoomInfo(roomids)
	if err != nil {
		log.Error("can not get room info by id", zap.Error(err), zap.Any("ids", roomids))
		return
	}
	var badges []*bilin.UserBabgeInfo
	babgeResp, err := clientcenter.ConfClient().BatchUserBabge(context.Background(), &bilin.UserBabgeReq{})
	if err != nil {
		log.Error("call BatchUserBabge fail", zap.Error(err))
	} else {
		badges = babgeResp.Userbabgeinfo
	}
	exRooms := make([]e.RoomE, len(rooms))
	owners := make([]uint64, len(rooms))
	for i, x := range rooms {
		info := getRoomInfo(roomInfos[roomids[i]])
		owners[i] = info.Owner

		exRooms[i].Id = x.Id
		exRooms[i].Name = x.Name
		exRooms[i].Live = x.Live
		/*if info.RoomType2 == officailRoom {
			exRooms[i].DisplayId = x.Id
		} else {
			exRooms[i].DisplayId = strconv.FormatUint(info.OwnerBilinID, 10)
		}*/
		exRooms[i].DisplayId = x.DisplayId
		exRooms[i].StartAt = strconv.FormatUint(info.Starttime, 10)
		exRooms[i].UserNum = strconv.FormatUint(info.Usernumber, 10)
	}
	ownerInfos, err := clientcenter.TakeUserInfo(owners)
	if err != nil {
		log.Error("can not get room owner info by id", zap.Error(err), zap.Any("ids", owners))
		return
	}
	docs = docs[:0]
	for i := range rooms {
		info := getUserInfo(ownerInfos[owners[i]])

		exRooms[i].Avatar = info.Avatar
		if badges != nil {
			for _, badge := range badges {
				if badge.Userid == info.Uid {
					exRooms[i].TagURL = append(exRooms[i].TagURL, badge.Url)
				}
			}
		}

		if exRooms[i].UserNum == "0" {
			continue
		}

		jsonBytes, _ := json.Marshal(exRooms[i])
		docs = append(docs, string(jsonBytes))
	}
	return
}

func extendSong(songDocs []string) (docs []string) {
	docs = make([]string, len(songDocs))
	copy(docs, songDocs)
	songs := make([]u.SongU, len(songDocs))
	for i := range songDocs {
		json.Unmarshal([]byte(songDocs[i]), &songs[i])
	}
	songids := make([]int64, len(songs))
	for i, x := range songs {
		id, err := strconv.ParseInt(x.Id, 10, 64)
		if err != nil {
			log.Error("invalid song id", zap.Error(err), zap.String("songid", x.Id))
			id = 0
		}
		songids[i] = id
	}
	docs = docs[:0]
	for i, x := range songs {
		ktv, found := clientcenter.GetBilinKtvById(songids[i])
		if !found {
			continue
		}

		var exSong e.SongE
		exSong.Id = x.Id
		exSong.Name = ktv.Name
		exSong.Artist = ktv.Artist
		exSong.Duration = strconv.FormatInt(ktv.Duration, 10)
		exSong.UploadBy = ktv.UploadBy
		exSong.Pkg = ktv.Pkg
		exSong.PkgMd5 = ktv.PkgMd5
		exSong.PkgLen = strconv.FormatInt(ktv.PkgLen, 10)

		jsonBytes, _ := json.Marshal(exSong)
		docs = append(docs, string(jsonBytes))
	}
	return
}

func (this *searchServantObj) GetRelatedHotSearches(ctx context.Context, req *bilin.GetRelatedHotSearchesReq) (rsp *bilin.GetRelatedHotSearchesRsp, oerr error) {
	const FuncName = "GetRelatedHotSearches: "
	var (
		request  *http.Request
		response *http.Response
		err      error
		body     []byte
		keywords []string
	)

	rsp = &bilin.GetRelatedHotSearchesRsp{
		HotSearches: []string{},
	}

	if req.Q == "" {
		return
	}

	if req.Typ == 0 {
		return
	}

	if request, err = http.NewRequest("GET", config.GetAppConfig().SearchURL, nil); err != nil {
		log.Error(FuncName+"http NewRequest fail", zap.Error(err))
		return
	}

	query := request.URL.Query()
	query.Add("q", req.Q)
	if req.Rows <= 0 {
		query.Add("rows", "10")
	} else if req.Rows > 100 {
		return
	} else {
		query.Add("rows", strconv.FormatInt(int64(req.Rows), 10))
	}
	if req.Start <= 0 {
		query.Add("start", "0")
	} else {
		query.Add("start", strconv.FormatInt(int64(req.Start), 10))
	}
	query.Add("typ", strconv.FormatInt(int64(req.Typ), 10))
	query.Add("app", "51")
	query.Add("v", "0")
	if req.Uid == "" {
		req.Uid = "0"
	}
	query.Add("uid", req.Uid)
	request.URL.RawQuery = query.Encode()

	log.Info(FuncName+"http GET", zap.String("url", request.URL.String()))

	if response, err = this.httpClient.Do(request); err != nil {
		log.Error(FuncName+"http client do request fail", zap.Error(err))
		return
	}
	defer response.Body.Close()

	if body, err = ioutil.ReadAll(response.Body); err != nil {
		log.Error(FuncName+"http read response body fail", zap.Error(err))
		return
	}

	log.Info(FuncName+"http GET response", zap.String("res", string(body)))

	irsp := &InternalSearchRsp{}
	if err := json.Unmarshal(body, irsp); err != nil {
		log.Error(FuncName+"json unmarshal response body fail", zap.Error(err))
		return
	}

	for k, v := range irsp.Data {
		switch k {
		case "room", "song":
			for i := range v.Docs {
				obj := &struct {
					Name string `json:"name"`
				}{}
				if err := json.Unmarshal(v.Docs[i], obj); err != nil {
					log.Error(FuncName+"json unmarshal doc fail", zap.Error(err))
				} else if obj.Name != "" {
					keywords = append(keywords, obj.Name)
				}
			}
		}
	}
	rsp.HotSearches = removeDuplicates(keywords)
	return
}

func (this *searchServantObj) GetAllHotSearches(ctx context.Context, req *bilin.GetAllHotSearchesReq) (rsp *bilin.GetAllHotSearchesRsp, oerr error) {
	const FuncName = "GetAllHotSearches: "

	rsp = &bilin.GetAllHotSearchesRsp{
		HotSearches: []string{
			"交友",
			"开黑吃鸡",
			"处对象",
			"听歌",
			"随便聊聊",
			"陪玩",
		},
	}
	return
}

func (this *searchServantObj) GetHotSongs(ctx context.Context, req *bilin.GetHotSongsReq) (rsp *bilin.GetHotSongsRsp, oerr error) {
	const FuncName = "GetHotSongs: "

	var songs []u.SongU
	for i := 1; i <= 150; i++ {
		songs = append(songs, u.SongU{Id: strconv.FormatInt(int64(i), 10)})
	}
	docs := make([]string, 0, len(songs))
	for _, s := range songs {
		doc, _ := json.Marshal(s)
		docs = append(docs, string(doc))
	}

	if req.Rows <= 0 {
		req.Rows = 10
	}
	if req.Start <= 0 {
		req.Start = 0
	}

	if req.Start > int32(len(docs)) {
		req.Start = int32(len(docs))
	}
	temp := req.Rows
	if req.Start+temp > int32(len(docs)) {
		temp = int32(len(docs)) - req.Start
	}

	docs = docs[req.Start : req.Start+temp]
	more := req.Rows == int32(len(docs))
	docs = extendSong(docs)
	rsp = &bilin.GetHotSongsRsp{
		Head: &bilin.SearchRspHeader{},
		Data: map[string]*bilin.SearchRspDataItem{
			"song": {
				NumFound: int32(len(songs)),
				Start:    req.Start,
				Docs:     docs,
				More:     more,
			},
		},
	}
	return
}
