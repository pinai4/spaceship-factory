package order_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/pinai4/spaceship-factory/order/internal/model"
	"github.com/pinai4/spaceship-factory/order/internal/repository"
	"github.com/pinai4/spaceship-factory/order/internal/repository/order"
)

type RepositorySuite struct {
	suite.Suite

	ctx context.Context

	repository repository.OrderRepository
}

func (s *RepositorySuite) SetupTest() {
	s.ctx = context.Background()

	s.repository = order.NewRepository()
}

func (s *RepositorySuite) TearDownTest() {
}

func TestRepositoryIntegration(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
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
