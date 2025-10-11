package part

import (
	"context"

	"github.com/pinai4/spaceship-factory/inventory/internal/model"
)

func (s *service) Get(ctx context.Context, uuid string) (model.Part, error) {
	return s.partRepository.Get(ctx, uuid)
}
