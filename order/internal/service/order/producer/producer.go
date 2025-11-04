package producer

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"google.golang.org/protobuf/proto"

	"github.com/pinai4/spaceship-factory/order/internal/model"
	def "github.com/pinai4/spaceship-factory/order/internal/service"
	eventsV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/events/v1"
)

var _ def.OrderProducer = (*producer)(nil)

const OrderPaidEventTopicKey = "order-paid-event"

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

func (p *producer) ProduceOrderPaid(ctx context.Context, event model.OrderPaidEvent) error {
	topic, ok := p.topics[OrderPaidEventTopicKey]
	if !ok {
		return fmt.Errorf("producer.ProduceOrderPaid failed to get topic for order paid event")
	}

	msg := &eventsV1.OrderPaid{
		EventUuid:       event.EventUUID,
		OrderUuid:       event.OrderUUID,
		UserUuid:        event.UserUUID,
		PaymentMethod:   event.PaymentMethod,
		TransactionUuid: event.TransactionUUID,
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		return fmt.Errorf("producer.ProduceOrderPaid failed to marshal OrderPaidEvent: %w", err)
	}

	if _, _, err := p.sarmaProducer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(event.EventUUID),
		Value: sarama.ByteEncoder(payload),
	}); err != nil {
		return fmt.Errorf("producer.ProduceOrderPaid failed to send message: %w", err)
	}

	return nil
}
