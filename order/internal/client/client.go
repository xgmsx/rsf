package client

import (
	"context"

	"github.com/google/uuid"

	"github.com/xgmsx/rsf/order/internal/model"
	genInventoryV1 "github.com/xgmsx/rsf/shared/pkg/proto/inventory/v1"
)

type InventoryClient interface {
	GetParts(ctx context.Context, uuids []uuid.UUID) (parts []*genInventoryV1.Part, err error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, userUUID, orderUUID uuid.UUID, paymentMethod model.PaymentMethod) (txUUID *uuid.UUID, err error)
}
