package service

import (
	"bilin/protocol"
	"testing"
)

var (
	karaoke_roomid uint64 = 400000367

	user1 = &bilin.UserInfo{
		Userid:    111,
		Nick:      "test1",
		Avatarurl: "avatarurl1",
		Fanscount: 0,
		From:      bilin.USERFROM_BROADCAST,
		Mute:      0,
		Sex:       0,
		Age:       0,
		CityName:  "guangzhou",
		Signature: "确认过眼神，我遇上对的人~",
	}
	song1 = &bilin.KaraokeSongInfo{Id: "aaaaa", SongName: "醉赤壁", Userinfo: user1, Status: bilin.KaraokeSongInfo_PREPARE}

	user2 = &bilin.UserInfo{
		Userid:    222,
		Nick:      "test2",
		Avatarurl: "avatarurl2",
		Fanscount: 0,
		From:      bilin.USERFROM_BROADCAST,
		Mute:      0,
		Sex:       0,
		Age:       0,
		CityName:  "guangzhou",
		Signature: "忘了是怎么开始，也许就是对你，有一种感觉~",
	}
	song2 = &bilin.KaraokeSongInfo{Id: "bbbbb", SongName: "爱很简单", Userinfo: user2, Status: bilin.KaraokeSongInfo_PREPARE}

	user3 = &bilin.UserInfo{
		Userid:    333,
		Nick:      "test3",
		Avatarurl: "avatarurl3",
		Fanscount: 0,
		From:      bilin.USERFROM_BROADCAST,
		Mute:      0,
		Sex:       0,
		Age:       0,
		CityName:  "guangzhou",
		Signature: "我的爱如潮水，爱如潮水将你推~",
	}
	song3 = &bilin.KaraokeSongInfo{Id: "ccccc", SongName: "爱如潮水", Userinfo: user3, Status: bilin.KaraokeSongInfo_PREPARE}

	user4 = &bilin.UserInfo{
		Userid:    444,
		Nick:      "test4",
		Avatarurl: "avatarurl4",
		Fanscount: 0,
		From:      bilin.USERFROM_BROADCAST,
		Mute:      0,
		Sex:       0,
		Age:       0,
		CityName:  "guangzhou",
		Signature: "葫芦娃，葫芦娃，一根藤上七个瓜~",
	}
	song4 = &bilin.KaraokeSongInfo{Id: "ddddd", SongName: "葫芦娃", Userinfo: user4, Status: bilin.KaraokeSongInfo_PREPARE}
)

func TestRedisCloseKaraoke(t *testing.T) {
	err := RedisCloseKaraoke(karaoke_roomid)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestRedisAddSong(t *testing.T) {
	err := RedisAddSong(karaoke_roomid, song1)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestRedisDelSong(t *testing.T) {
	err := RedisDelSong(karaoke_roomid, song1)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestRedisGetSongsList(t *testing.T) {
	result, err := RedisGetSongsList(karaoke_roomid)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(result)
}

func TestRedisGetSongInfoByID(t *testing.T) {
	result, err := RedisGetSongInfoByID(karaoke_roomid, "ddddd")
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(result)
}

func TestRedisGetDisplaySong(t *testing.T) {
	result, err := RedisGetDisplaySong(karaoke_roomid)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(result)
}

func TestRedisSetDisplaySong(t *testing.T) {
	result, err := RedisGetDisplaySong(karaoke_roomid)
	if err != nil {
		t.Error(err)
		return
	}

	result.Status = bilin.KaraokeSongInfo_SINGING

	err = RedisChangeDisplaySongStatus(karaoke_roomid, result)
	if err != nil {
		t.Error(err)
		return
	}

}

func TestRedisAddSongBeforePrepareSong(t *testing.T) {
	result, err := RedisGetDisplaySong(karaoke_roomid)
	if err != nil {
		t.Error(err)
		return
	}

	err = RedisAddSongBeforePrepareSong(karaoke_roomid, result, song2)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(result)
}

func TestRedisAddSongAfterSingSong(t *testing.T) {
	result, err := RedisGetDisplaySong(karaoke_roomid)
	if err != nil {
		t.Error(err)
		return
	}

	err = RedisAddSongAfterSingSong(karaoke_roomid, result, song3)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(result)
}
