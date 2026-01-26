package service

import (
	"context"

	"github.com/ploskirev/go-rocket/inventory/internal/model"
)

type PartService interface {
	GetPart(ctx context.Context, UUID string) (*model.Part, error)
	ListParts(ctx context.Context, filters *model.Filters) []*model.Part
}
