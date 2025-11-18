package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/pinai4/spaceship-factory/platform/pkg/logger"
)

const (
	defaultPort          = "50051"
	defaultContainerName = "app"

	startupTimeout = 5 * time.Minute
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...logger.Field)
	Error(ctx context.Context, msg string, fields ...logger.Field)
}

type Config struct {
	ContainerName string
	DockerfileDir string
	Dockerfile    string
	Port          string
	Env           map[string]string
	Networks      []string
	LogOutput     io.Writer
	StartupWait   wait.Strategy
	Logger        Logger

	MappedHost string
	MappedPort string
}

type Container struct {
	container testcontainers.Container
	cfg       *Config
}

func NewContainer(ctx context.Context, opts ...Option) (*Container, error) {
	cfg := &Config{
		ContainerName: defaultContainerName,
		Port:          defaultPort,
		Dockerfile:    "Dockerfile",
		DockerfileDir: ".",
		LogOutput:     io.Discard,
		StartupWait:   wait.ForListeningPort(defaultPort + "/tcp").WithStartupTimeout(startupTimeout),
		Env:           make(map[string]string),
		Logger:        &logger.NoopLogger{},
	}
	for _, opt := range opts {
		opt(cfg)
	}

	req := testcontainers.ContainerRequest{
		Name: cfg.ContainerName,
		FromDockerfile: testcontainers.FromDockerfile{
			Context:        cfg.DockerfileDir,
			Dockerfile:     cfg.Dockerfile,
			BuildLogWriter: cfg.LogOutput,
		},
		Networks:     cfg.Networks,
		Env:          cfg.Env,
		WaitingFor:   cfg.StartupWait,
		ExposedPorts: []string{cfg.Port + "/tcp"},
	}

	genericContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start app genericContainer: %w", err)
	}

	mappedPort, err := genericContainer.MappedPort(ctx, nat.Port(cfg.Port+"/tcp"))
	if err != nil {
		return nil, fmt.Errorf("failed to get mapped externalPort: %w", err)
	}
	cfg.MappedPort = mappedPort.Port()

	cfg.MappedHost, err = genericContainer.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get genericContainer externalHost: %w", err)
	}

	go streamContainerLogs(ctx, cfg.Logger, genericContainer, cfg.LogOutput)

	cfg.Logger.Info(ctx, fmt.Sprintf("%s container started", cfg.ContainerName), logger.String("uri:", net.JoinHostPort(cfg.MappedHost, cfg.MappedPort)))

	return &Container{
		container: genericContainer,
		cfg:       cfg,
	}, nil
}

func (a *Container) Address() string {
	return net.JoinHostPort(a.cfg.MappedHost, a.cfg.MappedPort)
}

func (a *Container) Config() *Config {
	return a.cfg
}

func (a *Container) Terminate(ctx context.Context) error {
	return a.container.Terminate(ctx)
}

func streamContainerLogs(ctx context.Context, log Logger, container testcontainers.Container, out io.Writer) {
	logs, err := container.Logs(ctx)
	if err != nil {
		log.Error(ctx, "failed to get container logs", logger.Error(err))
		return
	}
	defer func() {
		err = logs.Close()
		if err != nil {
			log.Error(ctx, "failed to close container logs", logger.Error(err))
		}
	}()

	_, err = io.Copy(out, logs)
	if err != nil && !errors.Is(err, io.EOF) {
		log.Error(ctx, "error copying container logs", logger.Error(err))
	}
}
