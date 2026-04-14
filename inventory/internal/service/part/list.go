package part

import (
	"context"
	"log"

	"github.com/ploskirev/go-rocket/inventory/internal/model"
)

func (ps *partService) ListParts(ctx context.Context, filters *model.Filters) ([]*model.Part, error) {
	parts, err := ps.pr.ListParts(ctx, filters)
	if err != nil {
		log.Printf("Error: List part from repository \n")
		return nil, err
	}

	return parts, nil
}
