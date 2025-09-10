package converter

import (
	"github.com/xgmsx/rsf/payment/internal/model"
	genPaymentV1 "github.com/xgmsx/rsf/shared/pkg/proto/payment/v1"
)

func PayInputFromRequest(request *genPaymentV1.PayOrderRequest) model.PayOrderInput {
	return model.PayOrderInput{
		OrderID:       request.GetOrderUuid(),
		PaymentMethod: model.PaymentMethod(request.GetPaymentMethod()),
	}
}

func PayOutputToResponse(output model.PayOrderOutput) *genPaymentV1.PayOrderResponse {
	return &genPaymentV1.PayOrderResponse{
		TransactionUuid: output.TransactionUUID,
	}
}
