package dao

import (
	"bilin/official/config"
	"bilin/official/service"
	"fmt"
	"github.com/jinzhu/gorm"
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
	hujiaoCallDB, err = gorm.Open("mysql", "blguildadmin:1oCrqMZYwr@tcp(14.116.173.148:6305)/Hujiao?readTimeout=500ms&charset=utf8mb4&parseTime=True&loc=Local")
	//hujiaoCallDB, err = gorm.Open("mysql", "bilin_admin:avYsLkYwQ@tcp(58.215.143.9:6307)/Hujiao?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println("open faile", err)
		return
	}

	BIDataDB, err = gorm.Open("mysql", "bilin_bi_data@bilin_bi_data:kT2pa8Yq@tcp(221.228.79.244:8066)/bilin_bi_data?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println("open faile", err)
		return
	}
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

func TestGetByOW(t *testing.T) {
	if info, err := GetGuildByOW(41127363); err != nil {
		t.Error(err)
	} else {
		t.Log(info)
	}
}
