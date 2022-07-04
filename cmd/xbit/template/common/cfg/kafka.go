package cfg

import (
	"github.com/segmentio/kafka-go"
	"github.com/xbitgo/components/MQ"
	"time"
)

type KafkaConsumer struct {
	Brokers []string `json:"brokers" yaml:"brokers"`
	Topic   string   `json:"topic" yaml:"topic"`
	GroupID string   `json:"group_id" yaml:"group_id"`
}

type KafkaProducer struct {
	Brokers  []string `json:"brokers" yaml:"brokers"`
	Topic    string   `json:"topic" yaml:"topic"`
	Balancer string   `json:"balancer" yaml:"balancer"`
	Async    bool     `json:"async" yaml:"async"`
}

func (k *KafkaConsumer) CreateInstance() (MQ.Consumer, error) {
	r := MQ.NewKafkaConsumer(kafka.ReaderConfig{
		Brokers:  k.Brokers,
		GroupID:  k.GroupID,
		Topic:    k.Topic,
		MinBytes: 10e3, // 10k
		MaxWait:  time.Second,
		MaxBytes: 10e6, // 10MB
	})
	return r, nil
}

func (k *KafkaProducer) CreateInstance() (MQ.Producer, error) {
	w := MQ.NewKafkaProducer(MQ.KafkaProducerConfig{
		Brokers:  k.Brokers,
		Topic:    k.Topic,
		Balancer: k.Balancer,
		Async:    k.Async,
	})
	return w, nil
}
