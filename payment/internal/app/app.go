package app

import (
	"context"
	"errors"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/pinai4/spaceship-factory/payment/internal/config"
	"github.com/pinai4/spaceship-factory/platform/pkg/closer"
	"github.com/pinai4/spaceship-factory/platform/pkg/grpc/health"
	"github.com/pinai4/spaceship-factory/platform/pkg/logger"
	paymentV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/payment/v1"
)

type App struct {
	config            *config.Config
	closer            *closer.Closer
	logger            logger.Logger
	diContainer       *diContainer
	grpcServerOptions []grpc.ServerOption
	grpcServer        *grpc.Server
	listener          net.Listener
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
		a.initListener,
		a.initGRPCServerOptions,
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

func (a *App) initGRPCServerOptions(ctx context.Context) error {
	printerInterceptor := func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err == nil {
			if r, ok := resp.(*paymentV1.PayOrderResponse); ok && r != nil {
				a.logger.Info(ctx, "Payment was successful", logger.String("transaction_id", r.GetTransactionUuid()))
			}
		}

		return resp, err
	}

	a.grpcServerOptions = append(a.grpcServerOptions, grpc.Creds(insecure.NewCredentials()), grpc.ChainUnaryInterceptor(
		printerInterceptor,
	))

	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(a.grpcServerOptions...)
	a.closer.AddNamed("gRPC server", func(ctx context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	reflection.Register(a.grpcServer)

	// Register health service for exposing server status
	health.RegisterService(a.grpcServer)

	paymentV1.RegisterPaymentServiceServer(a.grpcServer, a.diContainer.PaymentV1API(ctx))

	return nil
}

func (a *App) runGRPCServer(ctx context.Context) error {
	a.logger.Info(ctx, fmt.Sprintf("🚀 gRPC PaymentService server listening on %s", a.config.GRPCServer.Address()))

	err := a.grpcServer.Serve(a.listener)
	if err != nil {
		return err
	}

	return nil
}
