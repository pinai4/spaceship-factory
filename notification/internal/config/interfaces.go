package config

import "github.com/IBM/sarama"

type LoggerConfig interface {
	Level() string
	AsJSON() bool
}

type KafkaConfig interface {
	Brokers() []string
}

type OrderPaidEventConsumerConfig interface {
	Topic() string
	GroupID() string
	Config() *sarama.Config
}

type ShipAssembledEventConsumerConfig interface {
	Topic() string
	GroupID() string
	Config() *sarama.Config
}

type TelegramBotConfig interface {
	Token() string
	ChatID() int64
}
