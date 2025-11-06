package shipassembled

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"

	"github.com/pinai4/spaceship-factory/order/internal/model"
	kafkaConsumer "github.com/pinai4/spaceship-factory/platform/pkg/kafka/consumer"
	eventsV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/events/v1"
)

func (s *consumer) handler(ctx context.Context, msg kafkaConsumer.Message) error {
	event, err := s.decode(msg.Value)
	if err != nil {
		return fmt.Errorf("failed to decode 'ship.assembled' topic message: %w", err)
	}

	if err := s.orderService.Assemble(ctx, event.OrderUUID); err != nil {
		return fmt.Errorf("failed to assemble order: %w", err)
	}

	return nil
}

func (s *consumer) decode(data []byte) (model.ShipAssembledEvent, error) {
	var pb eventsV1.ShipAssembled
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.ShipAssembledEvent{}, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	orderUUID, err := uuid.Parse(pb.GetOrderUuid())
	if err != nil {
		return model.ShipAssembledEvent{}, fmt.Errorf("failed to parse order uuid: %w", err)
	}

	userUUID, err := uuid.Parse(pb.GetUserUuid())
	if err != nil {
		return model.ShipAssembledEvent{}, fmt.Errorf("failed to parse user uuid: %w", err)
	}

	return model.ShipAssembledEvent{
		EventUUID:    pb.GetEventUuid(),
		OrderUUID:    orderUUID,
		UserUUID:     userUUID,
		BuildTimeSec: pb.BuildTimeSec,
	}, nil
}
