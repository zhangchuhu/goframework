package rank_test

import (
	"bilin/apigateway/services/rank"
	_ "code.yy.com/yytars/goframework/tars/servant"
	"testing"
	//"time"
	//"fmt"
)

func TestGetChannelBillRankCode(t *testing.T) {
	code := rank.GetChannelBillRankCode()
	t.Logf("GetChannelBillRankCode: %s", code)	
}

func TestGetChannelBillRankList(t *testing.T) {
	err, info := rank.GetChannelBillRankList()
	if err != nil {
		t.Error("GetChannelBillRankList error:" + err.Error())
	} else {
		t.Logf("GetChannelBillRankList success, info: %v", info)
	}
}


