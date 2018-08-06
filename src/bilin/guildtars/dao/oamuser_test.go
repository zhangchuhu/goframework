/*
 * Copyright (c) 2018-07-26.
 * Author: kordenlu
 * 功能描述:${<VARIABLE_NAME>}
 */

package dao

import "testing"

func TestOAMUser_Create(t *testing.T) {
	oamuser := &OAMUser{
		Username: "lubaoquan",
		Passwd:   "123456",
		Role:     1,
	}
	if err := oamuser.Create(); err != nil {
		t.Error(err)
	}
}
