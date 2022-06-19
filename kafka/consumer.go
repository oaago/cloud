package kafka

import (
	"encoding/json"
	"time"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/oaago/component/logx"
)

func NewConsumer(callback ConsumerCallback) {
	consumer := ConsumerOptions
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	config.Consumer.Offsets.CommitInterval = 1 * time.Second
	config.Consumer.Offsets.Initial = sarama.OffsetOldest //初始从最新的offset开始
	var err error
	ConsumerGroup, err = cluster.NewConsumer(consumer.Nodes, consumer.GroupId, consumer.Topic, config)
	defer ConsumerGroup.Close()
	if err != nil {
		logx.Logger.Info(err.Error())
		return
	}
	go func() {
		for err := range ConsumerGroup.Errors() {
			logx.Logger.Error(err.Error())
		}
	}()
	go func() {
		for note := range ConsumerGroup.Notifications() {
			content, _ := json.Marshal(note)
			logx.Logger.Info(string(content))
		}
	}()
	for msg := range ConsumerGroup.Messages() {
		callback(msg, ConsumerGroup)
	}
}
