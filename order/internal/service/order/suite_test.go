package order_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	clientMocks "github.com/pinai4/spaceship-factory/order/internal/client/mocks"
	"github.com/pinai4/spaceship-factory/order/internal/model"
	repoMocks "github.com/pinai4/spaceship-factory/order/internal/repository/mocks"
	"github.com/pinai4/spaceship-factory/order/internal/service"
	"github.com/pinai4/spaceship-factory/order/internal/service/order"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	orderRepository *repoMocks.OrderRepository
	paymentClient   *clientMocks.PaymentClient
	inventoryClient *clientMocks.InventoryClient

	service service.OrderService
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.orderRepository = repoMocks.NewOrderRepository(s.T())
	s.paymentClient = clientMocks.NewPaymentClient(s.T())
	s.inventoryClient = clientMocks.NewInventoryClient(s.T())

	s.service = order.NewService(
		s.orderRepository,
		s.paymentClient,
		s.inventoryClient,
	)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
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
