package repository

import (
	"context"

	"github.com/pinai4/spaceship-factory/inventory/internal/model"
)

type PartRepository interface {
	Get(ctx context.Context, uuid string) (model.Part, error)
	List(ctx context.Context, filter model.PartsFilter) ([]model.Part, error)
	Add(ctx context.Context, part model.Part) error
}
