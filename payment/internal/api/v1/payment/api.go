package payment

import (
	"github.com/xgmsx/rsf/payment/internal/service"
	genPaymentV1 "github.com/xgmsx/rsf/shared/pkg/proto/payment/v1"
)

type paymentAPI struct {
	genPaymentV1.UnimplementedPaymentServiceServer

	service service.PaymentService
}

func NewPaymentAPI(service service.PaymentService) *paymentAPI {
	return &paymentAPI{service: service}
}
