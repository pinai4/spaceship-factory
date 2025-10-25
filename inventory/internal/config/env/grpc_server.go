package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type grpcServerEnvConfig struct {
	Host string `env:"GRPC_HOST,required"`
	Port string `env:"GRPC_PORT,required"`
}

type grpcServerConfig struct {
	raw grpcServerEnvConfig
}

func NewGRPCServerConfig() (*grpcServerConfig, error) {
	var raw grpcServerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &grpcServerConfig{raw: raw}, nil
}

func (cfg *grpcServerConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}
