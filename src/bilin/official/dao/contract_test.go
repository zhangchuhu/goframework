package dao

import (
	"bilin/official/service"
	"context"
	"testing"
	"time"
)

func TestContract_Create(t *testing.T) {
	hostuid := uint64(33424800)
	mycontract := Contract{
		GuildID: 39818246,
		HostUid: hostuid,
		//ContractStartTime:    time.Now(),
		//ContractEndTime:      time.Now().Add(60 * time.Minute * 24 * 365),
		GuildSharePercentage: 10,
		HostSharePercentage:  90,
	}
	//mycontract.ID = 4

	err := putContractToTurnOver(&mycontract)
	if err != nil {
		t.Error(err)
		return
	}

	turnovercontract, err := service.QueryContractByAnchor(context.TODO(), int64(hostuid))
	if err != nil {
		t.Error(err)
	}
	mycontract.ContractEndTime = time.Unix(turnovercontract.FinishTime/1000, 0)
	mycontract.ContractStartTime = time.Unix(turnovercontract.SignTime/1000, 0)
	mycontract.GuildID = uint64(turnovercontract.Sid)
	mycontract.HostUid = uint64(turnovercontract.LiveUid)
	mycontract.GuildSharePercentage = uint64(turnovercontract.Weight)
	mycontract.HostSharePercentage = uint64(100 - turnovercontract.Weight)

	if err := mycontract.Create(); err != nil {
		t.Error(err)
	}
	t.Log(mycontract)
}

func TestContract_Get(t *testing.T) {
	var (
		contract_ *Contract
		err       error
	)
	if contract_, err = GetContractByHostUid(17796525); err != nil {
		t.Error(err)
	}
	if contract_ == nil {
		t.Log("no contract")
	}
	t.Log(contract_)
}

func TestGetContractsByGuildID(t *testing.T) {
	info, err := GetContractsByGuildID(1)
	if err != nil {
		t.Error(err)
	}
	t.Log(info)
}
