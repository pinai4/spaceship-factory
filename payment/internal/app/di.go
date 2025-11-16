package app

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	paymentV1API "github.com/pinai4/spaceship-factory/payment/internal/api/payment/v1"
	"github.com/pinai4/spaceship-factory/payment/internal/config"
	"github.com/pinai4/spaceship-factory/payment/internal/service"
	paymentService "github.com/pinai4/spaceship-factory/payment/internal/service/payment"
	"github.com/pinai4/spaceship-factory/platform/pkg/closer"
	authV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/auth/v1"
	paymentV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	config *config.Config
	closer *closer.Closer

	paymentV1API paymentV1.PaymentServiceServer

	paymentService service.PaymentService

	authV1ClientGRPC authV1.AuthServiceClient
}

func NewDiContainer(config *config.Config, closer *closer.Closer) *diContainer {
	return &diContainer{config: config, closer: closer}
}

func (d *diContainer) PaymentV1API(ctx context.Context) paymentV1.PaymentServiceServer {
	if d.paymentV1API == nil {
		d.paymentV1API = paymentV1API.NewAPI(d.PaymentService(ctx))
	}

	return d.paymentV1API
}

func (d *diContainer) PaymentService(_ context.Context) service.PaymentService {
	if d.paymentService == nil {
		d.paymentService = paymentService.NewService()
	}

	return d.paymentService
}

func (d *diContainer) AuthV1ClientGRPC() authV1.AuthServiceClient {
	if d.authV1ClientGRPC == nil {
		conn, err := grpc.NewClient(
			d.Config().IAMGRPCClient.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to connect to IAM Service (GRPC): %s\n", err.Error()))
		}
		d.closer.AddNamed("IAMClient connection", func(ctx context.Context) error {
			return conn.Close()
		})

		d.authV1ClientGRPC = authV1.NewAuthServiceClient(conn)
	}

	return d.authV1ClientGRPC
}

func (d *diContainer) Config() *config.Config {
	return d.config
}
