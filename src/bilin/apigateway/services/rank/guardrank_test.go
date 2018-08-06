package rank_test

import (
	"bilin/apigateway/services/rank"
	_ "code.yy.com/yytars/goframework/tars/servant"
	"testing"
	//"time"
)

func TestGetGuardRankList(t *testing.T) {
	rank.InitThriftConnentPool("221.228.110.138:6903;")
	err, info := rank.GetGuardRankList()
	if err != nil {
		t.Error("GetGuardRankList error:" + err.Error())
	} else {
		t.Logf("GetGuardRankList success, info: %v", info)
	}
}


func TestGetGuardRankInfo(t *testing.T) {
	rank.InitThriftConnentPool("221.228.110.138:6903;")
	info, err := rank.GetGuardRankInfo()
	if err != nil {
		t.Error("GetGuardRankInfo error:" + err.Error())
	} else {
		t.Logf("GetGuardRankInfo success, info: %v", info)
	}
}

