package service

import (
	"bilin/bcserver/config"
	"bilin/bcserver/domain/entity"
	"bilin/protocol"
	"testing"
	"time"
)

var (
	roomid    uint64 = 123456
	userid    uint64 = 567891
	now              = time.Now().Unix()
	appconfig        = &config.AppConfig{
		//测试环境redis
		RedisAddr: "183.36.122.50:4019",

		//线上redis,谨慎使用
		//RedisAddr: "221.228.79.78:4000",
	}
)

func init() {
	config.SetTestAppConfig(appconfig)
	RedisInit()
}

func TestRedisIfRoomExist(t *testing.T) {
	ret, err := RedisIfRoomExist(roomid)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(ret)
}

func TestRedisAddRoom(t *testing.T) {
	err := RedisAddRoom(entity.NewRoom(roomid))

	if err != nil {
		t.Error(err)
		return
	}
}

func TestRedisRemoveRoom(t *testing.T) {
	err := RedisRemoveRoom(roomid)
	if err != nil {
		t.Error(err)
		return
	}

}

func TestRedisGetRoomInfo(t *testing.T) {
	ret, err := RedisGetRoomInfo(roomid)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", ret)
}

// 频道内用户基本进出事件
func TestRedisAddUser(t *testing.T) {
	err := RedisAddUser(roomid, &entity.User{
		RoomID:         roomid,
		UserID:         userid,
		Role:           entity.ROLE_AUDIENCE,
		Status:         entity.StatusUserJoined,
		BeginJoinTime:  uint64(now),
		NickName:       "测试昵称",
		AvatarURL:      "头像测试",
		IsMuted:        0,
		EnterBeginTime: uint64(now),
		LinkBeginTime:  0,
		Sex:            0,
		Age:            22,
		CityName:       "广州市",
		PraiseCount:    0,
	})
	if err != nil {
		t.Error(err)
		return
	}
}

func TestRedisGetUser(t *testing.T) {
	ret, err := RedisGetUser(roomid, userid)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", ret)
}

func TestRedisGetUserCount(t *testing.T) {
	ret, err := RedisGetUserCount(roomid)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", ret)
}

func TestRedisRemoveUser(t *testing.T) {
	err := RedisRemoveUser(roomid, userid)
	if err != nil {
		t.Error(err)
		return
	}

}

func TestRedisGetRoomUserList(t *testing.T) {
	ret, err := RedisGetRoomUserList(roomid)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", ret)
}

func TestRedisGetDisplayedUsers(t *testing.T) {
	ret, err := RedisGetDisplayedUsers(roomid, 40)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", ret)
}

func TestRedisSetPingTime(t *testing.T) {
	ret, err := RedisSetPingTime(roomid, userid)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", ret)
}

// 频道内禁止公屏聊天用户
func TestRedisGetForbidenStatus(t *testing.T) {
	ret, err := RedisGetForbidenStatus(roomid, userid)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", ret)
}

func TestRedisSetForbidenStatus(t *testing.T) {
	err := RedisSetForbidenStatus(roomid, userid, true)
	if err != nil {
		t.Error(err)
		return
	}

}

func TestRedisGetForbidenUserList(t *testing.T) {
	ret, err := RedisGetForbidenUserList(roomid)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", ret)
}

func TestRedisClearForbidenUserList(t *testing.T) {
	err := RedisClearForbidenUserList(roomid)
	if err != nil {
		t.Error(err)
		return
	}

}

// 频道内上麦用户信息，由于需要和老的版本兼容，所以需要三个hash 来存储麦上用户信息
// 分别是   用户id --> 用户信息
//		   麦位id --> 麦位信息 (空/有人/锁住)
func TestRedisAddUserToMike(t *testing.T) {
	err := RedisAddUserToMike(roomid, &entity.User{
		RoomID:         roomid,
		UserID:         userid,
		Role:           entity.ROLE_AUDIENCE,
		Status:         entity.StatusUserJoined,
		BeginJoinTime:  uint64(now),
		NickName:       "测试昵称",
		AvatarURL:      "头像测试",
		IsMuted:        0,
		EnterBeginTime: uint64(now),
		LinkBeginTime:  0,
		Sex:            0,
		Age:            22,
		CityName:       "广州市",
		PraiseCount:    0,
		MikeIndex:      1,
	})
	if err != nil {
		t.Error(err)
		return
	}

}

func TestRedisIfUserOnMike(t *testing.T) {
	ret, err := RedisIfUserOnMike(roomid, userid)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", ret)
}

func TestRedisGetUserOnMike(t *testing.T) {
	ret, err := RedisGetUserOnMike(roomid, userid)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", ret)
}

func TestRedisGetOnMikeUserList(t *testing.T) {
	ret, err := RedisGetOnMikeUserList(roomid)
	if err != nil {
		t.Error(err)
		return
	}

	for _, item := range ret {
		t.Logf("%+v", item)
	}

}

func TestRedisGetOnMikeUserCount(t *testing.T) {
	ret, err := RedisGetOnMikeUserCount(roomid)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", ret)
}

func TestRedisRemoveUserFromMike(t *testing.T) {
	err := RedisRemoveUserFromMike(roomid, &entity.User{
		RoomID:    roomid,
		UserID:    userid,
		MikeIndex: 1,
	})
	if err != nil {
		t.Error(err)
		return
	}

}

//锁定/解锁麦位  opt 2: 锁定   0： 解锁,麦位上有人是不能锁定和解锁的
func TestRedisLockUnlockMikeWheat(t *testing.T) {
	err := RedisLockUnlockMikeWheat(roomid, 1, bilin.MikeInfo_LOCK)
	if err != nil {
		t.Error(err)
		return
	}

}

func TestRedisGetAllMikeWheatStatus(t *testing.T) {
	ret, err := RedisGetAllMikeWheatStatus(roomid)
	if err != nil {
		t.Error(err)
		return
	}

	for key, item := range ret {
		t.Logf("%+v, %+v", key, item)
	}

}

// 排麦用户操作
func TestRedisAddUserToApplyMikeList(t *testing.T) {
	err := RedisAddUserToApplyMikeList(roomid, &entity.User{
		RoomID:         roomid,
		UserID:         userid,
		Role:           entity.ROLE_AUDIENCE,
		Status:         entity.StatusUserJoined,
		BeginJoinTime:  uint64(now),
		NickName:       "测试昵称",
		AvatarURL:      "头像测试",
		IsMuted:        0,
		EnterBeginTime: uint64(now),
		LinkBeginTime:  0,
		Sex:            0,
		Age:            22,
		CityName:       "广州市",
		PraiseCount:    0,
		MikeIndex:      1,
	})
	if err != nil {
		t.Error(err)
		return
	}

}

func TestRedisGetApplyMikeUserCount(t *testing.T) {
	ret, err := RedisGetApplyMikeUserCount(roomid)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", ret)

}

func TestRedisGetApplyMikeUserList(t *testing.T) {
	ret, err := RedisGetApplyMikeUserList(roomid)
	if err != nil {
		t.Error(err)
		return
	}

	for _, item := range ret {
		t.Logf("%+v", item)
	}

}

func TestRedisGetOneApplyMikeUser(t *testing.T) {
	ret, err := RedisGetOneApplyMikeUser(roomid)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", ret)

}

func TestRedisRemoveUserFromApplyMikeList(t *testing.T) {
	err := RedisRemoveUserFromApplyMikeList(roomid, userid)
	if err != nil {
		t.Error(err)
		return
	}

}

func TestRedisClearApplyMikeList(t *testing.T) {
	err := RedisClearApplyMikeList(roomid)
	if err != nil {
		t.Error(err)
		return
	}

}

func TestRedisHsacnRoomList(t *testing.T) {
	var cursor uint64 = 0
	for {
		_, cursor, _ = RedisHscanRoomList(cursor)
		if cursor == 0 {
			break
		}
	}
}

func TestRedisGetUserFansCount(t *testing.T) {
	count, err := RedisGetUserFansCount(17795535)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", count)
}

func TestRedisInitWheatLock(t *testing.T) {
	RedisClearMikeWheat(17795535)
	for idx := 1; idx <= 3; idx++ {
		RedisLockUnlockMikeWheat(17795535, uint32(idx), bilin.MikeInfo_EMPTY)
	}

	maplist, _ := RedisGetAllMikeWheatStatus(17795535)
	t.Logf("%+v", maplist)
}

//online test
func TestRedisChangeRoomMaixuSwitch(t *testing.T) {
	room, err := RedisGetRoomInfo(410298353)
	if err != nil {
		t.Error(err)
		return
	}

	room.Maixuswitch = bilin.BaseRoomInfo_OPENMAIXU
	RedisAddRoom(room)

	t.Log(room)
}
