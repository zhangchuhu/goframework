package rank_test

import (
	"bilin/apigateway/services/rank"
	_ "code.yy.com/yytars/goframework/tars/servant"
	"testing"
	//"time"
)

func TestGetContributeRank(t *testing.T) {
	rank.InitThriftConnentPool("221.228.110.138:6903;")
	err, info := rank.GetContributeRankList()
	if err != nil {
		t.Error("GetContributeRankList error:" + err.Error())
	} else {
		t.Logf("GetContributeRankList success, info: %v", info)
	}
}

func TestGetContributeRankInfo(t *testing.T) {
	rank.InitThriftConnentPool("221.228.110.138:6903;")
	info, err := rank.GetContributeRankInfo()
	if err != nil {
		t.Error("GetContributeRankInfo error:" + err.Error())
	} else {
		t.Logf("GetContributeRankInfo success, info: %v", info)
	}
}
