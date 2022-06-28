package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/Shopify/sarama"
)

type msgConsumerGroup struct{}

func (msgConsumerGroup) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (msgConsumerGroup) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h msgConsumerGroup) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		fmt.Printf("%s Message topic:%q partition:%d offset:%d  value:%s\n", msg.Topic, msg.Partition, msg.Offset, string(msg.Value))
		sess.MarkMessage(msg, "")
	}
	return nil
}

var handler msgConsumerGroup

func NewConsumer(callback ConsumerCallback) {
	consume := ConsumerOptions
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.CommitInterval = 1 * time.Second
	config.Consumer.Offsets.Initial = sarama.OffsetOldest //初始从最新的offset开始
	var err error
	consumer, err := sarama.NewConsumerGroup(consume.Nodes, consume.GroupId, config)
	if err != nil {
		fmt.Printf("consumer_test create consumer error %s\n", err.Error())
		return
	}
	//defer consumer.Close()
	for {
		err := consumer.Consume(context.Background(), []string{"testAutoSyncOffset"}, handler)
		if err != nil {
			fmt.Errorf(err.Error())
		}
	}
}
