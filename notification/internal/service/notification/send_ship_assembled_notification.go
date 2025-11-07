package notification

import (
	"context"

	"github.com/pinai4/spaceship-factory/notification/internal/model"
)

func (s *service) SendShipAssembledNotification(ctx context.Context, event model.ShipAssembledEvent) error {
	// fmt.Printf("Sending ship assembled notification. Event: %+v\n", event)
	return nil
}
