package kafka

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/oaago/component/logx"
	"github.com/tidwall/gjson"
)

func NewProducer(mode string) *ProducerType {
	//初始化配置
	p := ProducerList
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // 发送完数据需要leader和follow都确认
	//config.Producer.Partitioner = sarama.NewRandomPartitioner // 随机分配分区 partition
	config.Producer.Partitioner = sarama.NewManualPartitioner // 人工指定分区
	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回
	// 创建生产者，连接kafka
	var connectErr error
	p.Mode = mode
	if mode == "async" {
		// 异步
		p.AsyncProducer, connectErr = sarama.NewAsyncProducer(ProducerOptions.Nodes, config)
		if connectErr != nil {
			logx.Logger.Error(connectErr.Error(), "AsyncProducer")
			return Producers
		}
	} else {
		// 同步
		p.SyncProducer, connectErr = sarama.NewSyncProducer(ProducerOptions.Nodes, config)
		if connectErr != nil {
			logx.Logger.Error(connectErr.Error(), "SyncProducer")
			return Producers
		}
	}
	logx.Logger.Info("kafka producer 初始化成功", config)
	return p
}

func (p *ProducerType) SyncSendMessage(content string, topic string) {
	msg := &sarama.ProducerMessage{}
	msg.Topic = ProducerOptions.Topic
	if topic != "" {
		msg.Topic = topic
	}
	msg.Value = sarama.StringEncoder(content)
	partition := gjson.Get("content", "config.partition").Int()
	msg.Partition = int32(partition)
	// 发送消息
	pid, offset, err := p.SyncProducer.SendMessage(msg)
	if err != nil {
		logx.Logger.Info(err.Error())
		return
	}
	fmt.Printf("pid:%v offset:%v\n", pid, offset) //分区id，偏移id
}

func (p *ProducerType) Close() {
	if p.Mode == "async" {
		err := p.AsyncProducer.Close()
		if err != nil {
			panic("AsyncProducer关闭失败")
		}
	} else {
		err := p.SyncProducer.Close()
		if err != nil {
			panic("SyncProducer关闭失败")
		}
	}
}
