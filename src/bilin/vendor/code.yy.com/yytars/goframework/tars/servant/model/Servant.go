package model

import (
	"code.yy.com/yytars/goframework/jce/taf"
	"code.yy.com/yytars/goframework/tars/servant/protocol"
	"context"
)

type Servant interface {
	Taf_invoke(
		ctx context.Context,
		ctype byte,
		sFuncName string,
		buf []byte) (*taf.ResponsePacket, error)
}

type PbServant interface {
	Pb_invoke(
		ctx context.Context,
		ctype byte,
		sFuncName string,
		buf []byte,
		status map[string]string,
		context map[string]string) (*pbtaf.ResponsePacket, error)
}
