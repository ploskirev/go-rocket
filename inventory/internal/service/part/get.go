package part

import (
	"context"
	"log"

	"github.com/ploskirev/go-rocket/inventory/internal/model"
)

func (ps *partService) GetPart(ctx context.Context, uuid string) (*model.Part, error) {
	part, err := ps.pr.GetPart(ctx, uuid)
	if err != nil {
		log.Printf("Error: Get part from repository with uuid: %s \n", uuid)
		return nil, err
	}

	return part, nil
}
