package nacos

import (
	"encoding/json"
	"github.com/ghodss/yaml"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/oaago/component/logx"
	"github.com/oaago/component/op"
)

func GetNacosConfig(ini *NacosType) *op.Config {
	if NacosClient.error != nil {
		panic(NacosClient.error)
	}
	content, err := NacosClient.ConfigClient.GetConfig(vo.ConfigParam{
		DataId: ini.DataId,
		Group:  ini.Group,
	})

	err = NacosClient.ConfigClient.ListenConfig(vo.ConfigParam{
		DataId: ini.DataId,
		Group:  ini.Group,
		OnChange: func(namespace, group, dataId, data string) {
			logx.Logger.Info("group:" + group + ", dataId:" + dataId + ", data:" + data)
			if err != nil {
				print(err)
			}
			if len(data) == 0 {
				panic("获取配置失败")
			}
			var conf op.Config
			jsonByte, err := yaml.YAMLToJSON([]byte(data))
			err = json.Unmarshal(jsonByte, &conf)
			if err != nil {
				panic(err)
			}
		},
	})
	if err != nil {
		panic(err)
	}
	if len(content) == 0 {
		panic("获取配置失败")
	}
	var conf op.Config
	jsonByte, e := yaml.YAMLToJSON([]byte(content))
	if e != nil {
		panic("转译nacos失败" + e.Error())
	}
	err = json.Unmarshal(jsonByte, &conf)
	if err != nil {
		panic("配置转写结构失败")
	}
	return &conf
}
