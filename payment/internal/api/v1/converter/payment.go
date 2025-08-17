package converter

import (
	"github.com/xgmsx/rsf/payment/internal/model"
	paymentV1 "github.com/xgmsx/rsf/shared/pkg/proto/payment/v1"
)

func PayInputFromRequest(request *paymentV1.PayOrderRequest) model.PayOrderInput {
	return model.PayOrderInput{
		OrderID:       request.GetOrderUuid(),
		PaymentMethod: model.PaymentMethod(request.GetPaymentMethod()),
	}
}

func PayOutputToResponse(output model.PayOrderOutput) *paymentV1.PayOrderResponse {
	return &paymentV1.PayOrderResponse{
		TransactionUuid: output.TransactionUUID,
	}
}
