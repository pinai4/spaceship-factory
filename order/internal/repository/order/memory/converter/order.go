package converter

import (
	"github.com/pinai4/spaceship-factory/order/internal/model"
	repoModel "github.com/pinai4/spaceship-factory/order/internal/repository/order/memory/model"
)

func OrderToModel(order repoModel.Order) model.Order {
	var pm *model.OrderPaymentMethod
	if order.PaymentMethod != nil {
		m := model.OrderPaymentMethod(*order.PaymentMethod)
		pm = &m
	}

	return model.Order{
		UUID:            order.UUID,
		UserUUID:        order.UserUUID,
		PartUUIDs:       order.PartUUIDs,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   pm,
		Status:          model.OrderStatus(order.Status),
	}
}

func OrderToRepoModel(order model.Order) repoModel.Order {
	var pm *repoModel.OrderPaymentMethod
	if order.PaymentMethod != nil {
		m := repoModel.OrderPaymentMethod(*order.PaymentMethod)
		pm = &m
	}

	return repoModel.Order{
		UUID:            order.UUID,
		UserUUID:        order.UserUUID,
		PartUUIDs:       order.PartUUIDs,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   pm,
		Status:          repoModel.OrderStatus(order.Status),
	}
}
