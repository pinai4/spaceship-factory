package service

import (
	"context"
	"time"
)

type AssemblyService interface {
	AssembleShip(ctx context.Context, orderUUID, userUUID string, delay time.Duration) error
}
