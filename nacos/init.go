package nacos

import (
	"encoding/json"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/oaago/component/config"
)

type NacosType struct {
	ServiceName    string `json:"service_name"`
	IpAddr         string `json:"ip_addr"`
	NamespaceId    string `json:"namespace_id"`
	LogDir         string `json:"log_dir"`
	CacheDir       string `json:"cache_dir"`
	DataId         string `json:"data_id"`
	Group          string `json:"group"`
	Discover       bool
	Env            string
	Sc             []constant.ServerConfig
	Cc             constant.ClientConfig
	ConfigClient   config_client.IConfigClient
	DiscoverClient naming_client.INamingClient
	error          error
	CallBack       func(content string)
	Cluster        string
}

var NacosClient = &NacosType{}

var pd = map[string]string{
	"d":       "241e09bd-76e1-46c2-92db-b8dd1154a369",
	"t":       "665ca2f1-bd50-45bf-8e34-40ad4bd975e6",
	"o":       "a4a0021c-e696-4f0a-b375-4805a2f7f644",
	"p":       "16561b3c-d80a-444b-9249-9f26cb8f2f88",
	"test":    "665ca2f1-bd50-45bf-8e34-40ad4bd975e6",
	"offline": "a4a0021c-e696-4f0a-b375-4805a2f7f644",
	"prod":    "16561b3c-d80a-444b-9249-9f26cb8f2f88",
}

var ServerEnv string

func init() {
	nacosStr := config.Op.GetString("kafka.nacos")
	NacosClient.ServiceName = config.Op.GetString("server.name")
	NacosClient.Env = config.Op.GetString("server.env")
	NacosClient.DataId = NacosClient.ServiceName
	NacosClient.Group = NacosClient.ServiceName
	NacosClient.IpAddr = "ops.laodianhuang.cn"
	json.Unmarshal([]byte(nacosStr), NacosClient)
	NacosClient.Cluster = NacosClient.Env + "-" + NacosClient.ServiceName
}
func NewNacos() *NacosType {
	NacosClient.Sc = []constant.ServerConfig{
		{
			IpAddr: NacosClient.IpAddr,
			Scheme: "http",
			Port:   80,
		},
	}
	NacosClient.Cc = constant.ClientConfig{
		NamespaceId:         pd[ServerEnv], //namespace id
		LogDir:              "./nacos/logs/",
		CacheDir:            "./nacos/cache",
		TimeoutMs:           2000,
		NotLoadCacheAtStart: true,
		LogLevel:            "debug",
	}
	NacosClient.ConfigClient, NacosClient.error = clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &NacosClient.Cc,
			ServerConfigs: NacosClient.Sc,
		})
	NacosClient.DiscoverClient, _ = clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &NacosClient.Cc,
			ServerConfigs: NacosClient.Sc,
		},
	)
	return NacosClient
}
