package service

import (
	"context"

	"github.com/xgmsx/rsf/order/internal/model"
)

type OrderService interface {
	CancelOrder(ctx context.Context, orderUUID string) (model.Order, error)
	PayOrder(ctx context.Context, request model.PayOrderInput) (model.PayOrderOutput, error)
	CreateOrder(ctx context.Context, order model.CreateOrderInput) (model.CreateOrderOutput, error)
	GetOrder(ctx context.Context, orderUUID string) (model.Order, error)
}
