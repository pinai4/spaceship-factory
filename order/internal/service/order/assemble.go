package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/pinai4/spaceship-factory/order/internal/model"
)

func (s *service) Assemble(ctx context.Context, orderUUID uuid.UUID) error {
	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		return fmt.Errorf("OrderService.Assemble get order error: %w", err)
	}

	order.Status = model.OrderStatusAssembled

	if err := s.orderRepository.Update(ctx, orderUUID, order); err != nil {
		return fmt.Errorf("OrderService.Assemble update order error: %w", err)
	}

	return nil
}
