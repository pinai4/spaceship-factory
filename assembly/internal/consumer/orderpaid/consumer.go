package orderpaid

import (
	"context"
	"fmt"

	def "github.com/pinai4/spaceship-factory/assembly/internal/consumer"
	"github.com/pinai4/spaceship-factory/assembly/internal/service"
	"github.com/pinai4/spaceship-factory/platform/pkg/kafka"
)

var _ def.Consumer = (*consumer)(nil)

type consumer struct {
	platformConsumer kafka.Consumer
	assemblyService  service.AssemblyService
}

func New(platformConsumer kafka.Consumer, assemblyService service.AssemblyService) *consumer {
	return &consumer{
		platformConsumer: platformConsumer,
		assemblyService:  assemblyService,
	}
}

func (s *consumer) Run(ctx context.Context) error {
	err := s.platformConsumer.Consume(ctx, s.handler)
	if err != nil {
		return fmt.Errorf("failed to consume 'order.paid' topic: %w", err)
	}

	return nil
}
