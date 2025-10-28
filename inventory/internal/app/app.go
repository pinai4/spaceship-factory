package app

import (
	"context"
	"errors"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/pinai4/spaceship-factory/inventory/internal/config"
	"github.com/pinai4/spaceship-factory/inventory/internal/repository/seeds"
	"github.com/pinai4/spaceship-factory/platform/pkg/closer"
	"github.com/pinai4/spaceship-factory/platform/pkg/grpc/health"
	"github.com/pinai4/spaceship-factory/platform/pkg/logger"
	inventoryV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/inventory/v1"
)

type App struct {
	config      *config.Config
	closer      *closer.Closer
	logger      logger.Logger
	diContainer *diContainer
	grpcServer  *grpc.Server
	listener    net.Listener
}

func New(ctx context.Context, config *config.Config, closer *closer.Closer, logger logger.Logger) (*App, error) {
	a := &App{config: config, closer: closer, logger: logger}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	return a.runGRPCServer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initDBSeeds,
		a.initListener,
		a.initGRPCServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDiContainer(a.config, a.closer)
	return nil
}

func (a *App) initDBSeeds(ctx context.Context) error {
	seeder := seeds.New(a.diContainer.PartRepository(ctx))
	if err := seeder.Seed(ctx); err != nil {
		return err
	}

	return nil
}

func (a *App) initListener(_ context.Context) error {
	listener, err := net.Listen("tcp", a.config.GRPCServer.Address())
	if err != nil {
		return err
	}
	a.closer.AddNamed("TCP listener", func(ctx context.Context) error {
		lerr := listener.Close()
		if lerr != nil && !errors.Is(lerr, net.ErrClosed) {
			return lerr
		}

		return nil
	})

	a.listener = listener

	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	a.closer.AddNamed("gRPC server", func(ctx context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	reflection.Register(a.grpcServer)

	// Register health service for exposing server status
	health.RegisterService(a.grpcServer)

	inventoryV1.RegisterInventoryServiceServer(a.grpcServer, a.diContainer.InventoryV1API(ctx))

	return nil
}

func (a *App) runGRPCServer(ctx context.Context) error {
	a.logger.Info(ctx, fmt.Sprintf("🚀 gRPC InventoryService server listening on %s", a.config.GRPCServer.Address()))

	err := a.grpcServer.Serve(a.listener)
	if err != nil {
		return err
	}

	return nil
}
