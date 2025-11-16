package config

import (
	"time"

	"github.com/IBM/sarama"
)

type LoggerConfig interface {
	Level() string
	AsJSON() bool
}

type HTTPServerConfig interface {
	Address() string
	ReadTimeout() time.Duration
}

type PostgresConfig interface {
	DSN() string
}

type DBMigrationsConfig interface {
	Path() string
}

type InventoryGRPCClientConfig interface {
	Address() string
}

type PaymentGRPCClientConfig interface {
	Address() string
}

type IAMGRPCClientConfig interface {
	Address() string
}

type KafkaConfig interface {
	Brokers() []string
}

type OrderPaidEventProducerConfig interface {
	Topic() string
	Config() *sarama.Config
}

type ShipAssembledEventConsumerConfig interface {
	Topic() string
	GroupID() string
	Config() *sarama.Config
}
