package part

import (
	"context"

	"github.com/xgmsx/rsf/inventory/internal/model"
	"github.com/xgmsx/rsf/inventory/internal/repository"
	def "github.com/xgmsx/rsf/inventory/internal/service"
)

var _ def.PartService = (*partService)(nil)

type partService struct {
	repository repository.PartRepository
}

func NewPartService(repository repository.PartRepository) *partService {
	return &partService{
		repository: repository,
	}
}

func (r *partService) ListParts(ctx context.Context, filter *model.PartsFilter) ([]model.Part, error) {
	parts, err := r.repository.ListParts(ctx, filter)
	if err != nil {
		return nil, err
	}
	return parts, nil
}

func (r *partService) GetPart(ctx context.Context, uuid string) (model.Part, error) {
	part, err := r.repository.GetPart(ctx, uuid)
	if err != nil {
		return model.Part{}, err
	}
	return part, err
}
