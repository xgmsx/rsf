package model

import "errors"

var (
	ErrOrderNotFound    = errors.New("order not found")
	ErrOrderAlreadyPaid = errors.New("order already paid")

	ErrFailedToFetchInventory = errors.New("error while fetching inventory")
	ErrPartDoesNotExist       = errors.New("part does not exist")

	ErrPaymentMethodIsNotSupported = errors.New("payment method is not supported")
	ErrFailedToProcessPayment      = errors.New("failed to process payment")
)
