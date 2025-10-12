package v1

import (
	"context"
	"errors"

	"github.com/pinai4/spaceship-factory/order/internal/model"
	orderV1 "github.com/pinai4/spaceship-factory/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	if err := a.orderService.Cancel(ctx, params.OrderUUID); err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    404,
				Message: "Order with ID '" + params.OrderUUID.String() + "' not found",
			}, nil
		}
		if errors.Is(err, model.ErrOrderCancelNotAllowed) {
			return &orderV1.ConflictError{
				Code:    409,
				Message: model.ErrOrderCancelNotAllowed.Error(),
			}, nil
		}
		return nil, err
	}

	return &orderV1.CancelOrderNoContent{}, nil
}
