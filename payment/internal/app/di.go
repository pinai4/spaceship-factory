package app

import (
	"context"

	paymentV1API "github.com/pinai4/spaceship-factory/payment/internal/api/payment/v1"
	"github.com/pinai4/spaceship-factory/payment/internal/config"
	"github.com/pinai4/spaceship-factory/payment/internal/service"
	paymentService "github.com/pinai4/spaceship-factory/payment/internal/service/payment"
	"github.com/pinai4/spaceship-factory/platform/pkg/closer"
	paymentV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	config *config.Config
	closer *closer.Closer

	paymentV1API paymentV1.PaymentServiceServer

	paymentService service.PaymentService
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

func (d *diContainer) Config(_ context.Context) *config.Config {
	return d.config
}
