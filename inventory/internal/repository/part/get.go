package part

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/ploskirev/go-rocket/inventory/internal/model"
)

func (pr *partRepository) GetPart(ctx context.Context, uuid string) (*model.Part, error) {
	var part model.Part
	err := pr.collection.FindOne(ctx, bson.M{"uuid": uuid}).Decode(&part)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return &model.Part{}, model.ErrPartNotFound
		}

		return &model.Part{}, err
	}

	return &part, nil
}
