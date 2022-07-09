package preload

import (
	_ "github.com/oaago/cloud/config"
	_ "github.com/oaago/cloud/etcd/rpc"
	_ "github.com/oaago/cloud/kafka"
	_ "github.com/oaago/cloud/mysql"
	_ "github.com/oaago/cloud/nacos"
	"github.com/oaago/cloud/op"
	_ "github.com/oaago/cloud/redis"
)

func LoadConfig() *op.Config {
	return op.ConfigData
}
