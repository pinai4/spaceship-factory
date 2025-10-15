package v1_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	orderV1API "github.com/pinai4/spaceship-factory/order/internal/api/order/v1"
	"github.com/pinai4/spaceship-factory/order/internal/model"
	"github.com/pinai4/spaceship-factory/order/internal/service/mocks"
	orderV1 "github.com/pinai4/spaceship-factory/shared/pkg/openapi/order/v1"
)

type APISuite struct {
	suite.Suite

	ctx context.Context

	orderService *mocks.OrderService

	api orderV1.Handler
}

func (s *APISuite) SetupTest() {
	s.ctx = context.Background()

	s.orderService = mocks.NewOrderService(s.T())

	s.api = orderV1API.NewAPI(
		s.orderService,
	)
}

func (s *APISuite) TearDownTest() {
}

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(APISuite))
}

func buildTestOrder() model.Order {
	return model.Order{
		UUID:       uuid.New(),
		UserUUID:   uuid.New(),
		PartUUIDs:  []uuid.UUID{uuid.New(), uuid.New()},
		TotalPrice: 9.99,
		Status:     model.OrderStatusPendingPayment,
	}
}
