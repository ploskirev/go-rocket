package apiinventoryv1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ploskirev/go-rocket/inventory/internal/api/inventory/v1/converter"
	inventory_v1 "github.com/ploskirev/go-rocket/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, req *inventory_v1.ListPartsRequest) (*inventory_v1.ListPartsResponse, error) {
	parts, err := a.ps.ListParts(ctx, converter.FiltersToModel(req.Filter))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Part list get error: %s", err)
	}

	partsList := make([]*inventory_v1.Part, 0, len(parts))
	for _, p := range parts {
		partsList = append(partsList, converter.PartToProto(*p))
	}

	return &inventory_v1.ListPartsResponse{
		Parts: partsList,
	}, nil
}
