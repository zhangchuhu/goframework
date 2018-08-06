package dao

import (
	"bilin/official/config"
	"bilin/official/service"
	"fmt"
	"os"
	"testing"
)

const (
	testcontractipport    = "58.215.52.27:6907"
	productcontractipport = "101.226.22.226:6907"

	testturnover    = "58.215.52.27:6903"
	productturnover = "101.226.22.226:6903"
)

func TestMain(m *testing.M) {
	fmt.Println("enter TestMain")
	var err error

	conf := &config.AppConfig{
		ContractThriftAddr: testcontractipport,
		TurnOverThriftAddr: testturnover,
	}
	err = service.InitTurnOverService(conf)
	if err != nil {
		fmt.Println("InitTurnOverService err", err)
	}
	os.Exit(m.Run())
}

func TestGuild_Create(t *testing.T) {
	//newguild := Guild{
	//	OW:        17795053,
	//	Title:     "海内",
	//	Mobile:    "13832452345",
	//	Describle: "海内工会",
	//	GuildLogo: "http://tianya.jpg",
	//}
	//newguild.ID = 17795053

	newguild := Guild{
		OW:        39818246,
		Title:     "测试2",
		Mobile:    "13832452345",
		Describle: "测试2工会",
		GuildLogo: "http://tianya.jpg",
	}
	newguild.ID = uint(newguild.OW)

	if err := newguild.Create(); err != nil {
		t.Error(err)
	}
}

func TestGuild_Get(t *testing.T) {
	if myguild, err := GetByGuildID(17795287); err != nil {
		t.Error(err)
	} else {
		t.Log(myguild)
	}

}

func TestGuild_UpdateByOw(t *testing.T) {
	myguild := &Guild{
		Title:     "第吧工会111",
		GuildLogo: "xxx",
		Mobile:    "12311",
		Describle: "比邻第一工会111",
	}
	if err := myguild.UpdateByOw(41127363); err != nil {
		t.Error(err)
	}
	t.Log(myguild)
}
