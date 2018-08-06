package handler

import (
	"strconv"
	"strings"
	"testing"
	"time"

	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
)

func TestSortMapByValue(t *testing.T) {
	m := map[string]string{
		"1": `{"uid": 1, "timestamp": 100}`,
		"2": `{"uid": 2, "timestamp": 98}`,
		"3": `{"uid": 3, "timestamp": 101}`,
	}
	v := sortMapByValue(m)
	t.Log(v)
}

func doMatch(fromList UserList, toList UserList, matchCount int) {
	isPro := false
	startTime := time.Now().UnixNano() / 1e6

	toListNo := 0
	for fromIndex, fromValue := range fromList {
		var matchList UserList = nil
		counter := 0
		for i := toListNo; i < len(toList); i++ {
			matchList = append(matchList, toList[i])
			toListNo++
			counter++
			if counter == matchCount {
				log.Info("DoMatchSex ok",
					zap.Any("fromIndex", fromIndex),
					zap.Any("toListNo", toListNo),
					zap.Any("counter", counter),
					zap.Any("fromValue", fromValue),
					zap.Any("matchList", matchList),
					zap.Any("fromList length", len(fromList)),
					zap.Any("toList length", len(toList)),
					zap.Any("matchCount", matchCount),
					zap.Any("province", isPro))
				break
			}
		}
		if counter >= 1 && counter < matchCount {
			log.Info("DoMatchSex ok not enough male",
				zap.Any("fromIndex", fromIndex),
				zap.Any("toListNo", toListNo),
				zap.Any("counter", counter),
				zap.Any("fromValue", fromValue),
				zap.Any("matchList", matchList),
				zap.Any("fromList length", len(fromList)),
				zap.Any("toList length", len(toList)),
				zap.Any("matchCount", matchCount),
				zap.Any("province", isPro))
			break
		}
	}

	endTime := time.Now().UnixNano() / 1e6
	if duration := endTime - startTime; duration >= 900 {
		log.Warn("DoMatchSex cost time too long",
			zap.Any("fromList length", len(fromList)),
			zap.Any("toList length", len(toList)),
			zap.Any("matchCount", matchCount),
			zap.Any("province", isPro),
			zap.Any("duration", duration))
	}
}

func TestDoMatch(t *testing.T) {
	from := UserList{
		{uint32(1), 0, 0, 0, "haiwai", 1500000000000},
		{uint32(2), 0, 0, 0, "haiwai", 1500000000000},
	}
	to := UserList{
		{uint32(5), 0, 0, 0, "haiwai", 1500000000000},
		{uint32(6), 0, 0, 0, "haiwai", 1500000000000},
		{uint32(7), 0, 0, 0, "haiwai", 1500000000000},
		{uint32(8), 0, 0, 0, "haiwai", 1500000000000},
		{uint32(9), 0, 0, 0, "haiwai", 1500000000000},
		{uint32(10), 0, 0, 0, "haiwai", 1500000000000},
		{uint32(11), 0, 0, 0, "haiwai", 1500000000000},
	}
	doMatch(from, to, 3)
}

func TestParseUidAndCid(t *testing.T) {
	var cid, peeruid uint32
	items := strings.Split("12345,67890", ",")
	if len(items) > 1 {
		if temp, err := strconv.ParseUint(items[1], 10, 32); err == nil {
			cid = uint32(temp)
		}
	}
	if len(items) > 0 {
		if temp, err := strconv.ParseUint(items[0], 10, 32); err == nil {
			peeruid = uint32(temp)
		}
	}
	t.Logf("cid %d, peeruid %d\n", cid, peeruid)
}
