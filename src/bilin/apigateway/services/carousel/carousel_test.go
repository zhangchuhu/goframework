package carousel_test

import (
	"bilin/apigateway/services/carousel"
	_ "code.yy.com/yytars/goframework/tars/servant"
	"testing"
	//"time"
)

func TestGetEffectiveCarousel(t *testing.T) {
	info, err := carousel.GetEffectiveCarousel()
	if err != nil {
		t.Error("GetEffectiveCarousel error:" + err.Error())
	} else {
		t.Logf("GetEffectiveCarousel success, info: %v", info)
	}
}
