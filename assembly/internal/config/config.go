package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/pinai4/spaceship-factory/assembly/internal/config/env"
)

type Config struct {
	Logger                     LoggerConfig
	Kafka                      KafkaConfig
	ShipAssembledEventProducer ShipAssembledEventProducerConfig
	OrderPaidEventConsumer     OrderPaidEventConsumerConfig
}

func Load(path ...string) (*Config, error) {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return nil, err
	}

	kafkaCfg, err := env.NewKafkaConfig()
	if err != nil {
		return nil, err
	}

	shipAssembledEventProducerCfg, err := env.NewShipAssembledEventProducerConfig()
	if err != nil {
		return nil, err
	}

	orderPaidEventConsumerCfg, err := env.NewOrderPaidEventConsumerConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		Logger:                     loggerCfg,
		Kafka:                      kafkaCfg,
		ShipAssembledEventProducer: shipAssembledEventProducerCfg,
		OrderPaidEventConsumer:     orderPaidEventConsumerCfg,
	}, nil
}
