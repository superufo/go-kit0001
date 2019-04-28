package impl

import (
	"context"
	_ "strings"
)

type HelloServiceImpl struct{}

func (p *HelloServiceImpl) Hello(ctx context.Context, args *String,) (*String, error) {
	reply := &String{Value: "hello:" + args.GetValue()}
	return reply, nil
}