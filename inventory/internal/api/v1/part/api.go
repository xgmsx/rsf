package part

import (
	"github.com/xgmsx/rsf/inventory/internal/service"
	inventoryV1 "github.com/xgmsx/rsf/shared/pkg/proto/inventory/v1"
)

type partAPI struct {
	inventoryV1.UnimplementedInventoryServiceServer

	service service.PartService
}

func NewPartAPI(service service.PartService) *partAPI {
	return &partAPI{
		service: service,
	}
}
