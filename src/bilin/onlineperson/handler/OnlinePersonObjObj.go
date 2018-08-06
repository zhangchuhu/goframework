package handler

import (
	"context"
	"bilin/protocol"
)

type OnlinePersonObjObj struct {
}

func NewOnlinePersonObjObj() *OnlinePersonObjObj {
	return &OnlinePersonObjObj{}
}

func (p *OnlinePersonObjObj) SayGreeting(ctx context.Context, r *bilin.Request) (*bilin.Response, error) {
	return &bilin.Response{
		Greet: r.Name + ", Good day!",
	}, nil
}
