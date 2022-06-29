package etcd

import (
	"encoding/json"
	"github.com/oaago/cloud/config"
	"github.com/oaago/cloud/logx"
	"go.etcd.io/etcd/client/v3"
	"time"
)

type EtcdType struct {
	Endpoints []string
	Username  string
	Password  string
	Client    *clientv3.Client
}

var EtcdOptions *EtcdType

func init() {
	etcdStr := config.Op.GetString("etcd")
	if len(etcdStr) == 0 {
		return
	}
	json.Unmarshal([]byte(etcdStr), &EtcdOptions)
	var err error
	EtcdOptions.Client, err = clientv3.New(clientv3.Config{
		Endpoints:   EtcdOptions.Endpoints,
		DialTimeout: 5 * time.Second,
		Password:    EtcdOptions.Password,
		Username:    EtcdOptions.Username,
	})
	if err != nil {
		// handle error!
		logx.Logger.Error("connect to etcd failed, err:%v\n", err)
	}
	logx.Logger.Info("etcd 连接成功")
}
