package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type iamGRPCClientEnvConfig struct {
	Host string `env:"IAM_GRPC_HOST,required"`
	Port string `env:"IAM_GRPC_PORT,required"`
}

type iamGRPCClientConfig struct {
	raw iamGRPCClientEnvConfig
}

func NewIAMGRPCClientConfig() (*iamGRPCClientConfig, error) {
	var raw iamGRPCClientEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &iamGRPCClientConfig{raw: raw}, nil
}

func (cfg *iamGRPCClientConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}
