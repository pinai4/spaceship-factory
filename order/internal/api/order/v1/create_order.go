package v1

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/pinai4/spaceship-factory/order/internal/converter"
	"github.com/pinai4/spaceship-factory/order/internal/model"
	orderV1 "github.com/pinai4/spaceship-factory/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	if len(req.PartUuids) == 0 {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "parts list is empty",
		}, nil
	}

	orderUUID := uuid.New()
	totalPrice, err := a.orderService.Create(ctx, orderUUID, converter.CreateOrderRequestToModel(req))
	if err != nil {
		if errors.Is(err, model.ErrOrderedPartsNotAvailable) {
			return &orderV1.BadRequestError{
				Code:    400,
				Message: model.ErrOrderedPartsNotAvailable.Error(),
			}, nil
		}
		return nil, err
	}

	return &orderV1.CreateOrderResponse{
		OrderUUID:  orderUUID,
		TotalPrice: totalPrice,
	}, nil
}
