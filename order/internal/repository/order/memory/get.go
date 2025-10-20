package memory

import (
	"context"

	"github.com/google/uuid"

	"github.com/pinai4/spaceship-factory/order/internal/model"
	repoConverter "github.com/pinai4/spaceship-factory/order/internal/repository/order/memory/converter"
)

func (r *repository) Get(_ context.Context, orderUUID uuid.UUID) (model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	repoOrder, ok := r.data[orderUUID.String()]
	if !ok {
		return model.Order{}, model.ErrOrderNotFound
	}

	return repoConverter.OrderToModel(repoOrder), nil
}
