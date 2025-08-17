package repository

import (
	"context"

	"github.com/xgmsx/rsf/order/internal/model"
)

type OrderRepository interface {
	Get(ctx context.Context, orderUUID string) (model.Order, error)
	Update(ctx context.Context, order model.Order) error
}
