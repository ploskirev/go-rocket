package part

import (
	"context"

	"github.com/ploskirev/go-rocket/inventory/internal/repository/model"
)

func (pr *partRepository) InitStorage(ctx context.Context) error {
	part1 := &model.Part{
		UUID:  "q1w2e3r4t5",
		Name:  "Test detail",
		Price: 535.7,
	}
	part2 := &model.Part{
		UUID:  "z9x8c7v6b5",
		Name:  "Second detail",
		Price: 177.5,
	}

	_, err := pr.collection.InsertOne(ctx, part1)
	if err != nil {
		return err
	}
	_, err = pr.collection.InsertOne(ctx, part2)
	if err != nil {
		return err
	}

	return nil
}
