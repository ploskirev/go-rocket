//go:tag integration

package integration

import (
	"context"
	"os"

	"github.com/brianvoe/gofakeit/v7"
	"go.mongodb.org/mongo-driver/bson"

	inventory_v1 "github.com/ploskirev/go-rocket/shared/pkg/proto/inventory/v1"
)

// InsertTestPart — вставляет тестовeю деталь в коллекцию Mongo и возвращает ее UUID
func (env *TestEnvironment) InsertTestPart(ctx context.Context) (string, error) {
	partUUID := gofakeit.UUID()
	// now := time.Now()

	part := bson.M{
		"_id": partUUID,
		"info": bson.M{
			"name":  gofakeit.Name(),
			"price": gofakeit.Price(1, 1000),
		},
		// "created_at": primitive.NewDateTimeFromTime(now),
	}

	// Используем базу данных из переменной окружения MONGO_DATABASE
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "inventory" // fallback значение
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(collectionName).InsertOne(ctx, part)
	if err != nil {
		return "", err
	}

	return partUUID, nil
}

// InsertTestPartWithData — вставляет тестовeю деталь с заданными данными в коллекцию Mongo и возвращает ее UUID
func (env *TestEnvironment) InsertTestPartWithData(ctx context.Context, info *inventory_v1.Part) (string, error) {
	partUUID := gofakeit.UUID()
	// now := time.Now()

	// observedAt := info.GetObservedAt().AsTime()

	part := bson.M{
		"_id": partUUID,
		"info": bson.M{
			"name":  info.Name,
			"price": info.Price,
		},
		// "created_at": primitive.NewDateTimeFromTime(now),
	}

	// Используем базу данных из переменной окружения MONGO_DATABASE
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "inventory" // fallback значение
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(collectionName).InsertOne(ctx, part)
	if err != nil {
		return "", err
	}

	return partUUID, nil
}

// GetTestPartInfo — возвращает тестовую информацию о детали
func (env *TestEnvironment) GetTestPartInfo() *inventory_v1.Part {
	return &inventory_v1.Part{
		Name:  "Test part",
		Price: 150,
	}
}

// GetUpdatedPartInfo — возвращает обновленную информацию о детали
func (env *TestEnvironment) GetUpdatedPartInfo() *inventory_v1.Part {
	return &inventory_v1.Part{
		Name:  "Test part upd",
		Price: 151,
	}
}

// ClearPartsCollection — удаляет все записи из коллекции parts
func (env *TestEnvironment) ClearPartsCollection(ctx context.Context) error {
	// Используем базу данных из переменной окружения MONGO_DATABASE
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "inventory" // fallback значение
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(collectionName).DeleteMany(ctx, bson.M{})
	if err != nil {
		return err
	}

	return nil
}
