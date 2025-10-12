package client

import "context"

type PaymentClient interface {
	// PayOrder processes order payment returns payment transaction uuid
	PayOrder(ctx context.Context, OrderUUID, UserUUID, PaymentMethod string) (string, error)
}
