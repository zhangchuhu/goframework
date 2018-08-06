package service

import (
	"bilin/bcserver/config"
	"bilin/bcserver/domain/entity"
	"testing"
)

var (
	mysqlconfig = &config.AppConfig{
		MysqlAddr: "bilin:ZG7qEsNi2@tcp(183.36.124.123:6304)/bilin_hongbao?charset=utf8",
		KafkaAddr: []string{"14.17.103.229:20023", "116.31.112.143:20017", "14.215.104.221:20032"},
	}
)

func init() {
	config.SetTestAppConfig(mysqlconfig)
	MysqlInit()
}

func TestMysqlGetRoomInfo(t *testing.T) {
	ret, err := MysqlGetRoomInfo(11)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(ret)
}

func TestMysqlStorageRoomInfo(t *testing.T) {
	room := entity.NewRoom(10086)
	room.Title = "This is a test title!"
	room.Owner = 141155
	err := MysqlStorageRoomInfo(room)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(err)
}
