package cache

import (
	"sort"
	"testing"
)

func TestSortForPointer(t *testing.T) {
	var recs []*RecommandLivingInfo
	for i := 0; i < 10; i++ {
		recs = append(recs, &RecommandLivingInfo{
			UserCount:  uint64(i),
			SortWeight: int64(120 - i),
		})
	}

	sort.Sort(OnlineUserSlice(recs))
	t.Logf("recs:%v", recs)
	sort.Sort(RecommandLivingInfoSlice(recs))
	t.Log(recs)
}
