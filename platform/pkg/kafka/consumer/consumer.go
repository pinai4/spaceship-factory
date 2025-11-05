package consumer

import (
	"context"
	"errors"

	"github.com/IBM/sarama"

	"github.com/pinai4/spaceship-factory/platform/pkg/logger"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...logger.Field)
	Error(ctx context.Context, msg string, fields ...logger.Field)
}

type consumer struct {
	group       sarama.ConsumerGroup
	topics      []string
	logger      Logger
	middlewares []Middleware
}

// NewConsumer — creates a new consumer.
func NewConsumer(group sarama.ConsumerGroup, topics []string, logger Logger, middlewares ...Middleware) *consumer {
	return &consumer{
		group:       group,
		topics:      topics,
		logger:      logger,
		middlewares: middlewares,
	}
}

// Consume starts the consumer for the list of topics.
func (c *consumer) Consume(ctx context.Context, handler MessageHandler) error {
	newGroupHandler := NewGroupHandler(handler, c.logger, c.middlewares...)

	for {
		if err := c.group.Consume(ctx, c.topics, newGroupHandler); err != nil {
			if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				return nil
			}

			c.logger.Error(ctx, "Consume error", logger.Error(err))
			return err
		}

		if ctx.Err() != nil {
			return nil
		}

		c.logger.Info(ctx, "Consumer group rebalancing...")
	}
}
