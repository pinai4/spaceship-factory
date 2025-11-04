package service

import (
	"context"

	"github.com/pinai4/spaceship-factory/order/internal/model"
)

type OrderProducer interface {
	ProduceOrderPaid(ctx context.Context, event model.OrderPaidEvent) error
}
