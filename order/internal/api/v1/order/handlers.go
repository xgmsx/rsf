package order

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/xgmsx/rsf/order/internal/api/v1/converter"
	"github.com/xgmsx/rsf/order/internal/model"
	orderV1 "github.com/xgmsx/rsf/shared/pkg/openapi/order/v1"
)

// CreateOrder implements shared/pkg/openapi/order/v1.
func (h *orderApi) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	input := converter.CreateOrderInputFromRequest(*req)
	output, err := h.orderService.CreateOrder(ctx, input)
	if err != nil {
		if errors.Is(err, model.ErrPartDoesNotExist) {
			return &orderV1.BadRequestError{
				Code:    400,
				Message: err.Error(),
			}, nil
		}
		return nil, h.NewError(ctx, err)
	}

	return converter.CreateOrderOutputToResponse(output), nil
}

// GetOrder implements shared/pkg/openapi/order/v1.
func (h *orderApi) GetOrder(ctx context.Context, params orderV1.GetOrderParams) (orderV1.GetOrderRes, error) {
	output, err := h.orderService.GetOrder(ctx, params.OrderUUID.String())
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: "Order with UUID: '" + params.OrderUUID.String() + "' not found",
			}, nil
		}
		return nil, h.NewError(ctx, err)
	}

	return converter.GetOrderOutputToResponse(output), nil
}

// PayOrder implements shared/pkg/openapi/order/v1.
func (h *orderApi) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	input := converter.PayOrderInputFromRequest(*req, params)
	output, err := h.orderService.PayOrder(ctx, input)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: "Order with UUID: '" + params.OrderUUID.String() + "' not found",
			}, nil
		}
		if errors.Is(err, model.ErrPaymentMethodIsNotSupported) {
			return &orderV1.BadRequestError{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf("Payment method %v is not supported", req.PaymentMethod),
			}, nil
		}

		return &orderV1.InternalServerError{
			Code:    500,
			Message: err.Error(),
		}, err
	}

	return converter.PayOrderOutputToResponse(output), nil
}

// CancelOrder implements shared/pkg/openapi/order/v1.
func (h *orderApi) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	_, err := h.orderService.CancelOrder(ctx, params.OrderUUID.String())
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: "Order with UUID: '" + params.OrderUUID.String() + "' not found",
			}, nil
		}
		if errors.Is(err, model.ErrOrderAlreadyPaid) {
			return &orderV1.NotFoundError{
				Code:    http.StatusConflict,
				Message: "Order already paid and cannot be cancelled",
			}, nil
		}
		return nil, h.NewError(ctx, err)
	}

	return &orderV1.CancelOrderNoContent{}, nil
}
