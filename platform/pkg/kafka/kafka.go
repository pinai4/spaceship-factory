package kafka

import (
	"context"

	"github.com/pinai4/spaceship-factory/platform/pkg/kafka/consumer"
)

type Consumer interface {
	Consume(ctx context.Context, handler consumer.MessageHandler) error
}
