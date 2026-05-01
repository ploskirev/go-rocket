package app

import (
	"context"
	"fmt"
	"log"

	inventoryV1API "github.com/ploskirev/go-rocket/inventory/internal/api/inventory/v1"
	"github.com/ploskirev/go-rocket/inventory/internal/config"
	"github.com/ploskirev/go-rocket/platform/pkg/closer"
	inventoryV1 "github.com/ploskirev/go-rocket/shared/pkg/proto/inventory/v1"

	"github.com/ploskirev/go-rocket/inventory/internal/repository"
	inventoryRepository "github.com/ploskirev/go-rocket/inventory/internal/repository/part"
	"github.com/ploskirev/go-rocket/inventory/internal/service"
	inventoryService "github.com/ploskirev/go-rocket/inventory/internal/service/part"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type diContainer struct {
	inventoryV1API inventoryV1.InventoryServiceServer

	inventoryService service.PartService

	inventoryRepository repository.PartRepository

	mongoDBClient *mongo.Client
	mongoDBHandle *mongo.Database
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) InventoryV1API(ctx context.Context) inventoryV1.InventoryServiceServer {
	if d.inventoryV1API == nil {
		d.inventoryV1API = inventoryV1API.NewApi(d.InventoryService(ctx))
	}

	return d.inventoryV1API
}

func (d *diContainer) InventoryService(ctx context.Context) service.PartService {
	if d.inventoryService == nil {
		d.inventoryService = inventoryService.NewPartService(d.PartRepository(ctx))
	}

	return d.inventoryService
}

func (d *diContainer) PartRepository(ctx context.Context) repository.PartRepository {
	if d.inventoryRepository == nil {
		d.inventoryRepository = inventoryRepository.NewPartRepository(d.MongoDBHandle(ctx))
		err := d.inventoryRepository.InitStorage(ctx)
		if err != nil {
			log.Printf("failed to init repo %s \n", err)
		}
	}

	return d.inventoryRepository
}

func (d *diContainer) MongoDBClient(ctx context.Context) *mongo.Client {
	if d.mongoDBClient == nil {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
		if err != nil {
			panic(fmt.Sprintf("failed to connect to MongoDB: %s\n", err.Error()))
		}

		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			panic(fmt.Sprintf("failed to ping MongoDB: %v\n", err))
		}

		closer.AddNamed("MongoDB client", func(ctx context.Context) error {
			return client.Disconnect(ctx)
		})

		d.mongoDBClient = client
	}

	return d.mongoDBClient
}

func (d *diContainer) MongoDBHandle(ctx context.Context) *mongo.Database {
	if d.mongoDBHandle == nil {
		d.mongoDBHandle = d.MongoDBClient(ctx).Database(config.AppConfig().Mongo.DatabaseName())
	}

	return d.mongoDBHandle
}
