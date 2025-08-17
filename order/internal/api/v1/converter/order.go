package converter

import (
	"github.com/google/uuid"

	"github.com/xgmsx/rsf/order/internal/model"
	orderV1 "github.com/xgmsx/rsf/shared/pkg/openapi/order/v1"
)

func CreateOrderInputFromRequest(request orderV1.CreateOrderRequest) model.CreateOrderInput {
	return model.CreateOrderInput{
		UserUUID:  request.UserUUID,
		PartUuids: request.GetPartUuids(),
	}
}

func CreateOrderOutputToResponse(output model.CreateOrderOutput) *orderV1.CreateOrderResponse {
	return &orderV1.CreateOrderResponse{
		OrderUUID:  output.OrderUUID,
		TotalPrice: float32(output.TotalPrice),
	}
}

func PayOrderInputFromRequest(request orderV1.PayOrderRequest, params orderV1.PayOrderParams) model.PayOrderInput {
	return model.PayOrderInput{
		OrderUUID:     params.OrderUUID,
		PaymentMethod: model.PaymentMethod(request.PaymentMethod),
	}
}

func PayOrderOutputToResponse(output model.PayOrderOutput) *orderV1.PayOrderResponse {
	return &orderV1.PayOrderResponse{
		TransactionUUID: uuid.MustParse(output.TransactionUUID),
	}
}

func GetOrderOutputToResponse(order model.Order) *orderV1.Order {
	res := orderV1.Order{
		OrderUUID:  order.OrderUUID,
		UserUUID:   order.UserUUID,
		PartUuids:  order.PartUUIDs,
		TotalPrice: order.TotalPrice,
		Status:     orderV1.OrderStatus(order.Status),
	}
	if order.PaymentMethod != nil {
		res.PaymentMethod = orderV1.NewOptNilOrderPaymentMethod(
			orderV1.OrderPaymentMethod(*order.PaymentMethod),
		)
	}
	if order.TransactionUUID != nil {
		res.TransactionUUID = orderV1.NewOptNilUUID(*order.TransactionUUID)
	}
	return &res
}
