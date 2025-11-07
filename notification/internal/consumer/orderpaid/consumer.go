package orderpaid

import (
	"context"
	"fmt"

	def "github.com/pinai4/spaceship-factory/notification/internal/consumer"
	"github.com/pinai4/spaceship-factory/notification/internal/service"
	"github.com/pinai4/spaceship-factory/platform/pkg/kafka"
)

var _ def.Consumer = (*consumer)(nil)

type consumer struct {
	platformConsumer    kafka.Consumer
	notificationService service.NotificationService
}

func New(platformConsumer kafka.Consumer, notificationService service.NotificationService) *consumer {
	return &consumer{
		platformConsumer:    platformConsumer,
		notificationService: notificationService,
	}
}

func (s *consumer) Run(ctx context.Context) error {
	err := s.platformConsumer.Consume(ctx, s.handler)
	if err != nil {
		return fmt.Errorf("failed to consume 'order.paid' topic: %w", err)
	}

	return nil
}
