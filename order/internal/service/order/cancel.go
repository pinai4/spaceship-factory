package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/pinai4/spaceship-factory/order/internal/model"
)

func (s *service) Cancel(ctx context.Context, orderUUID uuid.UUID) error {
	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		return fmt.Errorf("OrderService.Cancel get order error: %w", err)
	}

	if order.Status == model.OrderStatusPaid {
		return model.ErrOrderCancelNotAllowed
	}

	order.Status = model.OrderStatusCancelled

	if err := s.orderRepository.Update(ctx, orderUUID, order); err != nil {
		return fmt.Errorf("OrderService.Cancel update order error: %w", err)
	}

	return nil
}
