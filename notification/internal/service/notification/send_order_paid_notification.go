package notification

import (
	"context"

	"github.com/pinai4/spaceship-factory/notification/internal/model"
)

func (s *service) SendOrderPaidNotification(ctx context.Context, event model.OrderPaidEvent) error {
	// fmt.Printf("Sending order paid notification. Event: %+v\n", event)
	return nil
}
