package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/pinai4/spaceship-factory/order/internal/model"
)

type OrderRepository interface {
	Create(ctx context.Context, order model.Order) error
	Get(ctx context.Context, orderUUID uuid.UUID) (model.Order, error)
	Update(ctx context.Context, orderUUID uuid.UUID, order model.Order) error
}
