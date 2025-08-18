package part

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/xgmsx/rsf/inventory/internal/api/v1/converter"
	"github.com/xgmsx/rsf/inventory/internal/model"
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
	filter := converter.PartFilterFromProto(req.Filter)
	parts, err := h.service.ListParts(ctx, filter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error: %v", err)
	}

	resp := &genInventoryV1.ListPartsResponse{
		Parts: make([]*genInventoryV1.Part, 0, len(parts)),
	}
	for _, part := range parts {
		resp.Parts = append(resp.Parts, converter.PartToProto(part))
	}
	return resp, nil
}
