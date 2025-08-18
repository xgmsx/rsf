package model

import "github.com/google/uuid"

type PayOrderInput struct {
	OrderUUID     uuid.UUID
	PaymentMethod PaymentMethod
}

type PayOrderOutput struct {
	TransactionUUID uuid.UUID
}

type CreateOrderInput struct {
	UserUUID  uuid.UUID
	PartUUIDs []uuid.UUID
}

type CreateOrderOutput struct {
	OrderUUID  uuid.UUID
	TotalPrice float64
}
