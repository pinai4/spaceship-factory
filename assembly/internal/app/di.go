package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"

	"github.com/pinai4/spaceship-factory/assembly/internal/config"
	"github.com/pinai4/spaceship-factory/assembly/internal/service"
	assemblyService "github.com/pinai4/spaceship-factory/assembly/internal/service/assembly"
	"github.com/pinai4/spaceship-factory/assembly/internal/service/assembly/producer"
	"github.com/pinai4/spaceship-factory/platform/pkg/closer"
)

type diContainer struct {
	config *config.Config
	closer *closer.Closer

	assemblyService  service.AssemblyService
	assemblyProducer service.AssemblyProducer

	syncProducer sarama.SyncProducer
}

func NewDiContainer(config *config.Config, closer *closer.Closer) *diContainer {
	return &diContainer{config: config, closer: closer}
}

func (d *diContainer) AssemblyService() service.AssemblyService {
	if d.assemblyService == nil {
		d.assemblyService = assemblyService.NewService(d.AssemblyProducer())
	}

	return d.assemblyService
}

func (d *diContainer) AssemblyProducer() service.AssemblyProducer {
	if d.assemblyProducer == nil {
		topics := map[string]string{
			producer.ShipAssembledEventTopicKey: d.Config().ShipAssembledEventProducer.Topic(),
		}
		d.assemblyProducer = producer.NewProducer(d.SyncProducer(), topics)
	}

	return d.assemblyProducer
}

func (d *diContainer) SyncProducer() sarama.SyncProducer {
	if d.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			d.Config().Kafka.Brokers(),
			d.Config().ShipAssembledEventProducer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create sync producer: %s\n", err.Error()))
		}
		d.closer.AddNamed("Sarama (kafka) sync producer", func(ctx context.Context) error {
			return p.Close()
		})

		d.syncProducer = p
	}

	return d.syncProducer
}

func (d *diContainer) Config() *config.Config {
	return d.config
}
