package handler

import ()

/*
func TestRecommentList(t *testing.T) {
	room := LivingRoom{}
	room.RoomId = 1111
	room.UserID = 2222
	room.NickName = "bobo"
	room.Count = 666

	rec := Carousel{}
	rec.CarouselId = 1
	rec.CarouselType = 2
	rec.BackgroudURL = "xxx"
	rec.TargetType = 3
	rec.TargetURL = "xxx"
	rec.Room = room

	list := []*Carousel{}
	list = append(list, &rec)

	resp := &HttpRetComm{
		IsEncrypt: "false",
		Data: HttpRetDataComm{
			Code: 0,
			Msg:  "success",
			Body: CarouselListBody{
				CarouselList: list,
			},
		},
	}

	byte, err := json.Marshal(resp)
	if err != nil {
		t.Error("failed to marshal", err)
	}
	t.Logf("byte:%v", string(byte))
}

func TestTodayRank(t *testing.T) {

	var user RankUser

	user.UserID = 1
	user.NickName = "url"
	user.Avatar = "gray avatar"

	var users []*RankUser
	users = append(users, &user)

	rank := RankInfo{
		Users:     users,
		TargetURL: "xxx",
	}

	today_rank := TadayRank{
		ContributeRank: rank,
		AnchorRank:     rank,
		GuardRank:      rank,
	}

	resp := &HttpRetComm{
		IsEncrypt: "false",
		Data: HttpRetDataComm{
			Code: 0,
			Msg:  "success",
			Body: TadayRankBody{
				TadayRank: today_rank,
			},
		},
	}

	byte, err := json.Marshal(resp)
	if err != nil {
		t.Error("failed to marshal", err)
	}
	t.Logf("byte:%v", string(byte))
}
*/
/*


type RankUser struct {
	UserID   uint64 `json:"user_id"`
	UserID string `json:"nick_name"`
	Avatar   string `json:"avatar"`
}

type RankInfo struct {
	Users      []*RankUser `json:"users"`
	TargetURL string      `json:"target_url"`
}

type TadayRank struct {
	ReachRank  RankInfo `json:"reach_rank"`  //土豪榜
	AnchorRank RankInfo `json:"anchor_rank"` //主播榜
	GuardRank  RankInfo `json:"guard_rank"`  //守护榜
}
*/
