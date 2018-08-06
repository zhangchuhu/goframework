package rank_test

import (
	"bilin/apigateway/services/rank"
	_ "code.yy.com/yytars/goframework/tars/servant"
	"testing"
	//"time"
)

func TestGetAnchorRankList(t *testing.T) {
	rank.InitThriftConnentPool("221.228.110.138:6903;")
	err, info := rank.GetAnchorRankList()
	if err != nil {
		t.Error("GetAnchorRankList error:" + err.Error())
	} else {
		t.Logf("GetAnchorRankList success, info: %v", info)
	}
}

func TestGetAnchorRankInfo(t *testing.T) {
	rank.InitThriftConnentPool("221.228.110.138:6903;")
	info, err := rank.GetAnchorRankInfo()
	if err != nil {
		t.Error("GetAnchorRankInfo error:" + err.Error())
	} else {
		t.Logf("GetAnchorRankInfo success, info: %v", info)
	}
}
