package service

import (
	"bilin/guildtars/config"
	"context"
	"fmt"
	"os"
	"testing"
	"time"
)

const (
	testcontractipport    = "58.215.52.27:6907"
	productcontractipport = "101.226.22.226:6907"

	testturnover    = "58.215.52.27:6903"
	productturnover = "101.226.22.226:6903"
)

func TestMain(m *testing.M) {
	conf := &config.AppConfig{
		ContractThriftAddr: testcontractipport,
	}
	err := InitTurnOverService(conf)
	if err != nil {
		fmt.Println("InitTurnOverService err", err)
	}
	os.Exit(m.Run())
}

func TestAddContractInfoExternal(t *testing.T) {
	ctx, cancle := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancle()
	ret, err := AddContractInfoExternal(ctx, 40859046, 33424800, 33424800, 10)
	if err != nil {
		t.Error(err)
	}
	t.Log(ret)
}

func TestQueryContractByAnchor(t *testing.T) {
	ctx, cancle := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancle()

	info, err := QueryContractByAnchor(ctx, 34991830)
	if err != nil {
		t.Error(err)
	}
	t.Log(info)
}
