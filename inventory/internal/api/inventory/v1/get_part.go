package v1

import (
	"context"

	"github.com/go-faster/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pinai4/spaceship-factory/inventory/internal/converter"
	"github.com/pinai4/spaceship-factory/inventory/internal/model"
	inventoryV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/inventory/v1"
)

func (a *api) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	part, err := a.partService.Get(ctx, req.GetUuid())
	if err != nil {
		if errors.Is(err, model.ErrPartNotFound) {
			return nil, status.Errorf(codes.NotFound, "part with UUID %s not found", req.GetUuid())
		}
		return nil, err
	}

	return &inventoryV1.GetPartResponse{Part: converter.PartToProto(part)}, nil
}
