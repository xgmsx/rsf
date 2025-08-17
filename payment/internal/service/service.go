package service

import (
	"context"

	"github.com/xgmsx/rsf/payment/internal/model"
)

type PaymentService interface {
	PayOrder(ctx context.Context, input model.PayOrderInput) (model.PayOrderOutput, error)
}
