package part

import (
	"context"

	"github.com/ploskirev/go-rocket/inventory/internal/model"
)

func (ps *partService) ListParts(ctx context.Context, filters *model.Filters) []*model.Part {
	return ps.pr.ListParts(ctx, filters)
}
