package dao

import (
	"bilin/official/service"
	"context"
	"testing"
	"time"
)

func TestContract_Create(t *testing.T) {
	hostuid := int64(33424800)
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
	mycontract.GuildID = int64(turnovercontract.Sid)
	mycontract.HostUid = int64(turnovercontract.LiveUid)
	mycontract.GuildSharePercentage = int64(turnovercontract.Weight)
	mycontract.HostSharePercentage = int64(100 - turnovercontract.Weight)

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
	t.Log(contract_)
}

func TestGetContractsByGuildID(t *testing.T) {
	info, err := GetContractsByGuildID(1)
	if err != nil {
		t.Error(err)
	}
	t.Log(info)
}

func TestGetContracts(t *testing.T) {
	info, err := GetAllContract()
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", info)
}

func TestMigrateContract(t *testing.T) {
	if err := MigrateContract(); err != nil {
		t.Error(err)
	}
}

func TestContract_Update(t *testing.T) {
	c := Contract{GuildSharePercentage: 20, HostSharePercentage: 80,
		ContractEndTime: time.Now().Add(time.Hour * 24 * 60),
		ContractState:   1,
	}
	c.ID = 1
	err := c.Update()
	if err != nil {
		t.Error(err)
	}
}

func TestContract_Delete(t *testing.T) {
	c := Contract{GuildSharePercentage: 20, HostSharePercentage: 80,
		ContractState: 1,
	}
	c.ID = 2
	if err := c.Delete(); err != nil {
		t.Error(err)
	}
}
