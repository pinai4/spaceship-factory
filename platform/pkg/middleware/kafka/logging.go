package kafka

import (
	"context"

	"github.com/pinai4/spaceship-factory/platform/pkg/kafka/consumer"
	"github.com/pinai4/spaceship-factory/platform/pkg/logger"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...logger.Field)
}

func Logging(log Logger) consumer.Middleware {
	return func(next consumer.MessageHandler) consumer.MessageHandler {
		return func(ctx context.Context, msg consumer.Message) error {
			log.Info(ctx, "Message received", logger.String("topic", msg.Topic), logger.Int64("offset", msg.Offset))
			res := next(ctx, msg)
			log.Info(ctx, "Message handled", logger.String("topic", msg.Topic), logger.Int64("offset", msg.Offset))
			return res
		}
	}
}
