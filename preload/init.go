package preload

import (
	_ "github.com/oaago/component/config"
	_ "github.com/oaago/component/mysql"
	"github.com/oaago/component/op"
	_ "github.com/oaago/component/redis"
)

func LoadConfig() *op.Config {
	return op.ConfigData
}
