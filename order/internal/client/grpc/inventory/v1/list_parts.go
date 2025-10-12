package v1

import (
	"context"
	"fmt"

	"github.com/pinai4/spaceship-factory/order/internal/client/converter"
	"github.com/pinai4/spaceship-factory/order/internal/model"
	inventoryV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/inventory/v1"
)

func (c *client) ListParts(ctx context.Context, partUUIDs []string) ([]model.Part, error) {
	resp, err := c.generatedClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{
			Uuids: partUUIDs,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("client API call error: %w", err)
	}

	return converter.PartsToModel(resp.GetParts()), nil
}
