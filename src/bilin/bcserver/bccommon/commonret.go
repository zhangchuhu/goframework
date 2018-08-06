// @author kordenlu
// @创建时间 2018/03/29 16:15
// 功能描述:

package bccommon

import (
	"bilin/protocol"
	strconv "strconv"
	"strings"
	"reflect"
)

func convertToString(b []byte) string {
	s := make([]string, len(b))
	for i := range b {
		s[i] = strconv.Itoa(int(b[i]))
	}
	return strings.Join(s, ",")
}

func Contains(slices []uint64, uid uint64) bool {
	for _, item := range slices {
		if item == uid {
			return true
		}
	}
	return false
}

func IsNil(i interface{}) bool {
	vi := reflect.ValueOf(i)
	if vi.Kind() == reflect.Ptr {
		return vi.IsNil()
	}
	return false
}

var (
	COMMERRORILLEGALDESC = "当前操作无效"
	ErrCodeAndDescMap    = map[bilin.CommonRetInfo_RETCODE]string{
		//公共错误
		bilin.CommonRetInfo_ILLEGAL_MESSAGE: "当前操作无效",

		//直播间错误
		bilin.CommonRetInfo_ENTER_ROOM_FAILED:          "服务器开小差了，再试试呗~",
		bilin.CommonRetInfo_ENTER_USER_NO_RIGHT:        "房主拒绝你进来哦~",
		bilin.CommonRetInfo_ENTER_ROOM_NOT_START:       "直播异常，主持人未开播",
		bilin.CommonRetInfo_ENTER_BAD_NETWORK:          "服务器开小差了，再试试呗~",
		bilin.CommonRetInfo_ENTER_ROOM_CLOSED:          "直播间已被关闭",
		bilin.CommonRetInfo_ENTER_ROOM_LOCKED:          "直播间被锁",
		bilin.CommonRetInfo_ENTER_ROOM_PWDERR:          "房间已上锁，请输入正确密码",
		bilin.CommonRetInfo_ENTER_ROOM_ALREADY_IN_ROOM: "用户已经在房间",
		bilin.CommonRetInfo_ENTER_ROOM_FORBIDDEN:       "直播间涉嫌违规",
	}
)

var (
	SUCCESSMESSAGE = &bilin.CommonRetInfo{
		Ret:  0,
		Desc: "成功",
	}

	UserDefinedFailed = func(ret bilin.CommonRetInfo_RETCODE, desc string) *bilin.CommonRetInfo {
		return &bilin.CommonRetInfo{
			Ret:  ret,
			Desc: desc,
			Show: true,
		}
	}
)

func SuccessOrFailedFun(ret int64) (code bool) {
	switch ret {
	case 0:
		return true
	default:
		return false
	}
}
