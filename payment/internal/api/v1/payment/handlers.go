package payment

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/xgmsx/rsf/payment/internal/model"
	"github.com/xgmsx/rsf/payment/internal/model/converter"
	genPaymentV1 "github.com/xgmsx/rsf/shared/pkg/proto/payment/v1"
)

func (h *paymentAPI) PayOrder(ctx context.Context, req *genPaymentV1.PayOrderRequest) (*genPaymentV1.PayOrderResponse, error) {
	input := converter.PayInputFromRequest(req)
	output, err := h.service.PayOrder(ctx, input)
	if err != nil {
		if errors.Is(err, model.ErrInvalidPaymentMethod) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return converter.PayOutputToResponse(output), nil
}
