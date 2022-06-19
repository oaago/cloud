package kafka

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/oaago/component/config"
)

type KafkaType struct {
	Consumer
	Producer
}

type ProducerType struct {
	sarama.SyncProducer
	sarama.AsyncProducer
	Mode string
}

type Producer struct {
	Nodes []string `yaml:"nodes"`
	Topic string   `yaml:"topic"`
}

var Producers *ProducerType

type Consumer struct {
	Nodes   []string `yaml:"nodes"`
	Topic   []string `yaml:"topic"`
	GroupId string   `yaml:"groupId"`
}

var ConsumerGroup *cluster.Consumer
var ConsumerOptions = &Consumer{}
var ProducerOptions = &Producer{}

var ProducerList *ProducerType

type ConsumerCallback func(*sarama.ConsumerMessage, *cluster.Consumer)

func init() {
	consumerStr := config.Op.GetString("kafka.consumer")
	producerStr := config.Op.GetString("kafka.producer")
	if len(consumerStr) > 0 {
		json.Unmarshal([]byte(consumerStr), ConsumerOptions)
	}
	if len(producerStr) > 0 {
		json.Unmarshal([]byte(producerStr), ProducerOptions)
	}
}
