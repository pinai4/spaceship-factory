package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/pinai4/spaceship-factory/order/internal/model"
)

func (s *service) Get(ctx context.Context, orderUUID uuid.UUID) (model.Order, error) {
	return s.orderRepository.Get(ctx, orderUUID)
}
