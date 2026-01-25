package part

import (
	"context"
	"log"

	"github.com/ploskirev/go-rocket/inventory/internal/model"
)

func (ps *partService) GetPart(ctx context.Context, UUID string) (*model.Part, error) {
	part, err := ps.pr.GetPart(ctx, UUID)
	if err != nil {
		log.Printf("Error: Get part from repository with uuid: %s \n", UUID)
		return nil, err
	}

	return part, nil
}
