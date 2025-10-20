package memory

import (
	"context"

	"github.com/google/uuid"

	"github.com/pinai4/spaceship-factory/order/internal/model"
	repoConverter "github.com/pinai4/spaceship-factory/order/internal/repository/order/memory/converter"
)

func (r *repository) Update(_ context.Context, orderUUID uuid.UUID, order model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.data[orderUUID.String()]
	if !ok {
		return model.ErrOrderNotFound
	}

	repoOrder := repoConverter.OrderToRepoModel(order)
	r.data[orderUUID.String()] = repoOrder

	return nil
}
