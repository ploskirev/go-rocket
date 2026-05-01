package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/ploskirev/go-rocket/order/internal/config"
	"github.com/ploskirev/go-rocket/platform/pkg/closer"
	"github.com/ploskirev/go-rocket/platform/pkg/logger"
	order_v1 "github.com/ploskirev/go-rocket/shared/pkg/openapi/order/v1"
)

const (
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
)

type App struct {
	diContainer *diContainer
	httpServer  *http.Server
	listener    net.Listener
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	return a.runHTTPServer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
		a.initHTTPServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	return logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJson(),
	)
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	api := a.diContainer.OrderV1API(ctx)
	orderServer, err := order_v1.NewServer(api)
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("ошибка создания Orders сервера: %v", err))
		return err
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

	a.httpServer = server

	return nil
}

func (a *App) runHTTPServer(ctx context.Context) error {
	// Запускаем сервер в отдельной горутине
	go func() {
		logger.Info(ctx, fmt.Sprintf("🚀 HTTP Order server started on %s", config.AppConfig().OrderHTTP.Address()))
		err := a.httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error(ctx, fmt.Sprintf("❌ Ошибка запуска сервера: %v\n", err))
			panic(fmt.Sprintf("❌ Ошибка запуска сервера: %v\n", err))
		}
	}()

	closer.AddNamed("HTTP server", func(ctx context.Context) error {
		err := a.httpServer.Shutdown(ctx)
		if err != nil {
			logger.Error(ctx, fmt.Sprintf("❌ Ошибка при остановке сервера: %v\n", err))
			return err
		}

		return nil
	})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info(ctx, "🛑 Завершение работы сервера...")

	return nil
}
