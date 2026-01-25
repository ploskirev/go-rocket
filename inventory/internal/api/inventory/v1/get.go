package apiinventoryv1

import (
	"context"
	"errors"

	"github.com/ploskirev/go-rocket/inventory/internal/api/inventory/v1/converter"
	"github.com/ploskirev/go-rocket/inventory/internal/model"
	inventory_v1 "github.com/ploskirev/go-rocket/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) GetPart(ctx context.Context, req *inventory_v1.GetPartRequest) (*inventory_v1.Part, error) {
	part, err := a.ps.GetPart(ctx, req.Uuid)
	if err != nil {
		if errors.Is(err, model.ErrPartNotFound) {
			return nil, status.Errorf(codes.NotFound, "Part with UUID %s not found", req.GetUuid())
		}
	}

	return converter.PartToProto(*part), nil
}
