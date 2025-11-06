package model

import "github.com/google/uuid"

type OrderPaidEvent struct {
	EventUUID       string
	OrderUUID       uuid.UUID
	UserUUID        uuid.UUID
	PaymentMethod   string
	TransactionUUID uuid.UUID
}

type ShipAssembledEvent struct {
	EventUUID    string
	OrderUUID    uuid.UUID
	UserUUID     uuid.UUID
	BuildTimeSec int64
}
