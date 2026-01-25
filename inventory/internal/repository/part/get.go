package part

import (
	"context"
	"log"

	"github.com/ploskirev/go-rocket/inventory/internal/model"
	"github.com/ploskirev/go-rocket/inventory/internal/repository/converter"
)

func (pr *partRepository) GetPart(_ context.Context, uuid string) (*model.Part, error) {
	pr.mu.RLock()
	defer pr.mu.RUnlock()

	part, ok := pr.storage[uuid]
	if !ok {
		log.Printf("Part with UUID %s not found\n", uuid)
		return nil, model.ErrPartNotFound
	}

	return converter.PartToModel(part), nil
}
