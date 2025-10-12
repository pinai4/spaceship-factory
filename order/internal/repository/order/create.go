package order

import (
	"context"

	"github.com/pinai4/spaceship-factory/order/internal/model"
	repoConverter "github.com/pinai4/spaceship-factory/order/internal/repository/converter"
)

func (r *repository) Create(_ context.Context, order model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	repoOrder := repoConverter.OrderToRepoModel(order)

	if _, ok := r.data[repoOrder.UUID.String()]; ok {
		return model.ErrOrderAlreadyExists
	}

	r.data[order.UUID.String()] = repoOrder

	return nil
}
