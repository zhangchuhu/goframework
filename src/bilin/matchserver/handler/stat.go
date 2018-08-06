package handler

import (
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
)

const (
	CallRecordDir   = "/data/bilin/matchserver/callrecord"
	CallOpStart     = 'o'
	CallOpEnd       = 'x'
	CallTypeUnknown = '-'
	CallTypeDirect  = 'D'
	CallTypeHetero  = 'O'
	CallTypeHomo    = 'S'

	OnlineCountDir = "/data/bilin/matchserver/onlinecount"

	PlayRecordDir = "/data/bilin/matchserver/playrecord"
	SexMale       = 'M'
	SexFemale     = 'F'

	SpamRecordDir = "/data/bilin/matchserver/spamrecord"
)

type Player struct {
	UserID      int64     // 用户ID
	UserSex     int       // 性别
	PlayType    int       // 玩法类型 @see CallType
	BeginTime   time.Time // 开始匹配
	Giveup1Time time.Time // 没匹配上：主动取消
	SelectTime  time.Time // 等待选择或被选择
	Giveup2Time time.Time // 没被选中或匹配中取消
	SuccessTime time.Time // 进入通话
}

var (
	playersLock sync.Mutex
	players     map[int64]*Player
)

func init() {
	if err := os.MkdirAll(CallRecordDir, 0755); err != nil {
		log.Error("init: MkdirAll", zap.Error(err))
	}
	if err := os.MkdirAll(OnlineCountDir, 0755); err != nil {
		log.Error("init: MkdirAll", zap.Error(err))
	}
	if err := os.MkdirAll(PlayRecordDir, 0755); err != nil {
		log.Error("init: MkdirAll", zap.Error(err))
	}
	if err := os.MkdirAll(SpamRecordDir, 0755); err != nil {
		log.Error("init: MkdirAll", zap.Error(err))
	}
	players = make(map[int64]*Player)
}

func WriteFile(filename string, data string) {
	const prefix = "WriteFile "
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Error(prefix+"OpenFile", zap.Error(err), zap.String("filename", filename), zap.String("data", data))
		return
	}
	_, err = f.Write([]byte(data))
	if err1 := f.Close(); err == nil {
		err = err1
	}
	if err != nil {
		log.Error(prefix+"Write", zap.Error(err), zap.String("filename", filename), zap.String("data", data))
	}
	return
}

func WriteCallRecord(op int, uid1, uid2 int64, rid int64, calltype int) {
	now := time.Now().Local()
	nowDate := now.Format("2006-01-02")
	nowStamp := now.Format("2006-01-02T15:04:05")
	fileName := path.Join(CallRecordDir, nowDate+".txt")
	WriteFile(fileName, fmt.Sprintf("%c %d %d %d %s %c\n", op, uid1, uid2, rid, nowStamp, calltype))
}

func WriteOnlineCount(onlineAll, malePlaying, femalePlaying, maleWaitingO, femaleWaitingO, maleWaitingS, femaleWaitingS int64) {
	now := time.Now().Local()
	nowDate := now.Format("2006-01-02")
	nowStamp := now.Format("2006-01-02T15:04:05")
	fileName := path.Join(OnlineCountDir, nowDate+".txt")
	WriteFile(fileName, fmt.Sprintf("%s %d %d %d %d %d %d %d\n", nowStamp,
		onlineAll, malePlaying, femalePlaying, maleWaitingO, femaleWaitingO, maleWaitingS, femaleWaitingS))
}

func formatTime(t time.Time) (s string) {
	if t.IsZero() {
		s = "-"
	} else {
		s = t.Format("2006-01-02T15:04:05")
	}
	return
}

func writePlayRecord(p *Player) {
	now := time.Now().Local()
	nowDate := now.Format("2006-01-02")
	fileName := path.Join(PlayRecordDir, nowDate+".txt")
	WriteFile(fileName, fmt.Sprintf("%d %c %c %s %s %s %s %s\n", p.UserID, p.UserSex, p.PlayType,
		formatTime(p.BeginTime), formatTime(p.Giveup1Time), formatTime(p.SelectTime),
		formatTime(p.Giveup2Time), formatTime(p.SuccessTime)))
}

func PlayerBegin(uid int64, sex int, playtype int) {
	playersLock.Lock()
	defer playersLock.Unlock()

	p, ok := players[uid]
	if ok {
		writePlayRecord(p)
	}
	players[uid] = &Player{
		UserID:    uid,
		UserSex:   sex,
		PlayType:  playtype,
		BeginTime: time.Now(),
	}
}

func PlayerGiveup1(uid int64) {
	playersLock.Lock()
	defer playersLock.Unlock()

	if p, ok := players[uid]; ok {
		p.Giveup1Time = time.Now()
		writePlayRecord(p)
		delete(players, uid)
	}
}

func PlayerSelect(uid int64) {
	playersLock.Lock()
	defer playersLock.Unlock()

	if p, ok := players[uid]; ok {
		p.SelectTime = time.Now()
	}
}

func PlayerGiveup2(uid int64) {
	playersLock.Lock()
	defer playersLock.Unlock()

	if p, ok := players[uid]; ok {
		p.Giveup2Time = time.Now()
		writePlayRecord(p)
		delete(players, uid)
	}
}

func PlayerSuccess(uid int64) {
	playersLock.Lock()
	defer playersLock.Unlock()

	if p, ok := players[uid]; ok {
		p.SuccessTime = time.Now()
		writePlayRecord(p)
		delete(players, uid)
	}
}

func WriteSpamRecord(uid int64, sex int, playtype int, level int32, cheat bool) {
	now := time.Now().Local()
	nowDate := now.Format("2006-01-02")
	nowStamp := now.Format("2006-01-02T15:04:05")
	fileName := path.Join(SpamRecordDir, nowDate+".txt")
	cheatInt := 0
	if cheat {
		cheatInt = 1
	}
	WriteFile(fileName, fmt.Sprintf("%s %d %c %c %d %d\n", nowStamp, uid, sex, playtype, level, cheatInt))
}
