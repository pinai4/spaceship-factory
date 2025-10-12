package converter

import (
	"github.com/pinai4/spaceship-factory/order/internal/model"
	inventoryV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/inventory/v1"
)

func PartsToModel(protoParts []*inventoryV1.Part) []model.Part {
	parts := make([]model.Part, 0, len(protoParts))
	for _, protoPart := range protoParts {
		if protoPart == nil {
			continue
		}

		part := model.Part{
			UUID:          protoPart.Uuid,
			Price:         protoPart.Price,
			StockQuantity: protoPart.StockQuantity,
		}
		parts = append(parts, part)
	}

	return parts
}
