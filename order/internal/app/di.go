package app

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/ploskirev/go-rocket/order/internal/config"
	"github.com/ploskirev/go-rocket/order/internal/migrator"
	"github.com/ploskirev/go-rocket/platform/pkg/closer"

	"github.com/ploskirev/go-rocket/order/internal/repository"
	orderRepository "github.com/ploskirev/go-rocket/order/internal/repository/order"
	"github.com/ploskirev/go-rocket/order/internal/service"
	orderService "github.com/ploskirev/go-rocket/order/internal/service/order"

	inventoryclient "github.com/ploskirev/go-rocket/order/internal/client/grpc/inventory/v1"
	paymentclient "github.com/ploskirev/go-rocket/order/internal/client/grpc/payment/v1"

	orderapiv1 "github.com/ploskirev/go-rocket/order/internal/api/order/v1"
)

type diContainer struct {
	orderV1API      *orderapiv1.Api
	orderService    service.OrderService
	orderRepository repository.OrderRepository
	postgresPool    *pgxpool.Pool
	inventoryClient *inventoryclient.InventoryClient
	paymentClient   *paymentclient.PaymentClient
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) OrderV1API(ctx context.Context) *orderapiv1.Api {
	if d.orderV1API == nil {
		d.orderV1API = orderapiv1.NewApi(d.OrderService(ctx))
	}

	return d.orderV1API
}

func (d *diContainer) OrderService(ctx context.Context) service.OrderService {
	if d.orderService == nil {
		d.orderService = orderService.NewOrderService(d.InventoryClient(ctx), d.PaymentClient(ctx), d.OrderRepository(ctx))
	}

	return d.orderService
}

func (d *diContainer) InventoryClient(ctx context.Context) *inventoryclient.InventoryClient {
	if d.inventoryClient == nil {
		inventoryClient, inventoryConn, err := inventoryclient.NewInventoryClient()
		if err != nil {
			panic(fmt.Sprintf("failed to connect inventory service: %v\n", err))
		}

		d.inventoryClient = inventoryClient

		closer.AddNamed("Inventory client", func(_ context.Context) error {
			return inventoryConn.Close()
		})
	}

	return d.inventoryClient
}

func (d *diContainer) PaymentClient(ctx context.Context) *paymentclient.PaymentClient {
	if d.paymentClient == nil {
		paymentClient, paymentConn, err := paymentclient.NewPaymentClient()
		if err != nil {
			panic(fmt.Sprintf("failed to connect payment service: %v\n", err))
		}

		d.paymentClient = paymentClient

		closer.AddNamed("Payment client", func(_ context.Context) error {
			return paymentConn.Close()
		})
	}

	return d.paymentClient
}

func (d *diContainer) OrderRepository(ctx context.Context) repository.OrderRepository {
	if d.orderRepository == nil {
		d.orderRepository = orderRepository.NewOrderRepo(d.PostgresPool(ctx))
	}

	return d.orderRepository
}

func (d *diContainer) PostgresPool(ctx context.Context) *pgxpool.Pool {
	if d.postgresPool == nil {
		dbURI := config.AppConfig().Postgres.URI()
		pool, err := pgxpool.New(ctx, dbURI)
		if err != nil {
			panic(fmt.Sprintf("failed to get pgxpool: %v\n", err))
		}
		d.postgresPool = pool

		closer.AddNamed("MongoDB client", func(_ context.Context) error {
			pool.Close()
			return nil
		})

		migrationsDir := config.AppConfig().Postgres.MigrationsDir()
		migratorRunner := migrator.NewMigrator(stdlib.OpenDBFromPool(pool), migrationsDir)

		err = migratorRunner.Up()
		if err != nil {
			log.Printf("Ошибка миграции базы данных: %v\n", err)
		}
	}

	return d.postgresPool
}
