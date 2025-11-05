package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"

	"github.com/pinai4/spaceship-factory/assembly/internal/config"
	"github.com/pinai4/spaceship-factory/assembly/internal/consumer"
	"github.com/pinai4/spaceship-factory/assembly/internal/consumer/orderpaid"
	"github.com/pinai4/spaceship-factory/platform/pkg/closer"
	platformConsumer "github.com/pinai4/spaceship-factory/platform/pkg/kafka/consumer"
	"github.com/pinai4/spaceship-factory/platform/pkg/logger"
	platformMiddleware "github.com/pinai4/spaceship-factory/platform/pkg/middleware/kafka"
)

type App struct {
	config      *config.Config
	closer      *closer.Closer
	logger      logger.Logger
	diContainer *diContainer

	saramaConsumerGroup sarama.ConsumerGroup
	consumer            consumer.Consumer
}

func New(ctx context.Context, config *config.Config, closer *closer.Closer, logger logger.Logger) (*App, error) {
	a := &App{config: config, closer: closer, logger: logger}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	return a.runConsumer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initSarmaConsumerGroup,
		a.initConsumer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDiContainer(a.config, a.closer)
	return nil
}

func (a *App) initSarmaConsumerGroup(_ context.Context) error {
	var err error
	a.saramaConsumerGroup, err = sarama.NewConsumerGroup(
		a.config.Kafka.Brokers(),
		a.config.OrderPaidEventConsumer.GroupID(),
		a.config.OrderPaidEventConsumer.Config(),
	)
	if err != nil {
		return fmt.Errorf("failed to create consumer group: %w", err)
	}
	a.closer.AddNamed("Sarama (kafka) consumer group", func(ctx context.Context) error {
		return a.saramaConsumerGroup.Close()
	})

	return nil
}

func (a *App) initConsumer(_ context.Context) error {
	pfmConsumer := platformConsumer.NewConsumer(
		a.saramaConsumerGroup,
		[]string{
			a.config.OrderPaidEventConsumer.Topic(),
		},
		a.logger,
		platformMiddleware.Logging(a.logger),
	)

	a.consumer = orderpaid.New(pfmConsumer, a.diContainer.AssemblyService())

	return nil
}

func (a *App) runConsumer(ctx context.Context) error {
	a.logger.Info(ctx, "🚀 Consumer running", logger.String("topic", a.config.OrderPaidEventConsumer.Topic()))

	if err := a.consumer.Run(ctx); err != nil {
		return err
	}

	return nil
}
