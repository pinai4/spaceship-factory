package payment_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/pinai4/spaceship-factory/payment/internal/service"
	"github.com/pinai4/spaceship-factory/payment/internal/service/payment"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	service service.PaymentService
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.service = payment.NewService()
}

func (s *ServiceSuite) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
