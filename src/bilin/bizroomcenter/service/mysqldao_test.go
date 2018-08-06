package service

import (
	"bilin/bizroomcenter/config"
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

func TestMysqlGetRoomInfo(t *testing.T) {
	ret, err := MysqlGetBizRoomInfo(11)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(ret)
}
