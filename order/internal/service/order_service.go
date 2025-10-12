package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/pinai4/spaceship-factory/order/internal/model"
)

type OrderService interface {
	Create(ctx context.Context, orderUUID uuid.UUID, createOrder model.CreateOrder) (float64, error)
	Get(ctx context.Context, orderUUID uuid.UUID) (model.Order, error)
	ProcessPayment(ctx context.Context, orderUUID uuid.UUID, paymentMethod model.OrderPaymentMethod) (uuid.UUID, error)
	Cancel(ctx context.Context, orderUUID uuid.UUID) error
}
