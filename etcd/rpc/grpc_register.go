package rpc

import (
	"context"
	"fmt"
	"github.com/oaago/cloud/config"
	"github.com/oaago/cloud/logx"
	"time"

	uuid "github.com/satori/go.uuid"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
)

//var client *clientv3.Client

const (
	prefix = "service"
)

type EtcdType struct {
	client   *clientv3.Client
	Username string
	Password string
}

var Etcd = &EtcdType{}

func init() {
	endpoints := config.Op.GetStringSlice("etcd.endpoints")
	enable := config.Op.GetBool("etcd.enable")
	if len(endpoints) > 0 && enable {
		client, err := clientv3.New(clientv3.Config{
			Endpoints:   endpoints,
			DialTimeout: 3 * time.Second,
			//Username:    config.Op.GetString("etcd.username"),
			//Password:    config.Op.GetString("etcd.password"),
		})
		if err != nil {
			panic(err)
		}
		Etcd.client = client
	}
}

func Register(ctx context.Context, serviceName, addr string) error {
	logx.Logger.Info("Try register to etcd ...")
	// 创建一个租约
	lease := clientv3.NewLease(Etcd.client)
	cancelCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	leaseResp, err := lease.Grant(cancelCtx, 3)
	if err != nil {
		return err
	}

	leaseChannel, err := lease.KeepAlive(ctx, leaseResp.ID) // 长链接, 不用设置超时时间
	if err != nil {
		return err
	}

	em, err := endpoints.NewManager(Etcd.client, prefix)
	if err != nil {
		return err
	}

	cancelCtx, cancel = context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	if err := em.AddEndpoint(cancelCtx, fmt.Sprintf("%s/%s/%s", prefix, serviceName, uuid.NewV4().String()), endpoints.Endpoint{
		Addr: addr,
	}, clientv3.WithLease(leaseResp.ID)); err != nil {
		return err
	}
	logx.Logger.Info("Register etcd success")

	del := func() {
		logx.Logger.Info("Register close")

		cancelCtx, cancel = context.WithTimeout(ctx, time.Second*3)
		defer cancel()
		em.DeleteEndpoint(cancelCtx, serviceName)

		lease.Close()
	}
	// 保持注册状态(连接断开重连)
	keepRegister(ctx, leaseChannel, del, serviceName, addr)

	return nil
}

func keepRegister(ctx context.Context, leaseChannel <-chan *clientv3.LeaseKeepAliveResponse, cleanFunc func(), serviceName, addr string) {
	go func() {
		failedCount := 0
		for {
			select {
			case resp := <-leaseChannel:
				if resp != nil {
					//log.Println("keep alive success.")
				} else {
					logx.Logger.Error("keep alive failed.")
					failedCount++
					for failedCount > 3 {
						cleanFunc()
						if err := Register(ctx, serviceName, addr); err != nil {
							time.Sleep(time.Second)
							continue
						}
						return
					}
					continue
				}
			case <-ctx.Done():
				cleanFunc()
				Etcd.client.Close()
				return
			}
		}
	}()
}
