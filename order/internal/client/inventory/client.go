package inventory

import (
	"context"
	"log"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	def "github.com/xgmsx/rsf/order/internal/client"
	genInventoryV1 "github.com/xgmsx/rsf/shared/pkg/proto/inventory/v1"
)

var _ def.InventoryClient = (*client)(nil)

type client struct {
	generatedClient genInventoryV1.InventoryServiceClient
}

func NewClient(addr string) *client {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect to %s: %v\n", addr, err)
	}

	return &client{
		generatedClient: genInventoryV1.NewInventoryServiceClient(conn),
	}
}

func (c *client) GetParts(ctx context.Context, uuids []uuid.UUID) ([]*genInventoryV1.Part, error) {
	uuidsStr := make([]string, len(uuids))
	for i, uid := range uuids {
		uuidsStr[i] = uid.String()
	}

	partsFilter := &genInventoryV1.PartsFilter{
		Uuids: uuidsStr,
	}

	res, err := c.generatedClient.ListParts(ctx, &genInventoryV1.ListPartsRequest{
		Filter: partsFilter,
	})
	if err != nil {
		return nil, err
	}

	return res.Parts, nil
}
