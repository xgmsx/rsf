package model

import "github.com/google/uuid"

type PayOrderInput struct {
	OrderUUID     uuid.UUID
	PaymentMethod PaymentMethod
}

type PayOrderOutput struct {
	TransactionUUID string
}

type CreateOrderInput struct {
	UserUUID  uuid.UUID
	PartUuids []string
}

type CreateOrderOutput struct {
	OrderUUID  uuid.UUID
	TotalPrice float64
}
