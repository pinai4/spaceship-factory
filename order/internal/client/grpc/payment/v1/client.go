package v1

import (
	def "github.com/pinai4/spaceship-factory/order/internal/client"
	paymentV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/payment/v1"
)

var _ def.PaymentClient = (*client)(nil)

type client struct {
	generatedClient paymentV1.PaymentServiceClient
}

func NewClient(paymentClient paymentV1.PaymentServiceClient) *client {
	return &client{generatedClient: paymentClient}
}
