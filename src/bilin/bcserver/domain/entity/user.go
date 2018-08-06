// @author kordenlu
// @创建时间 2018/03/29 16:20
// 功能描述:

package entity

//角色定义
const (
	ROLE_ADMIN    = 1
	ROLE_HOST     = 2
	ROLE_AUDIENCE = 3
)

const (
	// StatusUserOnComing 老的bc_server对用户进频道有两个状态，发第一个enter请求的时候状态为coming
	StatusUserOnComing = 0
	// StatusUserJoined 发init请求的时候，状态就改为joined。新系统合成一个请求了，那就直接用USER_JOINED
	StatusUserJoined = 1
)

// User 序列化之后，存在redis的直播间hash里
type User struct {
	RoomID         uint64 `json:"roomid"`
	UserID         uint64 `json:"user_id"`
	Role           uint32 `json:"role"`
	Status         int    `json:"status"`
	BeginJoinTime  uint64 `json:"begin_join_timestamp"`
	NickName       string `json:"nick"`
	AvatarURL      string `json:"head_url"`
	IsMuted        uint32 `json:"is_muted"`
	EnterBeginTime uint64 `json:"enter_begin_time"`
	LinkBeginTime  int    `json:"link_begin_time"`
	Sex            int32  `json:"sex"`
	Age            int32  `json:"age"`
	CityName       string `json:"cityName"`
	PraiseCount    uint32 `json:"praise_count"`
	MikeIndex      uint32 `json:"mike_index"`
	FansCount      uint32 `json:"fans_count"`
	OnMikeTime     uint64 `json:"on_mike_timestamp"`
	Version        string `json:"version"`
	Signature      string `json:"signature"`
}

type UserSortByBeginJoinTimeSlice []*User

func (c UserSortByBeginJoinTimeSlice) Len() int {
	return len(c)
}
func (c UserSortByBeginJoinTimeSlice) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c UserSortByBeginJoinTimeSlice) Less(i, j int) bool {
	return c[i].BeginJoinTime > c[j].BeginJoinTime
}

type UserSortByOnMikeTimeSlice []*User

func (c UserSortByOnMikeTimeSlice) Len() int {
	return len(c)
}
func (c UserSortByOnMikeTimeSlice) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c UserSortByOnMikeTimeSlice) Less(i, j int) bool {
	return c[i].OnMikeTime < c[j].OnMikeTime
}

func (user *User) Less(other *User) bool {
	if user.BeginJoinTime < other.BeginJoinTime {
		return true
	}

	return false
}
