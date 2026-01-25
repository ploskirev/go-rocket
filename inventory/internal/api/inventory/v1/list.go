package apiinventoryv1

import (
	"context"

	"github.com/ploskirev/go-rocket/inventory/internal/api/inventory/v1/converter"
	inventory_v1 "github.com/ploskirev/go-rocket/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, req *inventory_v1.ListPartsRequest) (*inventory_v1.ListPartsResponse, error) {
	parts := a.ps.ListParts(ctx, converter.FiltersToModel(req.Filter))
	partsList := make([]*inventory_v1.Part, 0, len(parts))
	for _, p := range parts {
		partsList = append(partsList, converter.PartToProto(*p))
	}

	return &inventory_v1.ListPartsResponse{
		Parts: partsList,
	}, nil
}
