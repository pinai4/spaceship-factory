package service

import (
	"context"

	"github.com/pinai4/spaceship-factory/notification/internal/model"
)

type NotificationService interface {
	SendOrderPaidNotification(ctx context.Context, event model.OrderPaidEvent) error
	SendShipAssembledNotification(ctx context.Context, event model.ShipAssembledEvent) error
}
