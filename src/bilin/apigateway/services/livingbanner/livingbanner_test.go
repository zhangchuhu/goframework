package livingbanner_test

import (
	_ "code.yy.com/yytars/goframework/tars/servant"
	"testing"
	//"time"
	"bilin/apigateway/services/livingbanner"
)

func TestGetEffectiveCarousel(t *testing.T) {
	info, err := livingbanner.GetEffectiveLivingBanner()
	if err != nil {
		t.Error("GetEffectiveLivingBanner error:" + err.Error())
	} else {
		t.Logf("GetEffectiveLivingBanner success, info: %v", info)
	}
}
