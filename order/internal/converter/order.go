package converter

import (
	"github.com/pinai4/spaceship-factory/order/internal/model"
	orderV1 "github.com/pinai4/spaceship-factory/shared/pkg/openapi/order/v1"
)

func OrderToOpenAPI(order model.Order) *orderV1.Order {
	var tranUUID orderV1.OptUUID
	if order.TransactionUUID != nil {
		tranUUID.SetTo(*order.TransactionUUID)
	}

	var mp orderV1.OptOrderPaymentMethod
	if order.PaymentMethod != nil {
		mp.SetTo(orderV1.OrderPaymentMethod(*order.PaymentMethod))
	}

	return &orderV1.Order{
		OrderUUID:       order.UUID,
		UserUUID:        order.UserUUID,
		PartUuids:       order.PartUUIDs,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: tranUUID,
		PaymentMethod:   mp,
		Status:          orderV1.OrderStatus(order.Status),
	}
}

func CreateOrderRequestToModel(req *orderV1.CreateOrderRequest) model.CreateOrder {
	return model.CreateOrder{
		UserUUID:  req.UserUUID,
		PartUUIDs: req.PartUuids,
	}
}
