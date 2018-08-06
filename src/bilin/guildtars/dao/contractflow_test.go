/*
 * Copyright (c) 2018-07-23.
 * Author: kordenlu
 * 功能描述:${<VARIABLE_NAME>}
 */

package dao

import "testing"

func TestContractFlow_Create(t *testing.T) {
	cf := ContractFlow{
		OperationUid: 100,
		Operation:    1,
		ContractId:   123,
	}
	if err := cf.Create(); err != nil {
		t.Error(err)
	}
}
