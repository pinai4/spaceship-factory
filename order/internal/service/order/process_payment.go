package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/pinai4/spaceship-factory/order/internal/model"
)

func (s *service) ProcessPayment(ctx context.Context, orderUUID uuid.UUID, paymentMethod model.OrderPaymentMethod) (uuid.UUID, error) {
	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("OrderService.ProcessPayment get order error: %w", err)
	}

	tranID, err := s.paymentClient.PayOrder(ctx, orderUUID.String(), order.UserUUID.String(), string(paymentMethod))
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("OrderService.ProcessPayment payment API client error: %w", err)
	}

	tranUUID, err := uuid.Parse(tranID)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("OrderService.ProcessPayment transaction uuid parse error: %w", err)
	}

	order.TransactionUUID = &tranUUID
	order.PaymentMethod = &paymentMethod
	order.Status = model.OrderStatusPaid

	if err := s.orderRepository.Update(ctx, orderUUID, order); err != nil {
		return uuid.UUID{}, fmt.Errorf("OrderService.ProcessPayment update order error: %w", err)
	}

	return tranUUID, nil
}
