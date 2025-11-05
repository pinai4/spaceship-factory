package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type shipAssembledEventProducerEnvConfig struct {
	TopicName string `env:"SHIP_ASSEMBLED_TOPIC_NAME,required"`
}

type shipAssembledEventProducerConfig struct {
	raw shipAssembledEventProducerEnvConfig
}

func NewShipAssembledEventProducerConfig() (*shipAssembledEventProducerConfig, error) {
	var raw shipAssembledEventProducerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &shipAssembledEventProducerConfig{raw: raw}, nil
}

func (cfg *shipAssembledEventProducerConfig) Topic() string {
	return cfg.raw.TopicName
}

// Config return configuration for sarama producer
func (cfg *shipAssembledEventProducerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Producer.Return.Successes = true

	return config
}
