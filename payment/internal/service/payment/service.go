package payment

import (
	"context"

	"github.com/google/uuid"

	"github.com/xgmsx/rsf/payment/internal/model"
	"github.com/xgmsx/rsf/payment/internal/service"
)

var _ service.PaymentService = (*paymentService)(nil)

type paymentService struct{}

func NewService() *paymentService {
	return &paymentService{}
}

func (s *paymentService) PayOrder(ctx context.Context, input model.PayOrderInput) (model.PayOrderOutput, error) {
	if input.PaymentMethod == model.PaymentMethod_UNSPECIFIED {
		return model.PayOrderOutput{}, model.ErrInvalidPaymentMethod
	}
	transactionUUID := uuid.New()
	return model.PayOrderOutput{
		TransactionUUID: transactionUUID.String(),
	}, nil
}
