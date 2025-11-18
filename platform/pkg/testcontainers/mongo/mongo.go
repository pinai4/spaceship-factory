package mongo

import (
	"context"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/pinai4/spaceship-factory/platform/pkg/logger"
)

const (
	port           = "27017"
	startupTimeout = 5 * time.Minute

	defaultContainerName = "mongo"

	envUsernameKey = "MONGO_INITDB_ROOT_USERNAME"
	envPasswordKey = "MONGO_INITDB_ROOT_PASSWORD" //nolint:gosec
)

type Container struct {
	container testcontainers.Container
	client    *mongo.Client
	cfg       *Config
}

func NewContainer(ctx context.Context, opts ...Option) (*Container, error) {
	cfg := buildConfig(opts...)

	container, err := startContainer(ctx, cfg)
	if err != nil {
		return nil, err
	}

	success := false
	defer func() {
		if !success {
			if err = container.Terminate(ctx); err != nil {
				cfg.Logger.Error(ctx, "failed to terminate mongo container", logger.Error(err))
			}
		}
	}()

	cfg.MappedHost, cfg.MappedPort, err = getContainerHostPort(ctx, container)
	if err != nil {
		return nil, err
	}

	uri := buildMongoURI(cfg)

	client, err := connectMongoClient(ctx, uri)
	if err != nil {
		return nil, err
	}

	cfg.Logger.Info(ctx, "Mongo container started", logger.String("uri", uri))
	success = true

	return &Container{
		container: container,
		client:    client,
		cfg:       cfg,
	}, nil
}

func (c *Container) Client() *mongo.Client {
	return c.client
}

func (c *Container) Config() *Config {
	return c.cfg
}

func (c *Container) Terminate(ctx context.Context) error {
	if err := c.client.Disconnect(ctx); err != nil {
		c.cfg.Logger.Error(ctx, "failed to disconnect mongo client", logger.Error(err))
	}

	if err := c.container.Terminate(ctx); err != nil {
		c.cfg.Logger.Error(ctx, "failed to terminate mongo container", logger.Error(err))
	}

	c.cfg.Logger.Info(ctx, "Mongo container terminated")

	return nil
}
