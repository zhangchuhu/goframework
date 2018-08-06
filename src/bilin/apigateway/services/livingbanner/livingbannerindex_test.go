package livingbanner_test

import (
	"bilin/apigateway/services/livingbanner"
	_ "code.yy.com/yytars/goframework/tars/servant"
	"testing"
)

func TestMatchUserCarousel(t *testing.T) {
	info := livingbanner.MatchLivingBanner("", "", 0, 2)
	t.Logf("MatchLivingBanner success, info: %v", info)

}
