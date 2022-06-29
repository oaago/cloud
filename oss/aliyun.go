package oss

import (
	"encoding/json"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/oaago/cloud/config"
	"github.com/oaago/cloud/logx"
)

type AliyunType struct {
	Endpoint        string `yaml:"endpoint"`
	AccessKeyId     string `yaml:"accessKeyId"`
	AccessKeySecret string `yaml:"accessKeySecret"`
	BucketName      string `yaml:"bucketName"`
}

var AliyunOptions = &AliyunType{}

func init() {
	ossStr := config.Op.GetString("oss")
	json.Unmarshal([]byte(ossStr), AliyunOptions)
}

func NewOssBucket() *oss.Bucket {
	client, err := oss.New(AliyunOptions.Endpoint, AliyunOptions.AccessKeyId, AliyunOptions.AccessKeySecret)
	if err != nil {
		logx.Logger.Error("阿里云连接失败")
	}
	bucket, err := client.Bucket(AliyunOptions.BucketName)
	if err != nil {
		logx.Logger.Error("获取bucket失败检查是否一致")
	}
	return bucket
}
