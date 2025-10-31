package mongo

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func startContainer(ctx context.Context, cfg *Config) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Name:     cfg.ContainerName,
		Image:    cfg.ImageName,
		Networks: []string{cfg.NetworkName},
		Env: map[string]string{
			envUsernameKey: cfg.Username,
			envPasswordKey: cfg.Password,
		},
		WaitingFor: wait.ForListeningPort(port + "/tcp").WithStartupTimeout(startupTimeout),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start mongo container: %w", err)
	}

	return container, nil
}

func getContainerHostPort(ctx context.Context, container testcontainers.Container) (string, string, error) {
	host, err := container.Host(ctx)
	if err != nil {
		return "", "", fmt.Errorf("failed to get container host: %w", err)
	}

	port, err := container.MappedPort(ctx, port+"/tcp")
	if err != nil {
		return "", "", fmt.Errorf("failed to get mapped port: %w", err)
	}

	return host, port.Port(), nil
}

func buildMongoURI(cfg *Config) string {
	return fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/%s?authSource=%s",
		cfg.Username,
		cfg.Password,
		cfg.MappedHost,
		cfg.MappedPort,
		cfg.Database,
		cfg.AuthDB,
	)
}
