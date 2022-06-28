package kafka

import (
	"encoding/json"
	"github.com/Shopify/sarama"
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

var ConsumerOptions = &Consumer{}
var ProducerOptions = &Producer{}

var ProducerList = &ProducerType{}

type ConsumerCallback func(*sarama.ConsumerMessage, *sarama.Consumer)

func init() {
	consumerEnable := config.Op.GetBool("kafka.consumer.enable")
	producerEnable := config.Op.GetBool("kafka.producer.enable")
	if consumerEnable {
		consumerStr := config.Op.GetStringMapStringSlice("kafka.consumer")
		marshal, err := json.Marshal(consumerStr)
		if err != nil {
			return
		}
		json.Unmarshal(marshal, &ConsumerOptions)
	}
	if producerEnable {
		producerStr := config.Op.GetStringMap("kafka.producer")
		marshal, _ := json.Marshal(producerStr)
		json.Unmarshal(marshal, ProducerOptions)
	}
}
