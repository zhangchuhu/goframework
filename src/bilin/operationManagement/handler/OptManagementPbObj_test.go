package handler

import (
	"bilin/bcserver/config"
	"bilin/operationManagement/service"
	"bilin/protocol"
	"context"
	"testing"
)

var (
	spb       *OptManagementPbObj
	appconfig = &config.AppConfig{
		RedisAddr: "183.36.122.50:4019",
		MysqlAddr: "bilin:ZG7qEsNi2@tcp(183.36.124.123:6304)/bilin_hongbao?charset=utf8",
	}
)

func init() {
	config.SetTestAppConfig(appconfig)
	spb = NewOptManagementPbObj()

	service.MysqlInit()
}

func TestActDistributionHeadgear(t *testing.T) {
	req := &bilin.ActDistributionHeadgearRequest{Hinfo: &bilin.HeadgearInfo{
		Uid:        9999999,
		Headgear:   "https://bilinoperationmanagement.bs2dl.yy.com/2fec7fa157e60a8b33c636f9522083fe.jpg",
		Effecttime: 1528989011,
		Expiretime: 1528999011,
	}}
	resp, _ := spb.ActDistributionHeadgear(context.TODO(), req)

	t.Logf("done", resp)
}

func TestGetUserHeadgearInfo(t *testing.T) {
	req := &bilin.GetUserHeadgearInfoReq{
		Uid: 9999999}
	resp, _ := spb.GetUserHeadgearInfo(context.TODO(), req)

	t.Logf("done", resp)
}
