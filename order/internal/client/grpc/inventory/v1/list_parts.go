package inventoryv1

import (
	"context"
	"log"

	"github.com/ploskirev/go-rocket/order/internal/client/converter"
	"github.com/ploskirev/go-rocket/order/internal/model"
	inventory_v1 "github.com/ploskirev/go-rocket/shared/pkg/proto/inventory/v1"
)

func (c *inventoryClient) ListParts(ctx context.Context) ([]*model.Part, error) {
	res, err := (c.ic).ListParts(ctx, &inventory_v1.ListPartsRequest{})
	if err != nil {
		log.Printf("ERROR: Get part list from inventory service client: %s", err)
		return nil, err
	}

	partsList := make([]*model.Part, 0, len(res.Parts))

	for _, p := range res.Parts {
		partsList = append(partsList, converter.PartToModel(p))
	}

	return partsList, nil
}
