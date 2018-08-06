package service

import (
	"bilin/thrift/gen-go/turnover"
	"context"
	"testing"
	"time"
)

func TestQueryRevenueRecord(t *testing.T) {
	info, count, err := QueryRevenueRecord(context.TODO(), 17795069, 0, "20180601000000", "20180613235959",
		1, 1, 10, 0)
	if err != nil {
		t.Error(err)
	}
	t.Log(count, len(info), info)
}

func TestUserAccountByUidAndType(t *testing.T) {
	info, err := UserAccountByUidAndType(context.TODO(), 100, turnover.TAppId_Bilin, turnover.TCurrencyType_Bilin_Profit)
	if err != nil {
		t.Error(err)
	}
	t.Log(info)
}

func TestQueryCurMonthRevenueRecordPaging(t *testing.T) {
	startTime, err := time.Parse(TimeLayoutOthers, "2018-06-01 00:00:00")
	if err != nil {
		t.Error(err)
	}

	endTime, err := time.Parse(TimeLayoutOthers, "2018-06-10 00:00:00")
	if err != nil {
		t.Error(err)
	}

	info, err := QueryRevenueRecordPaging(context.TODO(),
		17795052,
		10,
		1,
		2,
		startTime.Unix()*1000,
		endTime.Unix()*1000,
		0,
		0,
	)
	if err != nil {
		t.Error(err)
	}
	t.Log(info)
}

func TestQueryCurMonthRevenueRecord(t *testing.T) {
	info, err := QueryCurMonthRevenueRecord(context.TODO(), 0, 18, 1, 1, 10, 17796525)
	if err != nil {
		t.Error(err)
	}
	t.Log(info)
}
