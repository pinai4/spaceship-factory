package config

import "github.com/IBM/sarama"

type LoggerConfig interface {
	Level() string
	AsJSON() bool
}

type KafkaConfig interface {
	Brokers() []string
}

type ShipAssembledEventProducerConfig interface {
	Topic() string
	Config() *sarama.Config
}

type OrderPaidEventConsumerConfig interface {
	Topic() string
	GroupID() string
	Config() *sarama.Config
}
