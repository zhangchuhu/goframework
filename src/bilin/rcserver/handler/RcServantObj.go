package handler

import (
	"context"
	"bilin/protocol"
)

type RcServantObj struct {
}

func NewRcServantObj() *RcServantObj {
	return &RcServantObj{}
}

func (this *RcServantObj) StartRandomCall(ctx context.Context, req *bilin.StartRandomCallReq) (*bilin.StartRandomCallResp, error) {
	return &bilin.StartRandomCallResp{
		Greet: ", Good day!",
	}, nil
}
