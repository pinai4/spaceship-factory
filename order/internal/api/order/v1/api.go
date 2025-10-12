package v1

import (
	"github.com/pinai4/spaceship-factory/order/internal/service"
)

type api struct {
	orderService service.OrderService
	// orderV1.UnimplementedHandler
}

func NewAPI(orderService service.OrderService) *api {
	return &api{orderService: orderService}
}
