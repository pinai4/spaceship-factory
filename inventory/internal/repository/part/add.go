package part

import (
	"context"

	"github.com/pinai4/spaceship-factory/inventory/internal/model"
	repoConverter "github.com/pinai4/spaceship-factory/inventory/internal/repository/converter"
)

func (r *repository) Add(_ context.Context, part model.Part) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	repoPart := repoConverter.PartToRepoModel(part)

	if _, ok := r.data[repoPart.UUID]; ok {
		return model.ErrPartAlreadyExists
	}

	r.data[repoPart.UUID] = repoPart

	return nil
}
