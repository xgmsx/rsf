package order

import (
	"context"
	"net/http"

	"github.com/xgmsx/rsf/order/internal/service"
	genOrderV1 "github.com/xgmsx/rsf/shared/pkg/openapi/order/v1"
)

var _ genOrderV1.Handler = (*orderApi)(nil)

type orderApi struct {
	orderService service.OrderService
}

func NewOrderAPI(orderService service.OrderService) *orderApi {
	return &orderApi{
		orderService: orderService,
	}
}

// NewError создает новую ошибку в формате GenericError
func (h *orderApi) NewError(_ context.Context, err error) *genOrderV1.GenericErrorStatusCode {
	return &genOrderV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: genOrderV1.GenericError{
			Code:    genOrderV1.NewOptInt(http.StatusInternalServerError),
			Message: genOrderV1.NewOptString(err.Error()),
		},
	}
}
