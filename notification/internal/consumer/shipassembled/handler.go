package shipassembled

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
		return fmt.Errorf("failed to decode 'ship.assembled' topic message: %w", err)
	}

	if err := s.notificationService.SendShipAssembledNotification(ctx, event); err != nil {
		return fmt.Errorf("failed to send ship assembled notification: %w", err)
	}

	return nil
}

func (s *consumer) decode(data []byte) (model.ShipAssembledEvent, error) {
	var pb eventsV1.ShipAssembled
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.ShipAssembledEvent{}, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	return model.ShipAssembledEvent{
		EventUUID:    pb.GetEventUuid(),
		OrderUUID:    pb.GetOrderUuid(),
		UserUUID:     pb.GetUserUuid(),
		BuildTimeSec: pb.GetBuildTimeSec(),
	}, nil
}
