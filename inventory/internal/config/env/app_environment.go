package env

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

const (
	appEnvDev  string = "dev"
	appEnvTest string = "test"
	appEnvProd string = "prod"
)

type appEnvironmentEnvConfig struct {
	AppEnv string `env:"APP_ENV,required"`
}

type appEnvironmentConfig struct {
	raw appEnvironmentEnvConfig
}

func NewAppEnvironmentConfig() (*appEnvironmentConfig, error) {
	var raw appEnvironmentEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	if raw.AppEnv != appEnvDev && raw.AppEnv != appEnvTest && raw.AppEnv != appEnvProd {
		return nil, fmt.Errorf("invalid app environment: %s", raw.AppEnv)
	}

	return &appEnvironmentConfig{raw: raw}, nil
}

func (cfg *appEnvironmentConfig) IsDev() bool {
	return cfg.raw.AppEnv == appEnvDev
}

func (cfg *appEnvironmentConfig) IsTest() bool {
	return cfg.raw.AppEnv == appEnvTest
}

func (cfg *appEnvironmentConfig) IsProd() bool {
	return cfg.raw.AppEnv == appEnvProd
}
