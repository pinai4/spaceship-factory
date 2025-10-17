package service

import (
	"context"

	"github.com/pinai4/spaceship-factory/inventory/internal/model"
)

type PartService interface {
	Get(ctx context.Context, uuid string) (model.Part, error)
	List(ctx context.Context, filter model.PartsFilter) ([]model.Part, error)
}
