package v1

import (
	"context"

	paymentV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, _ *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	tranID, err := a.paymentService.PayOrder(ctx)
	if err != nil {
		return nil, err
	}

	return &paymentV1.PayOrderResponse{TransactionUuid: tranID}, nil
}
