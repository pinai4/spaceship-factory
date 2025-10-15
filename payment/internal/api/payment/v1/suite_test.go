package v1_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	v1 "github.com/pinai4/spaceship-factory/payment/internal/api/payment/v1"
	"github.com/pinai4/spaceship-factory/payment/internal/service/mocks"
	paymentV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/payment/v1"
)

type APISuite struct {
	suite.Suite

	ctx context.Context

	paymentService *mocks.PaymentService

	api paymentV1.PaymentServiceServer
}

func (s *APISuite) SetupTest() {
	s.ctx = context.Background()

	s.paymentService = mocks.NewPaymentService(s.T())

	s.api = v1.NewAPI(
		s.paymentService,
	)
}

func (s *APISuite) TearDownTest() {
}

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(APISuite))
}
