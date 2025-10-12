package model

import "errors"

var (
	ErrOrderNotFound            = errors.New("order not found")
	ErrOrderAlreadyExists       = errors.New("order already exists")
	ErrOrderedPartsNotAvailable = errors.New("ordered parts are not available")
	ErrOrderCancelNotAllowed    = errors.New("order can't be canceled, it was already paid")
)
