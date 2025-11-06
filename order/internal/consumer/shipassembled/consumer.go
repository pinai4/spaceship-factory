package shipassembled

import (
	"context"
	"fmt"

	def "github.com/pinai4/spaceship-factory/order/internal/consumer"
	"github.com/pinai4/spaceship-factory/order/internal/service"
	"github.com/pinai4/spaceship-factory/platform/pkg/kafka"
)

var _ def.Consumer = (*consumer)(nil)

type consumer struct {
	platformConsumer kafka.Consumer
	orderService     service.OrderService
}

func New(platformConsumer kafka.Consumer, orderService service.OrderService) *consumer {
	return &consumer{
		platformConsumer: platformConsumer,
		orderService:     orderService,
	}
}

func (s *consumer) Run(ctx context.Context) error {
	err := s.platformConsumer.Consume(ctx, s.handler)
	if err != nil {
		return fmt.Errorf("failed to consume 'ship.assembled' topic: %w", err)
	}

	return nil
}
