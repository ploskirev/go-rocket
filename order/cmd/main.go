package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/ploskirev/go-rocket/order/internal/config"

	orderapiv1 "github.com/ploskirev/go-rocket/order/internal/api/order/v1"
	inventoryclient "github.com/ploskirev/go-rocket/order/internal/client/grpc/inventory/v1"
	paymentclient "github.com/ploskirev/go-rocket/order/internal/client/grpc/payment/v1"
	"github.com/ploskirev/go-rocket/order/internal/migrator"
	orderrepo "github.com/ploskirev/go-rocket/order/internal/repository/order"
	orderservice "github.com/ploskirev/go-rocket/order/internal/service/order"
	order_v1 "github.com/ploskirev/go-rocket/shared/pkg/openapi/order/v1"
)

const (
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
	configPath        = "./deploy/compose/order/.env"
)

func main() {
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	ctx := context.Background()

	dbURI := config.AppConfig().Postgres.URI()

	// Создаем пул соединений с базой данных
	pool, err := pgxpool.New(ctx, dbURI)
	if err != nil {
		log.Printf("failed to connect to database: %v\n", err)
		return
	}
	defer pool.Close()

	migrationsDir := config.AppConfig().Postgres.MigrationsDir()
	fmt.Println("MIGR DIR: ", dbURI)
	migratorRunner := migrator.NewMigrator(stdlib.OpenDBFromPool(pool), migrationsDir)

	err = migratorRunner.Up()
	if err != nil {
		log.Printf("Ошибка миграции базы данных: %v\n", err)
		return
	}

	inventoryClient, inventoryConn, err := inventoryclient.NewInventoryClient()
	if err != nil {
		log.Printf("failed to connect inventory service: %v\n", err)
		return
	}
	defer func() {
		if cerr := inventoryConn.Close(); cerr != nil {
			log.Printf("failed to close connect inventory service: %v", cerr)
		}
	}()
	paymentClient, paymentConn, err := paymentclient.NewPaymentClient()
	if err != nil {
		log.Printf("failed to connect payment service: %v\n", err)
		return
	}
	defer func() {
		if cerr := paymentConn.Close(); cerr != nil {
			log.Printf("failed to close connect payment service: %v", cerr)
		}
	}()

	orderRepo := orderrepo.NewOrderRepo(pool)
	orderService := orderservice.NewOrderService(inventoryClient, paymentClient, orderRepo)
	api := orderapiv1.NewApi(orderService)

	orderServer, err := order_v1.NewServer(api)
	if err != nil {
		log.Printf("ошибка создания Orders сервера: %v", err)
	}

	// Инициализируем роутер Chi
	r := chi.NewRouter()

	// Добавляем middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
	// r.Use(customMiddleware.RequestLogger)

	// Монтируем обработчики OpenAPI
	r.Mount("/", orderServer)

	// Запускаем HTTP-сервер
	server := &http.Server{
		Addr:              config.AppConfig().OrderHTTP.Address(),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout, // Защита от Slowloris атак - тип DDoS-атаки, при которой
		// атакующий умышленно медленно отправляет HTTP-заголовки, удерживая соединения открытыми и истощая
		// пул доступных соединений на сервере. ReadHeaderTimeout принудительно закрывает соединение,
		// если клиент не успел отправить все заголовки за отведенное время.
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		log.Printf("🚀 HTTP-сервер запущен gо адресу %s\n", config.AppConfig().OrderHTTP.Address())
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	// Создаем контекст с таймаутом для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
