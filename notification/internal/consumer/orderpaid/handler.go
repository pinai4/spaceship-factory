package orderpaid

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/pinai4/spaceship-factory/notification/internal/model"
	kafkaConsumer "github.com/pinai4/spaceship-factory/platform/pkg/kafka/consumer"
	eventsV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/events/v1"
)

func (s *consumer) handler(ctx context.Context, msg kafkaConsumer.Message) error {
	event, err := s.decode(msg.Value)
	if err != nil {
		return fmt.Errorf("failed to decode 'order.paid' topic message: %w", err)
	}

	if err := s.notificationService.SendOrderPaidNotification(ctx, event); err != nil {
		return fmt.Errorf("failed to send order paid notification: %w", err)
	}

	return nil
}

func (s *consumer) decode(data []byte) (model.OrderPaidEvent, error) {
	var pb eventsV1.OrderPaid
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.OrderPaidEvent{}, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	return model.OrderPaidEvent{
		EventUUID:       pb.GetEventUuid(),
		OrderUUID:       pb.GetOrderUuid(),
		UserUUID:        pb.GetUserUuid(),
		PaymentMethod:   pb.GetPaymentMethod(),
		TransactionUUID: pb.GetTransactionUuid(),
	}, nil
}
