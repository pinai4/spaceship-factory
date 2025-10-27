package env

import (
	"net"
	"time"

	"github.com/caarlos0/env/v11"
)

type httpServerEnvConfig struct {
	Host        string        `env:"HTTP_HOST,required"`
	Port        string        `env:"HTTP_PORT,required"`
	ReadTimeout time.Duration `env:"HTTP_READ_TIMEOUT,required"`
}

type httpServerConfig struct {
	raw httpServerEnvConfig
}

func NewHTTPServerConfig() (*httpServerConfig, error) {
	var raw httpServerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &httpServerConfig{raw: raw}, nil
}

func (cfg *httpServerConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}

func (cfg *httpServerConfig) ReadTimeout() time.Duration {
	return cfg.raw.ReadTimeout
}
