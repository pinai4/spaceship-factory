package env

import (
	"github.com/caarlos0/env/v11"
)

type dbMigrationsEnvConfig struct {
	Path string `env:"MIGRATION_DIRECTORY,required"`
}

type dbMigrationsConfig struct {
	raw dbMigrationsEnvConfig
}

func NewDBMigrationsConfig() (*dbMigrationsConfig, error) {
	var raw dbMigrationsEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &dbMigrationsConfig{raw: raw}, nil
}

func (cfg *dbMigrationsConfig) Path() string {
	return cfg.raw.Path
}
