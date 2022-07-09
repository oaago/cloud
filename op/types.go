package op

import (
	"fmt"
	"github.com/oaago/cloud/etcd/rpc"

	"github.com/oaago/cloud/config"
	"github.com/oaago/cloud/elastic"
	"github.com/oaago/cloud/kafka"
	"github.com/oaago/cloud/logx"
	"github.com/oaago/cloud/oss"
	"github.com/oaago/cloud/redis"
)

type Nacos struct {
	IpAddr      string `json:"ip_addr"`
	NamespaceId string `json:"namespace_id"`
	LogDir      string `json:"log_dir"`
	CacheDir    string `json:"cache_dir"`
	DataId      string `json:"data_id"`
	Group       string `json:"group"`
}

type Server struct {
	Name     string `yaml:"name"`
	Port     int    `yaml:"port"`
	Env      string `yaml:"env"`
	Version  string `json:"version"` //服务版本
	Weight   int64  `json:"weight"`  //服务权重
	BasePath string `yaml:"basePath"`
}

type Docker struct {
	Harbor struct {
		Url string `json:"url"`
	}
}

type Config struct {
	Server  `json:"server"`
	Nacos   `json:"nacos"`
	Kafka   kafka.KafkaType     `json:"kafka"`
	Mysql   map[string]string   `json:"mysql,omitempty"`
	Redis   redis.RedisType     `json:"redis"`
	Logger  logx.LoggerType     `json:"logger"`
	Elastic elastic.ElasticType `json:"elastic"`
	OSS     oss.AliyunType      `json:"oss"`
	Etcd    rpc.EtcdType        `json:"etcd"`
	UCSDK   struct {
		AuthServerURL string `yaml:"authServerURL" json:"auth_server_url,omitempty"`
		ClientID      string `yaml:"clientID" json:"client_id,omitempty"`
		ClientSecret  string `yaml:"clientSecret" json:"client_secret,omitempty"`
		RedirectURL   string `yaml:"redirectURL" json:"redirect_url,omitempty"`
	} `json:"ucsdk"`
	CodeMap  map[int]string         `json:"code_map,omitempty"`
	SelfData map[string]interface{} `json:"self_data,omitempty"`
}

var ConfigData *Config

func init() {
	err := config.Op.Unmarshal(&ConfigData)
	if err != nil {
		fmt.Errorf(err.Error())
		return
	}
	fmt.Println(ConfigData)
}
