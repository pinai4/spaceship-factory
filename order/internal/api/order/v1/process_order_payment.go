package v1

import (
	"context"
	"errors"

	"github.com/pinai4/spaceship-factory/order/internal/model"
	orderV1 "github.com/pinai4/spaceship-factory/shared/pkg/openapi/order/v1"
)

func (a *api) ProcessOrderPayment(
	ctx context.Context,
	req *orderV1.ProcessOrderPaymentRequest,
	params orderV1.ProcessOrderPaymentParams,
) (orderV1.ProcessOrderPaymentRes, error) {
	tranUUID, err := a.orderService.ProcessPayment(ctx, params.OrderUUID, model.OrderPaymentMethod(req.PaymentMethod))
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    404,
				Message: "Order with ID '" + params.OrderUUID.String() + "' not found",
			}, nil
		}
		return nil, err
	}

	return &orderV1.ProcessOrderPaymentResponse{
		TransactionUUID: tranUUID,
	}, nil
}
