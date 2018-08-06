package collector

import (
	"bilin/bcserver/domain/entity"
	"bilin/bcserver/domain/service"
	"bilin/protocol"
	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"github.com/rs/xid"
)

func GenXid() string {
	id := xid.New()
	return id.String()
}

func AllRoomKaraokeInfo(room *entity.Room) *bilin.AllRoomKaraokeInfo {
	return &bilin.AllRoomKaraokeInfo{
		Songs: GetSongsList(room),
	}

}

func CloseKaraoke(room *entity.Room) (err error) {
	const prefix = "CloseKaraoke "

	service.RedisCloseKaraoke(room.Roomid)
	room.Karaokeswitch = bilin.BaseRoomInfo_CLOSEKARAOKE
	log.Debug(prefix, zap.Any("roomid", room.Roomid))
	return
}

func OpenKaraoke(room *entity.Room) (err error) {
	const prefix = "OpenKaraoke "

	room.Karaokeswitch = bilin.BaseRoomInfo_OPENKARAOKE

	log.Debug(prefix, zap.Any("roomid", room.Roomid))
	return
}

//func GetPrepareSong(room *entity.Room) (prepareSong *bilin.KaraokeSongInfo) {
//	const prefix = "GetPrepareSong "
//
//	prepareSong, _ = service.RedisGetDisplaySong(room.Roomid)
//	if prepareSong == nil || prepareSong.Status != bilin.KaraokeSongInfo_PREPARE {
//		log.Info(prefix+"nil", zap.Any("roomid", room.Roomid), zap.Any("prepareSong", prepareSong))
//		return
//	}
//
//	log.Debug(prefix, zap.Any("roomid", room.Roomid), zap.Any("prepareSong", prepareSong))
//	return
//}
//
//func GetSingingSong(room *entity.Room) (singingSong *bilin.KaraokeSongInfo) {
//	const prefix = "GetSingingSong "
//
//	singingSong, _ = service.RedisGetDisplaySong(room.Roomid)
//	if singingSong == nil || singingSong.Status != bilin.KaraokeSongInfo_SINGING {
//		log.Info(prefix+"nil", zap.Any("roomid", room.Roomid), zap.Any("singingSong", singingSong))
//		return
//	}
//
//	log.Debug(prefix, zap.Any("roomid", room.Roomid), zap.Any("singingSong", singingSong))
//	return
//}

func GetDisplaySong(room *entity.Room) (displaysong *bilin.KaraokeSongInfo) {
	const prefix = "GetDisplaySong "

	displaysong, _ = service.RedisGetDisplaySong(room.Roomid)

	log.Debug(prefix, zap.Any("roomid", room.Roomid), zap.Any("displaysong", displaysong))
	return
}

func ChangeDisplaySongStatus(room *entity.Room, displaysong *bilin.KaraokeSongInfo) (err error) {
	const prefix = "ChangeDisplaySongStatus "

	err = service.RedisChangeDisplaySongStatus(room.Roomid, displaysong)

	log.Debug(prefix, zap.Any("roomid", room.Roomid), zap.Any("displaysong", displaysong))
	return
}

func GetSongsCount(roomid uint64)(count int64) {
	const prefix = "GetSongsCount "

	count,_= service.RedisGetSongsCount(roomid)

	log.Debug(prefix, zap.Any("roomid", roomid), zap.Any("count", count))
	return
}

func GetSongsList(room *entity.Room) (songsList []*bilin.KaraokeSongInfo) {
	const prefix = "GetSongsList "

	songsList, _ = service.RedisGetSongsList(room.Roomid)

	log.Debug(prefix, zap.Any("roomid", room.Roomid), zap.Any("songsList", songsList))
	return
}

func GetSongInfoByID(room *entity.Room, id string) (song *bilin.KaraokeSongInfo) {
	const prefix = "GetSongInfoByID "

	song, _ = service.RedisGetSongInfoByID(room.Roomid, id)

	log.Debug(prefix, zap.Any("roomid", room.Roomid), zap.Any("song", song))
	return
}

func AddSong(room *entity.Room, song *bilin.KaraokeSongInfo) (err error) {
	const prefix = "AddSong "

	err = service.RedisAddSong(room.Roomid, song)

	log.Debug(prefix, zap.Any("roomid", room.Roomid), zap.Any("song", song))
	return
}

func DelSong(room *entity.Room, song *bilin.KaraokeSongInfo) (err error) {
	const prefix = "DelSong "

	err = service.RedisDelSong(room.Roomid, song)

	log.Debug(prefix, zap.Any("roomid", room.Roomid), zap.Any("song", song))
	return
}

func SetTopSong(room *entity.Room, song *bilin.KaraokeSongInfo) (err error) {
	const prefix = "SetTopSong "

	displaySong := GetDisplaySong(room)
	if displaySong == nil {
		err = service.RedisAddSong(room.Roomid, song)
		log.Info(prefix, zap.Any("roomid", room.Roomid), zap.Any("song", song))
		return
	}

	if displaySong.Id == song.Id {
		log.Info(prefix+"equal id, no need to set top", zap.Any("roomid", room.Roomid), zap.Any("song", song))
		return
	}
	//先从redis中删除数据，然后再插入
	err = service.RedisDelSong(room.Roomid, song)
	if err != nil {
		log.Error(prefix+"RedisDelSong", zap.Any("song", song))
		return
	}

	if displaySong.Status == bilin.KaraokeSongInfo_PREPARE { //如果当前没有正在播放的歌曲，则插入到第0位
		err = service.RedisAddSongBeforePrepareSong(room.Roomid, displaySong, song)
	} else { //如果当前有正在播放或者暂停的歌曲，则插入到第1位
		err = service.RedisAddSongAfterSingSong(room.Roomid, displaySong, song)
	}

	log.Debug(prefix, zap.Any("roomid", room.Roomid), zap.Any("song", song))
	return
}
