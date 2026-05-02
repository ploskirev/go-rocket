package part

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type partRepository struct {
	collection *mongo.Collection
}

// func NewPartRepository(db *mongo.Database) *partRepository {
// 	collection := db.Collection("parts")

// 	return &partRepository{
// 		collection: collection,
// 	}
// }

func NewPartRepository(db *mongo.Database) *partRepository {

	return &partRepository{
		collection: nil,
	}
}
