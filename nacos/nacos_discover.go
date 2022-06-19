package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func (c NacosType) RegisterInstance() {
	c.DiscoverClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          NacosClient.IpAddr,
		Port:        8848,
		ServiceName: NacosClient.ServiceName,
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"idc": "shanghai"},
		ClusterName: c.ServiceName, // 默认值DEFAULT
		GroupName:   c.ServiceName, // 默认值DEFAULT_GROUP
	})
}

func (c NacosType) DeregisterInstance() (bool, error) {
	success, err := c.DiscoverClient.DeregisterInstance(vo.DeregisterInstanceParam{
		Ip:          c.IpAddr,
		Port:        8848,
		ServiceName: c.ServiceName,
		Ephemeral:   true,
		Cluster:     c.ServiceName, // 默认值DEFAULT
		GroupName:   c.ServiceName, // 默认值DEFAULT_GROUP
	})
	return success, err
}

func (c NacosType) SelectOneHealthyInstance() (*model.Instance, error) {
	instance, err := c.DiscoverClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: c.ServiceName,
		GroupName:   c.ServiceName,           // 默认值DEFAULT_GROUP
		Clusters:    []string{c.ServiceName}, // 默认值DEFAULT
	})
	return instance, err
}
