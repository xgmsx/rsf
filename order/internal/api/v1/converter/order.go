package converter

import (
	"github.com/xgmsx/rsf/order/internal/model"
	genOrderV1 "github.com/xgmsx/rsf/shared/pkg/openapi/order/v1"
)

func CreateOrderInputFromRequest(request genOrderV1.CreateOrderRequest) model.CreateOrderInput {
	return model.CreateOrderInput{
		UserUUID:  request.UserUUID,
		PartUUIDs: request.PartUuids,
	}
}

func CreateOrderOutputToResponse(output model.CreateOrderOutput) *genOrderV1.CreateOrderResponse {
	return &genOrderV1.CreateOrderResponse{
		OrderUUID:  output.OrderUUID,
		TotalPrice: float32(output.TotalPrice),
	}
}

func PayOrderInputFromRequest(request genOrderV1.PayOrderRequest, params genOrderV1.PayOrderParams) model.PayOrderInput {
	return model.PayOrderInput{
		OrderUUID:     params.OrderUUID,
		PaymentMethod: model.PaymentMethod(request.PaymentMethod),
	}
}

func PayOrderOutputToResponse(output model.PayOrderOutput) *genOrderV1.PayOrderResponse {
	return &genOrderV1.PayOrderResponse{
		TransactionUUID: output.TransactionUUID,
	}
}

func GetOrderOutputToResponse(order model.Order) *genOrderV1.Order {
	res := genOrderV1.Order{
		OrderUUID:  order.OrderUUID,
		UserUUID:   order.UserUUID,
		PartUuids:  order.PartUUIDs,
		TotalPrice: order.TotalPrice,
		Status:     genOrderV1.OrderStatus(order.Status),
	}
	if order.PaymentMethod != nil {
		res.PaymentMethod = genOrderV1.NewOptNilOrderPaymentMethod(
			genOrderV1.OrderPaymentMethod(*order.PaymentMethod),
		)
	}
	if order.TransactionUUID != nil {
		res.TransactionUUID = genOrderV1.NewOptNilUUID(*order.TransactionUUID)
	}
	return &res
}
