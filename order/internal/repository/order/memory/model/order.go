package model

import (
	"slices"

	"github.com/google/uuid"
)

type OrderPaymentMethod string

const (
	OrderPaymentMethodUnknown       OrderPaymentMethod = "UNKNOWN"
	OrderPaymentMethodCard          OrderPaymentMethod = "CARD"
	OrderPaymentMethodSBP           OrderPaymentMethod = "SBP"
	OrderPaymentMethodCreditCard    OrderPaymentMethod = "CREDIT_CARD"
	OrderPaymentMethodInvestorMoney OrderPaymentMethod = "INVESTOR_MONEY"
)

// AllValues returns all OrderPaymentMethod values.
func (OrderPaymentMethod) AllValues() []OrderPaymentMethod {
	return []OrderPaymentMethod{
		OrderPaymentMethodUnknown,
		OrderPaymentMethodCard,
		OrderPaymentMethodSBP,
		OrderPaymentMethodCreditCard,
		OrderPaymentMethodInvestorMoney,
	}
}

func (pm OrderPaymentMethod) IsValid() bool {
	return slices.Contains(pm.AllValues(), pm)
}

type OrderStatus string

const (
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid           OrderStatus = "PAID"
	OrderStatusCancelled      OrderStatus = "CANCELLED"
)

// AllValues returns all OrderStatus values.
func (OrderStatus) AllValues() []OrderStatus {
	return []OrderStatus{
		OrderStatusPendingPayment,
		OrderStatusPaid,
		OrderStatusCancelled,
	}
}

func (s OrderStatus) IsValid() bool {
	return slices.Contains(s.AllValues(), s)
}

type Order struct {
	UUID            uuid.UUID
	UserUUID        uuid.UUID
	PartUUIDs       []uuid.UUID
	TotalPrice      float64
	TransactionUUID *uuid.UUID
	PaymentMethod   *OrderPaymentMethod
	Status          OrderStatus
}
