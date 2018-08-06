package rank_test

import (
	"bilin/apigateway/config"
	"bilin/apigateway/services/rank"
	_ "code.yy.com/yytars/goframework/tars/servant"
	"testing"
	//"time"
)

func TestReloadRankInfo(t *testing.T) {
	if err := config.InitAndSubConfig("appconfig.json"); err != nil {
		t.Error("InitAndSubConfig failed")
		t.Logf("InitAndSubConfig failed\n\n\n\n")
		return
	}
	rank.InitThriftConnentPool("221.228.110.138:6903;")
	rank.LoadContributeRankInfo()
	rank.LoadAnchorRankInfo()
	rank.LoadGuardRankInfo()
	rank.LoadChannelBillRankInfo()

	info := rank.GetContributeRank()
	if info == nil {
		t.Error("GetContributeRank error")
	} else {
		t.Logf("GetContributeRank success, info: %v", *info)
	}

	info = rank.GetAnchorRank()
	if info == nil {
		t.Error("GetAnchorRank error")
	} else {
		t.Logf("GetAnchorRank success, info: %v", *info)
	}

	info = rank.GetGuardRank()
	if info == nil {
		t.Error("GetGuardRank error")
	} else {
		t.Logf("GetGuardRank success, info: %v", *info)
	}

	channel_bill_info := rank.GetChannelBillRank()
	if info == nil {
		t.Error("GetChannelBillRank error")
	} else {
		t.Logf("GetChannelBillRank success, info: %v", *channel_bill_info)
	}
}
