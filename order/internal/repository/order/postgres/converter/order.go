package converter

import (
	"database/sql"

	"github.com/google/uuid"

	"github.com/pinai4/spaceship-factory/order/internal/model"
	repoModel "github.com/pinai4/spaceship-factory/order/internal/repository/order/postgres/model"
)

func OrderToModel(order repoModel.Order) model.Order {
	var pm *model.OrderPaymentMethod
	if order.PaymentMethod.Valid {
		v := model.OrderPaymentMethod(order.PaymentMethod.String)
		pm = &v
	}

	var tranUUID *uuid.UUID
	if order.TransactionID.Valid {
		tranUUID = &order.TransactionID.UUID
	}

	return model.Order{
		UUID:            order.ID,
		UserUUID:        order.UserID,
		PartUUIDs:       order.PartIDs,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: tranUUID,
		PaymentMethod:   pm,
		Status:          model.OrderStatus(order.Status),
	}
}

func OrderToRepoModel(order model.Order) repoModel.Order {
	var pm sql.NullString
	if order.PaymentMethod != nil {
		pm = sql.NullString{String: string(*order.PaymentMethod), Valid: true}
	}

	var tranID uuid.NullUUID
	if order.TransactionUUID != nil {
		tranID = uuid.NullUUID{UUID: *order.TransactionUUID, Valid: true}
	}

	return repoModel.Order{
		ID:            order.UUID,
		UserID:        order.UserUUID,
		PartIDs:       order.PartUUIDs,
		TotalPrice:    order.TotalPrice,
		TransactionID: tranID,
		PaymentMethod: pm,
		Status:        string(order.Status),
	}
}
