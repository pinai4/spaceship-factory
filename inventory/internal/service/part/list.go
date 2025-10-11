package part

import (
	"context"

	"github.com/pinai4/spaceship-factory/inventory/internal/model"
)

func (s *service) List(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	return s.partRepository.List(ctx, filter)
}
