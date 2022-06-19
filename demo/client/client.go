package main

import (
	"context"
	"github.com/oaago/component/demo/client/pb"
	"github.com/oaago/component/logx"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"time"
)

type RpcClientType struct {
	*clientv3.Client
	EtcdAddr          string
	Url               string
	RemoteServiceName string
}

func NewRpcClient(rpc RpcClientType) *pb.GreetResp {
	cli, err := clientv3.NewFromURL(rpc.EtcdAddr)
	if err != nil {
		panic("无法获取连接 etcd")
	}
	rpc.Client = cli
	builder, err := resolver.NewBuilder(rpc.Client)
	conn, err := grpc.Dial("etcd:///service/go-demo",
		grpc.WithResolvers(builder),
		grpc.WithInsecure())
	if err != nil {
		logx.Logger.Error(err)
	}
	fooClient := pb.NewFooClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	resp, err := fooClient.Greet(ctx, &pb.GreetReq{
		MyName: "Bar",
		Msg:    "Hello, World",
	})
	defer cancel()
	return resp
}

func main() {
	cli, err := clientv3.NewFromURL("http://localhost:2379")
	if err != nil {
		logx.Logger.Error(err)
	}
	builder, err := resolver.NewBuilder(cli)
	if err != nil {
		logx.Logger.Error(err)
	}
	conn, err := grpc.Dial("etcd:///service/go-demo",
		grpc.WithResolvers(builder),
		grpc.WithInsecure())
	if err != nil {
		logx.Logger.Error(err)
	}
	fooClient := pb.NewFooClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	resp, err := fooClient.Greet(ctx, &pb.GreetReq{
		MyName: "Bar",
		Msg:    "Hello, World",
	})
	if err != nil {
		logx.Logger.Error(err)
	}

	logx.Logger.Info(resp.Msg)
}
