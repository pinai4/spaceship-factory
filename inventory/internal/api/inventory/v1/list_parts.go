package v1

import (
	"context"

	"github.com/pinai4/spaceship-factory/inventory/internal/converter"
	inventoryV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	parts, err := a.partService.List(ctx, converter.PartsFilterToModel(req.GetFilter()))
	if err != nil {
		return nil, err
	}

	return &inventoryV1.ListPartsResponse{Parts: converter.PartsToProto(parts)}, nil
}
