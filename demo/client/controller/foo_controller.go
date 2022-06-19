package controller

import (
	"context"
	"fmt"
	"github.com/oaago/component/demo/client/pb"
)

type FooController struct {
}

func NewFooController() *FooController {
	f := &FooController{}
	return f
}

func (f *FooController) Greet(ctx context.Context, in *pb.GreetReq) (*pb.GreetResp, error) {
	reply := fmt.Sprintf("Hello %s, I got your msg:%s", in.GetMyName(), in.GetMsg())
	out := &pb.GreetResp{}
	out.Msg = reply
	return out, nil
}
