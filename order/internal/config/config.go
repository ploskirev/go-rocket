package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/ploskirev/go-rocket/order/internal/config/env"
)

var appConfig *config

type config struct {
	Logger    LoggerConfig
	OrderHTTP OrderHTTPConfig
	Postgres  PostgresConfig
	Inventory InventoryGRPCConfig
	Payment   PaymentGRPCConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	orderHTTPCfg, err := env.NewOrderGRPCConfig()
	if err != nil {
		return err
	}

	postgresCFG, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	inventoryGRPCCfg, err := env.NewInventoryGRPCConfig()
	if err != nil {
		return err
	}

	paymentGRPCCfg, err := env.NewPaymentGRPCConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:    loggerCfg,
		OrderHTTP: orderHTTPCfg,
		Postgres:  postgresCFG,
		Inventory: inventoryGRPCCfg,
		Payment:   paymentGRPCCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
