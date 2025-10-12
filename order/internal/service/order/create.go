package order

import (
	"context"
	"fmt"
	"slices"

	"github.com/google/uuid"

	"github.com/pinai4/spaceship-factory/order/internal/model"
)

func (s *service) Create(ctx context.Context, orderUUID uuid.UUID, createOrder model.CreateOrder) (float64, error) {
	orderedPartIDs := make([]string, len(createOrder.PartUUIDs))
	for i, p := range createOrder.PartUUIDs {
		orderedPartIDs[i] = p.String()
	}

	parts, err := s.inventoryClient.ListParts(ctx, orderedPartIDs)
	if err != nil {
		return 0, fmt.Errorf("OrderService.Create inventory API client error: %w", err)
	}
	if len(parts) != len(orderedPartIDs) {
		return 0, model.ErrOrderedPartsNotAvailable
	}

	var totalPrice float64
	for _, p := range parts {
		if !slices.Contains(orderedPartIDs, p.UUID) || p.StockQuantity == 0 {
			return 0, model.ErrOrderedPartsNotAvailable
		}
		totalPrice += p.Price
	}

	order := model.Order{
		UUID:       orderUUID,
		UserUUID:   createOrder.UserUUID,
		PartUUIDs:  createOrder.PartUUIDs,
		TotalPrice: totalPrice,
		Status:     model.OrderStatusPendingPayment,
	}

	if err := s.orderRepository.Create(ctx, order); err != nil {
		return 0, fmt.Errorf("OrderService.Create create order error: %w", err)
	}

	return totalPrice, nil
}
