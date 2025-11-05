package service

import (
	"context"

	"github.com/pinai4/spaceship-factory/assembly/internal/model"
)

type AssemblyProducer interface {
	ProduceShipAssembled(ctx context.Context, event model.ShipAssembledEvent) error
}
