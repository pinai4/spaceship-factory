package producer

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"google.golang.org/protobuf/proto"

	"github.com/pinai4/spaceship-factory/assembly/internal/model"
	def "github.com/pinai4/spaceship-factory/assembly/internal/service"
	eventsV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/events/v1"
)

var _ def.AssemblyProducer = (*producer)(nil)

const ShipAssembledEventTopicKey = "ship-assembled-event"

type producer struct {
	sarmaProducer sarama.SyncProducer
	topics        map[string]string
}

func NewProducer(sarmaProducer sarama.SyncProducer, topics map[string]string) *producer {
	return &producer{
		sarmaProducer: sarmaProducer,
		topics:        topics,
	}
}

func (p *producer) ProduceShipAssembled(_ context.Context, event model.ShipAssembledEvent) error {
	topic, ok := p.topics[ShipAssembledEventTopicKey]
	if !ok {
		return fmt.Errorf("producer.ProduceShipAssembled failed to get topic for ShipAssembledEvent")
	}

	msg := &eventsV1.ShipAssembled{
		EventUuid:    event.EventUUID,
		OrderUuid:    event.OrderUUID,
		UserUuid:     event.UserUUID,
		BuildTimeSec: event.BuildTimeSec,
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		return fmt.Errorf("producer.ProduceShipAssembled failed to marshal ShipAssembledEvent: %w", err)
	}

	if _, _, err := p.sarmaProducer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(event.EventUUID),
		Value: sarama.ByteEncoder(payload),
	}); err != nil {
		return fmt.Errorf("producer.ProduceShipAssembled failed to send message: %w", err)
	}

	return nil
}
