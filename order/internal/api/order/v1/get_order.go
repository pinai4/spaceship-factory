package v1

import (
	"context"
	"errors"

	"github.com/pinai4/spaceship-factory/order/internal/converter"
	"github.com/pinai4/spaceship-factory/order/internal/model"
	orderV1 "github.com/pinai4/spaceship-factory/shared/pkg/openapi/order/v1"
)

func (a *api) GetOrder(ctx context.Context, params orderV1.GetOrderParams) (orderV1.GetOrderRes, error) {
	order, err := a.orderService.Get(ctx, params.OrderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    404,
				Message: "Order with ID '" + params.OrderUUID.String() + "' not found",
			}, nil
		}
		return nil, err
	}

	return converter.OrderToOpenAPI(order), nil
}
