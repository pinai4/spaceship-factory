package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/pinai4/spaceship-factory/inventory/internal/config/env"
)

type Config struct {
	Logger     LoggerConfig
	GRPCServer GRPCServerConfig
	Mongo      MongoConfig
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

	mongoCfg, err := env.NewMongoConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		Logger:     loggerCfg,
		GRPCServer: grpcServerCfg,
		Mongo:      mongoCfg,
	}, nil
}
