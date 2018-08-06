package index

import (
	"strings"
)

type MatchMap map[string]bool

type Index struct {
	Cache map[uint64]MatchMap
}

func (i *Index) InitCache() {
	i.Cache = make(map[uint64]MatchMap)
}
func (i *Index) addItem(id uint64, value string) {
	if m, ok := i.Cache[id]; !ok {
		match_map := make(MatchMap)
		i.Cache[id] = match_map
		match_map[value] = true
	} else {
		m[value] = true
	}
}

func (i *Index) Match(id uint64, matchData string) bool {
	m, ok := i.Cache[id]
	if !ok {
		return false
	}

	if _, ok := m["match_all"]; ok {
		return true
	}

	if _, ok = m[matchData]; ok {
		return true
	}

	return false
}

//initData 英文逗号分割的数据列表,  matchAll 匹配所有的字符模式
func (i *Index) Add(id uint64, initData, matchAll string) {
	data := strings.Replace(initData, " ", "", -1)
	data = strings.Replace(initData, "\n", "", -1)
	data = strings.Replace(initData, "\r", "", -1)
	data = strings.Replace(initData, "\t", "", -1)

	if data == matchAll {
		i.addItem(id, "match_all")
	}

	s := strings.Split(data, ",")
	for _, v := range s {
		i.addItem(id, v)
	}
}
