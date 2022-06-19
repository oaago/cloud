package main

import (
	"context"
	"fmt"
	"github.com/oaago/component/config"
	_ "github.com/oaago/component/config"
	"github.com/oaago/component/demo/cache/dao"
	"github.com/oaago/component/demo/client/controller"
	"github.com/oaago/component/demo/client/pb"
	_ "github.com/oaago/component/etcd"
	"github.com/oaago/component/etcd/rpc"
	"github.com/oaago/component/logx"
	"github.com/oaago/component/mysql"
	_ "github.com/oaago/component/op"
	"github.com/oaago/component/redis"
	"google.golang.org/grpc"
	"net"
	"sync"
	"time"
)

func main() {
	fmt.Println("111", config.Op)
	logx.Logger.Info("111")
	db := mysql.GetDBByName("")
	redis := redis.Client
	fmt.Println(db.Name(), redis)
	//grpcDemo()
	logx.Logger.Info("111")
	logx.Logger.Info("111")
	getFromCache()
}

func grpcDemo() {
	addr := "127.0.0.1:9900"
	ctx := context.Background()
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logx.Logger.Error("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterFooServer(s, controller.NewFooController())
	if err := rpc.Register(ctx, "go-demo", addr); err != nil { // 服务注册名: go-demo
		logx.Logger.Error("register %s failed:%v", "go-demo", err)
	}
	fmt.Printf("start grpc server:%s", addr)
	if err := s.Serve(lis); err != nil {
		logx.Logger.Error("failed to serve: %v", err)
	}
}

func getFromCache() {
	var wait sync.WaitGroup
	wait.Add(20)
	for i := 0; i < 20; i++ {
		go getData(i, &wait)
		time.Sleep(time.Millisecond * time.Duration(500))
	}
	wait.Wait()
}

func getData(i int, wait *sync.WaitGroup) {
	logx.Logger.Info("进入协程:", i)
	model := dao.InitFrontStaffUserListModel(dao.Where{
		IntoTimeStart: "2022-02-01 00:00:00",
		IntoTimeEnd:   "2022-02-08 00:00:00",
		All:           "1",
	})
	model.GetFrontStaffUserExcel()
	//list := model.FrontStaffUserListResult.UserListSearchResultList
	//logx.Logger.Info("查询出来的list:", list)
	wait.Done()
}
