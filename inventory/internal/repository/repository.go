package repository

import (
	"context"

	"github.com/ploskirev/go-rocket/inventory/internal/model"
)

type PartRepository interface {
	GetPart(ctx context.Context, uuid string) (*model.Part, error)
	ListParts(ctx context.Context, filters *model.Filters) ([]*model.Part, error)
}
