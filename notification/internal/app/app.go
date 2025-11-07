package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"

	"github.com/pinai4/spaceship-factory/notification/internal/config"
	"github.com/pinai4/spaceship-factory/notification/internal/consumer"
	"github.com/pinai4/spaceship-factory/notification/internal/consumer/orderpaid"
	"github.com/pinai4/spaceship-factory/notification/internal/consumer/shipassembled"
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

	orderPaidSaramaConsumerGroup sarama.ConsumerGroup
	orderPaidConsumer            consumer.Consumer

	shipAssembledSaramaConsumerGroup sarama.ConsumerGroup
	shipAssembledConsumer            consumer.Consumer
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
	// Channel for receiving errors from components
	errCh := make(chan error, 2)

	// Context to stop all goroutines
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// orderPaidConsumer
	go func() {
		if err := a.runOrderPaidConsumer(ctx); err != nil {
			errCh <- fmt.Errorf("order paid consumer crashed: %w", err)
		}
	}()

	// shipAssembledConsumer
	go func() {
		if err := a.runShipAssembledConsumer(ctx); err != nil {
			errCh <- fmt.Errorf("ship assembled consumer crashed: %w", err)
		}
	}()

	// Wait either for an error or for context completion (e.g., SIGINT/SIGTERM)
	select {
	case <-ctx.Done():
		a.logger.Info(ctx, "Shutdown signal received")
	case err := <-errCh:
		a.logger.Error(ctx, "Component crashed, shutting down", logger.Error(err))
		// Trigger cancel to stop the other component
		cancel()
		// Wait for all tasks to finish (if graceful shutdown is implemented internally)
		<-ctx.Done()
		return err
	}

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initOrderPaidSaramaConsumerGroup,
		a.initOrderPaidConsumer,
		a.initShipAssembledSaramaConsumerGroup,
		a.initShipAssembledConsumer,
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

func (a *App) initOrderPaidSaramaConsumerGroup(_ context.Context) error {
	var err error
	a.orderPaidSaramaConsumerGroup, err = sarama.NewConsumerGroup(
		a.config.Kafka.Brokers(),
		a.config.OrderPaidEventConsumer.GroupID(),
		a.config.OrderPaidEventConsumer.Config(),
	)
	if err != nil {
		return fmt.Errorf("failed to create consumer (order paid) group: %w", err)
	}
	a.closer.AddNamed("Sarama (kafka) consumer (order paid) group", func(ctx context.Context) error {
		return a.orderPaidSaramaConsumerGroup.Close()
	})

	return nil
}

func (a *App) initOrderPaidConsumer(_ context.Context) error {
	pfmConsumer := platformConsumer.NewConsumer(
		a.orderPaidSaramaConsumerGroup,
		[]string{
			a.config.OrderPaidEventConsumer.Topic(),
		},
		a.logger,
		platformMiddleware.Logging(a.logger),
	)

	a.orderPaidConsumer = orderpaid.New(pfmConsumer, a.diContainer.NotificationService())

	return nil
}

func (a *App) runOrderPaidConsumer(ctx context.Context) error {
	a.logger.Info(ctx, "🚀 Consumer running", logger.String("topic", a.config.OrderPaidEventConsumer.Topic()))

	if err := a.orderPaidConsumer.Run(ctx); err != nil {
		return err
	}

	return nil
}

func (a *App) initShipAssembledSaramaConsumerGroup(_ context.Context) error {
	var err error
	a.shipAssembledSaramaConsumerGroup, err = sarama.NewConsumerGroup(
		a.config.Kafka.Brokers(),
		a.config.ShipAssembledEventConsumer.GroupID(),
		a.config.ShipAssembledEventConsumer.Config(),
	)
	if err != nil {
		return fmt.Errorf("failed to create consumer (ship assembled) group: %w", err)
	}
	a.closer.AddNamed("Sarama (kafka) consumer (ship assembled) group", func(ctx context.Context) error {
		return a.shipAssembledSaramaConsumerGroup.Close()
	})

	return nil
}

func (a *App) initShipAssembledConsumer(_ context.Context) error {
	pfmConsumer := platformConsumer.NewConsumer(
		a.shipAssembledSaramaConsumerGroup,
		[]string{
			a.config.ShipAssembledEventConsumer.Topic(),
		},
		a.logger,
		platformMiddleware.Logging(a.logger),
	)

	a.shipAssembledConsumer = shipassembled.New(pfmConsumer, a.diContainer.NotificationService())

	return nil
}

func (a *App) runShipAssembledConsumer(ctx context.Context) error {
	a.logger.Info(ctx, "🚀 Consumer running", logger.String("topic", a.config.ShipAssembledEventConsumer.Topic()))

	if err := a.shipAssembledConsumer.Run(ctx); err != nil {
		return err
	}

	return nil
}
