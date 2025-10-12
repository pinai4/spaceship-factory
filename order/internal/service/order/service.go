package order

import (
	"github.com/pinai4/spaceship-factory/order/internal/client"
	"github.com/pinai4/spaceship-factory/order/internal/repository"
	def "github.com/pinai4/spaceship-factory/order/internal/service"
)

var _ def.OrderService = (*service)(nil)

type service struct {
	orderRepository repository.OrderRepository
	paymentClient   client.PaymentClient
	inventoryClient client.InventoryClient
}

func NewService(
	orderRepository repository.OrderRepository,
	paymentClient client.PaymentClient,
	inventoryClient client.InventoryClient,
) *service {
	return &service{
		orderRepository: orderRepository,
		paymentClient:   paymentClient,
		inventoryClient: inventoryClient,
	}
}
