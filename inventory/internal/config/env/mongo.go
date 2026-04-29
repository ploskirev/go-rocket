package env

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type postgresEnvConfig struct {
	Host         string `env:"MONGO_HOST,required"`
	Port         string `env:"MONGO_PORT,required"`
	ExternalPort string `env:"EXTERNAL_MONGO_PORT,required"`
	Database     string `env:"MONGO_DATABASE,required"`
	User         string `env:"MONGO_INITDB_ROOT_USERNAME,required"`
	Password     string `env:"MONGO_INITDB_ROOT_PASSWORD,required"`
	AuthDB       string `env:"MONGO_AUTH_DB,required"`
}

type postgresConfig struct {
	raw postgresEnvConfig
}

func NewPostgresConfig() (*postgresConfig, error) {
	var raw postgresEnvConfig
	err := env.Parse(&raw)
	if err != nil {
		return nil, err
	}

	return &postgresConfig{raw: raw}, nil
}

func (cfg *postgresConfig) URI() string {
	return fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/%s?authSource=%s",
		cfg.raw.User,
		cfg.raw.Password,
		cfg.raw.Host,
		cfg.raw.Port,
		cfg.raw.Database,
		cfg.raw.AuthDB,
	)
}

func (cfg *postgresConfig) DatabaseName() string {
	return cfg.raw.Database
}
