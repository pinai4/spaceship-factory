package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/pinai4/spaceship-factory/iam/internal/config/env"
)

type Config struct {
	Logger       LoggerConfig
	GRPCServer   GRPCServerConfig
	Postgres     PostgresConfig
	DBMigrations DBMigrationsConfig
}

func Load(path ...string) (*Config, error) {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return nil, err
	}

	grpcServerCfg, err := env.NewGRPCServerConfig()
	if err != nil {
		return nil, err
	}

	postgresCfg, err := env.NewPostgresConfig()
	if err != nil {
		return nil, err
	}

	dbMigrationsCfg, err := env.NewDBMigrationsConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		Logger:       loggerCfg,
		GRPCServer:   grpcServerCfg,
		Postgres:     postgresCfg,
		DBMigrations: dbMigrationsCfg,
	}, nil
}
