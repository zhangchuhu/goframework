/*
 * Copyright (c) 2018-07-23.
 * Author: kordenlu
 * 功能描述:签约流水
 */

package handler

type ContractFlow struct {
	OperationUid int64 `json:"operation_uid"`
	ContractId   int64 `json:"contract_id"`
	Operation    int32 `json:"operation"`
}
