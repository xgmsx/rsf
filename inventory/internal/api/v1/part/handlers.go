package part

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/xgmsx/rsf/inventory/internal/model"
	"github.com/xgmsx/rsf/inventory/internal/model/converter"
	genInventoryV1 "github.com/xgmsx/rsf/shared/pkg/proto/inventory/v1"
)

func (h *partAPI) GetPart(ctx context.Context, req *genInventoryV1.GetPartRequest) (*genInventoryV1.GetPartResponse, error) {
	part, err := h.service.GetPart(ctx, req.GetUuid())
	if err != nil {
		if errors.Is(err, model.ErrPartDoesNotExist) {
			return nil, status.Errorf(codes.NotFound, "part not found")
		}
		return nil, status.Errorf(codes.Internal, "internal error: %v", err)
	}

	return &genInventoryV1.GetPartResponse{
		Part: converter.PartToProto(part),
	}, nil
}

func (h *partAPI) ListParts(ctx context.Context, req *genInventoryV1.ListPartsRequest) (*genInventoryV1.ListPartsResponse, error) {
	err := req.ValidateAll()
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validate error: %v", err)
	}

	parts, err := h.service.ListParts(ctx, converter.PartFilterFromProto(req.Filter))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error: %v", err)
	}

	if len(parts) == 0 {
		return nil, status.Errorf(codes.NotFound, "not found error")
	}

	resp := &genInventoryV1.ListPartsResponse{
		Parts: converter.PartsToProto(parts),
	}

	err = resp.Validate()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error: %v", err)
	}

	return resp, nil
}
