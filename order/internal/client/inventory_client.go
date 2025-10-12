package client

import (
	"context"

	"github.com/pinai4/spaceship-factory/order/internal/model"
)

type InventoryClient interface {
	ListParts(ctx context.Context, partUUIDs []string) ([]model.Part, error)
}
