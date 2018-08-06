package service

import (
	"bilin/protocol"
	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
)

const (
	RedisKaraokeKeyPrefix = "karaoke_songs_"
)

func RedisCloseKaraoke(roomid uint64) (err error) {
	const prefix = "RedisCloseKaraoke "

	redisKey := RedisKaraokeKeyPrefix + fmt.Sprintf("%d", roomid)

	//清空歌曲列表
	if err = RedisClient.Del(redisKey).Err(); err != nil {
		log.Error(prefix+"redis.Del", zap.Any("err", err), zap.Any("redisKey", redisKey))
	}

	log.Info(prefix+"end", zap.Any("roomid", roomid))
	return
}

//插到队列尾部
func RedisAddSong(roomid uint64, song *bilin.KaraokeSongInfo) (err error) {
	const prefix = "RedisAddSong "

	redisKey := RedisKaraokeKeyPrefix + fmt.Sprintf("%d", roomid)

	songBytes, err := json.Marshal(song)
	if err != nil {
		log.Error(prefix+"json.Marshal(song)", zap.Any("err", err))
		return
	}

	if err = RedisClient.RPush(redisKey, string(songBytes)).Err(); err != nil {
		log.Error(prefix+"redis.RPush", zap.Any("err", err))
		return
	}

	log.Info(prefix+"end", zap.Any("roomid", roomid), zap.Any("song", song))
	return
}

func RedisGetSongsCount(roomid uint64) (count int64, err error) {
	const prefix = "RedisGetSongsCount "

	redisKey := RedisKaraokeKeyPrefix + fmt.Sprintf("%d", roomid)
	if count, err = RedisClient.LLen(redisKey).Result(); err != nil {
		log.Error(prefix+"redis.LLen", zap.Any("err", err))
	}

	log.Debug(prefix, zap.Any("roomid", roomid), zap.Any("count", count))
	return
}

func RedisDelSong(roomid uint64, song *bilin.KaraokeSongInfo) (err error) {
	const prefix = "RedisDelSong "

	redisKey := RedisKaraokeKeyPrefix + fmt.Sprintf("%d", roomid)

	songBytes, err := json.Marshal(song)
	if err != nil {
		log.Error(prefix+"json.Marshal(user)", zap.Any("err", err))
		return
	}

	if err = RedisClient.LRem(redisKey, 0, string(songBytes)).Err(); err != nil {
		log.Error(prefix+"Redis.LRem", zap.Any("err", err))
		return
	}

	log.Info(prefix+"end", zap.Any("roomid", roomid), zap.Any("song", song))
	return
}

func RedisGetSongsList(roomid uint64) (songsList []*bilin.KaraokeSongInfo, err error) {
	const prefix = "RedisGetSongsList "

	redisKey := RedisKaraokeKeyPrefix + fmt.Sprintf("%d", roomid)

	var redisVal []string
	if redisVal, err = RedisClient.LRange(redisKey, 0, -1).Result(); err != nil && err != redis.Nil {
		log.Error(prefix+"redis.LRange", zap.Any("err", err))
		return
	}

	for _, value := range redisVal {
		user := &bilin.KaraokeSongInfo{}
		if err = json.Unmarshal([]byte(value), user); err != nil {
			log.Warn(prefix+"json.Unmarshal", zap.Any("value", value), zap.Any("err", err))
			return nil, err
		}

		songsList = append(songsList, user)
	}

	log.Debug(prefix, zap.Any("roomid", roomid), zap.Any("songsList", songsList))
	return
}

func RedisGetSongInfoByID(roomid uint64, id string) (song *bilin.KaraokeSongInfo, err error) {
	const prefix = "RedisGetSongInfoByID "

	allInfo, _ := RedisGetSongsList(roomid)
	for _, value := range allInfo {
		if value.Id == id {
			song = value
			break
		}
	}

	log.Info(prefix+"end", zap.Any("roomid", roomid), zap.Any("song", song))
	return
}

//获取队列的第一首歌, 不需要出队列
func RedisGetDisplaySong(roomid uint64) (song *bilin.KaraokeSongInfo, err error) {
	const prefix = "RedisGetDisplaySong "
	song = &bilin.KaraokeSongInfo{}

	var redisVal string
	redisKey := RedisKaraokeKeyPrefix + fmt.Sprintf("%d", roomid)
	if redisVal, err = RedisClient.LIndex(redisKey, 0).Result(); err != nil && err != redis.Nil {
		log.Error(prefix+"Redis.LIndex", zap.Any("err", err))
		return nil, err
	}

	if len(redisVal) == 0 {
		log.Info(prefix+"room not find in redis", zap.Any("roomid", roomid))
		return nil, nil
	}

	if err = json.Unmarshal([]byte(redisVal), song); err != nil {
		log.Warn(prefix+"Unmarshal failed", zap.Any("roomid", roomid))
		return nil, err
	}

	log.Info(prefix+"end", zap.Any("roomid", roomid), zap.Any("song", song))
	return
}

//设置displaysong状态，不需要出队列
func RedisChangeDisplaySongStatus(roomid uint64, song *bilin.KaraokeSongInfo) (err error) {
	const prefix = "RedisChangeDisplaySongStatus "

	redisKey := RedisKaraokeKeyPrefix + fmt.Sprintf("%d", roomid)

	songBytes, err := json.Marshal(song)
	if err != nil {
		log.Error(prefix+"json.Marshal(user)", zap.Any("err", err))
		return
	}

	if err = RedisClient.LSet(redisKey, 0, string(songBytes)).Err(); err != nil {
		log.Error(prefix+"redis.LSet", zap.Any("err", err))
		return
	}

	log.Info(prefix+"end", zap.Any("roomid", roomid), zap.Any("song", song))
	return
}

//当前有预告歌曲，插到前面
func RedisAddSongBeforePrepareSong(roomid uint64, prepareSong *bilin.KaraokeSongInfo, song *bilin.KaraokeSongInfo) (err error) {
	const prefix = "RedisAddSongBeforePrepareSong "

	redisKey := RedisKaraokeKeyPrefix + fmt.Sprintf("%d", roomid)

	prepareSongBytes, err := json.Marshal(prepareSong)
	if err != nil {
		log.Error(prefix+"json.Marshal(prepareSong)", zap.Any("err", err), zap.Any("prepareSong", prepareSong))
		return
	}

	songBytes, err := json.Marshal(song)
	if err != nil {
		log.Error(prefix+"json.Marshal(song)", zap.Any("err", err), zap.Any("song", song))
		return
	}
	if err = RedisClient.LInsertBefore(redisKey, string(prepareSongBytes), string(songBytes)).Err(); err != nil {
		log.Error(prefix+"redis.LInsertBefore", zap.Any("err", err))
		return
	}

	log.Info(prefix+"end", zap.Any("roomid", roomid), zap.Any("song", song), zap.Any("prepareSong", prepareSong))
	return
}

//当前有singsing或者pause的歌曲，插到后面
func RedisAddSongAfterSingSong(roomid uint64, singsingSong *bilin.KaraokeSongInfo, song *bilin.KaraokeSongInfo) (err error) {
	const prefix = "RedisAddSongAfterSingSong "

	redisKey := RedisKaraokeKeyPrefix + fmt.Sprintf("%d", roomid)

	singsingSongBytes, err := json.Marshal(singsingSong)
	if err != nil {
		log.Error(prefix+"json.Marshal(singsingSong)", zap.Any("err", err), zap.Any("singsingSong", singsingSong))
		return
	}

	songBytes, err := json.Marshal(song)
	if err != nil {
		log.Error(prefix+"json.Marshal(song)", zap.Any("err", err), zap.Any("song", song))
		return
	}
	if err = RedisClient.LInsertAfter(redisKey, string(singsingSongBytes), string(songBytes)).Err(); err != nil {
		log.Error(prefix+"redis.LInsertBefore", zap.Any("err", err))
		return
	}

	log.Info(prefix+"end", zap.Any("roomid", roomid), zap.Any("song", song), zap.Any("singsingSong", singsingSong))
	return
}
