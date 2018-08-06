package handler

import (
	"context"
	"%{AppName}/protocol"
)

type %{ServantName}Obj struct {
}

func New%{ServantName}Obj() *%{ServantName}Obj {
	return &%{ServantName}Obj{}
}

func (p *%{ServantName}Obj) SayGreeting(ctx context.Context, r *%{AppName}.Request) (*%{AppName}.Response, error) {
	return &%{AppName}.Response{
		Greet: r.Name + ", Good day!",
	}, nil
}
