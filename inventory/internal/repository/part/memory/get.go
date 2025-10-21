package memory

import (
	"context"

	"github.com/pinai4/spaceship-factory/inventory/internal/model"
	repoConverter "github.com/pinai4/spaceship-factory/inventory/internal/repository/part/memory/converter"
)

func (r *repository) Get(_ context.Context, uuid string) (model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	repoPart, ok := r.data[uuid]
	if !ok {
		return model.Part{}, model.ErrPartNotFound
	}

	return repoConverter.PartToModel(repoPart), nil
}
