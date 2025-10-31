package mongo

import (
	"context"

	"github.com/pinai4/spaceship-factory/platform/pkg/logger"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...logger.Field)
	Error(ctx context.Context, msg string, fields ...logger.Field)
}

type Config struct {
	// standard MongoDB port, can not be changed through Options
	Port          string
	NetworkName   string
	ContainerName string
	ImageName     string
	Database      string
	Username      string
	Password      string
	AuthDB        string
	Logger        Logger

	MappedHost string
	MappedPort string
}

func buildConfig(opts ...Option) *Config {
	cfg := &Config{
		ContainerName: defaultContainerName,
		Port:          port,
		Logger:        &logger.NoopLogger{},
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}
