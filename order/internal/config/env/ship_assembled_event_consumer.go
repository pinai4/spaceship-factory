package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type shipAssembledEventConsumerEnvConfig struct {
	Topic   string `env:"SHIP_ASSEMBLED_TOPIC_NAME,required"`
	GroupID string `env:"SHIP_ASSEMBLED_CONSUMER_GROUP_ID,required"`
}

type shipAssembledEventConsumerConfig struct {
	raw shipAssembledEventConsumerEnvConfig
}

func NewShipAssembledEventConsumerConfig() (*shipAssembledEventConsumerConfig, error) {
	var raw shipAssembledEventConsumerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &shipAssembledEventConsumerConfig{raw: raw}, nil
}

func (cfg *shipAssembledEventConsumerConfig) Topic() string {
	return cfg.raw.Topic
}

func (cfg *shipAssembledEventConsumerConfig) GroupID() string {
	return cfg.raw.GroupID
}

func (cfg *shipAssembledEventConsumerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config
}
