package repository

import (
	"context"

	"github.com/ploskirev/go-rocket/inventory/internal/model"
)

type PartRepository interface {
	GetPart(_ context.Context, uuid string) (*model.Part, error)
	ListParts(_ context.Context, filters *model.Filters) []*model.Part
}
