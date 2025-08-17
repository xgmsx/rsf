package model

type PayOrderInput struct {
	OrderID       string
	PaymentMethod PaymentMethod
}

type PayOrderOutput struct {
	TransactionUUID string
}
