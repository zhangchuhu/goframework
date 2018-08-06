package service

import (
	"bilin/bcserver/config"
	"testing"
)

var (
	mysqlconfig = &config.AppConfig{
		MysqlAddr: "bilin:ZG7qEsNi2@tcp(183.36.124.123:6304)/bilin_hongbao?charset=utf8",
	}
)

func init() {
	config.SetTestAppConfig(mysqlconfig)
	MysqlInit()
}

func TestMysqlGetVipUser(t *testing.T) {
	ret, err := MysqlGetVipUser(17795535)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(ret)
}
