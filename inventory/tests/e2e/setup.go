//go:tag integration

package integration

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"

	"github.com/docker/go-connections/nat"
	"github.com/ploskirev/go-rocket/platform/pkg/logger"
	"github.com/ploskirev/go-rocket/platform/pkg/testcontainers"
	"github.com/ploskirev/go-rocket/platform/pkg/testcontainers/app"
	"github.com/ploskirev/go-rocket/platform/pkg/testcontainers/mongo"
	"github.com/ploskirev/go-rocket/platform/pkg/testcontainers/network"
	"github.com/ploskirev/go-rocket/platform/pkg/testcontainers/path"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	// Параметры для контейнеров
	inventoryAppName    = "inventory-app"
	inventoryDockerfile = "deploy/docker/inventory/Dockerfile"

	// Переменные окружения приложения
	grpcPortKey = "GRPC_PORT"

	// Значения переменных окружения
	loggerLevelValue = "debug"
	startupTimeout   = 20 * time.Minute
)

// TestEnvironment — структура для хранения ресурсов тестового окружения
type TestEnvironment struct {
	Network *network.Network
	Mongo   *mongo.Container
	App     *app.Container
}

// setupTestEnvironment — подготавливает тестовое окружение: сеть, контейнеры и возвращает структуру с ресурсами
func setupTestEnvironment(ctx context.Context) *TestEnvironment {
	logger.Info(ctx, "🚀 Подготовка тестового окружения...")

	// Шаг 1: Создаём общую Docker-сеть
	generatedNetwork, err := network.NewNetwork(ctx, projectName)
	if err != nil {
		logger.Fatal(ctx, "не удалось создать общую сеть", zap.Error(err))
	}
	logger.Info(ctx, "✅ Сеть успешно создана")

	// Получаем переменные окружения для MongoDB с проверкой на наличие
	mongoUsername := getEnvWithLogging(ctx, testcontainers.MongoUsernameKey)
	mongoPassword := getEnvWithLogging(ctx, testcontainers.MongoPasswordKey)
	mongoImageName := getEnvWithLogging(ctx, testcontainers.MongoImageNameKey)
	mongoDatabase := getEnvWithLogging(ctx, testcontainers.MongoDatabaseKey)

	// Получаем порт gRPC для waitStrategy
	grpcPort := getEnvWithLogging(ctx, grpcPortKey)

	// Шаг 2: Запускаем контейнер с MongoDB
	generatedMongo, err := mongo.NewContainer(ctx,
		mongo.WithNetworkName(generatedNetwork.Name()),
		mongo.WithContainerName(testcontainers.MongoContainerName),
		mongo.WithImageName(mongoImageName),
		mongo.WithDatabase(mongoDatabase),
		mongo.WithAuth(mongoUsername, mongoPassword),
		mongo.WithLogger(logger.Logger()),
	)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork})
		logger.Fatal(ctx, "не удалось запустить контейнер MongoDB", zap.Error(err))
	}
	logger.Info(ctx, fmt.Sprintf("✅ Контейнер MongoDB успешно запущен. %s, %+v", generatedMongo.Config().Host, generatedMongo.Config()))

	// Шаг 3: Запускаем контейнер с приложением
	projectRoot := path.GetProjectRoot()

	appEnv := map[string]string{
		// Передаём в контейнер обязательные переменные из окружения теста
		grpcPortKey:                           grpcPort,
		"LOGGER_LEVEL":                       getEnvWithLogging(ctx, "LOGGER_LEVEL"),
		"LOGGER_AS_JSON":                     getEnvWithLogging(ctx, "LOGGER_AS_JSON"),
		testcontainers.MongoImageNameKey:      mongoImageName,
		testcontainers.MongoPortKey:           getEnvWithLogging(ctx, testcontainers.MongoPortKey),
		testcontainers.MongoDatabaseKey:       mongoDatabase,
		testcontainers.MongoAuthDBKey:         getEnvWithLogging(ctx, testcontainers.MongoAuthDBKey),
		testcontainers.MongoUsernameKey:       mongoUsername,
		testcontainers.MongoPasswordKey:       mongoPassword,

		// Переопределяем хост MongoDB для подключения к контейнеру из testcontainers
		testcontainers.MongoHostKey: generatedMongo.Config().ContainerName,
	}

	// // Создаем настраиваемую стратегию ожидания с увеличенным таймаутом
	waitStrategy := wait.ForListeningPort(nat.Port(grpcPort + "/tcp")).
		WithStartupTimeout(startupTimeout).WithPollInterval(1 * time.Second)

	appContainer, err := app.NewContainer(ctx,
		app.WithName(inventoryAppName),
		app.WithPort(grpcPort),
		app.WithDockerfile(projectRoot, inventoryDockerfile),
		app.WithNetwork(generatedNetwork.Name()),
		app.WithEnv(appEnv),
		app.WithLogOutput(os.Stdout),
		app.WithStartupWait(waitStrategy),
		app.WithLogger(logger.Logger()),
	)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork, Mongo: generatedMongo})
		logger.Fatal(ctx, "не удалось запустить контейнер приложения", zap.Error(err))
	}
	logger.Info(ctx, "✅ Контейнер приложения успешно запущен")

	logger.Info(ctx, "🎉 Тестовое окружение готово")
	return &TestEnvironment{
		Network: generatedNetwork,
		Mongo:   generatedMongo,
		App:     appContainer,
	}
}

// getEnvWithLogging возвращает значение переменной окружения с логированием
func getEnvWithLogging(ctx context.Context, key string) string {
	value := os.Getenv(key)
	if value == "" {
		logger.Warn(ctx, "Переменная окружения не установлена", zap.String("key", key))
	}

	return value
}
