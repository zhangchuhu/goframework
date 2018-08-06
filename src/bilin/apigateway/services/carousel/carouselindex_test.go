package carousel_test

import (
	"bilin/apigateway/services/carousel"
	_ "code.yy.com/yytars/goframework/tars/servant"
	"testing"
	//"time"
)

func TestMatchUserCarousel(t *testing.T) {
	info := carousel.MatchUserCarousel("","", 0)
	t.Logf("MatchUserCarousel success, info: %v", info)
	
}

func TestGetAllCarousel(t *testing.T) {
	info := carousel.GetAllCarousel()
	t.Logf("GetAllCarousel success, info: %v", info)
	
}
