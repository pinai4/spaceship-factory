package v1

import (
	"context"
	"fmt"

	paymentV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/payment/v1"
)

func (c *client) PayOrder(ctx context.Context, orderUUID, userUUID, paymentMethod string) (string, error) {
	parsePaymentMethod := func(s string) paymentV1.PaymentMethod {
		if val, ok := paymentV1.PaymentMethod_value[s]; ok {
			return paymentV1.PaymentMethod(val)
		}
		return paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}

	resp, err := c.generatedClient.PayOrder(ctx, &paymentV1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: parsePaymentMethod(fmt.Sprintf("PAYMENT_METHOD_%s", paymentMethod)),
	})
	if err != nil {
		return "", fmt.Errorf("client API call error: %w", err)
	}

	return resp.GetTransactionUuid(), nil
}
