package assembly

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/pinai4/spaceship-factory/assembly/internal/model"
)

func (s *service) AssembleShip(ctx context.Context, orderUUID, userUUID string, delay time.Duration) error {
	start := time.Now()

	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		// assembling finished
	}

	event := model.ShipAssembledEvent{
		EventUUID:    uuid.NewString(),
		OrderUUID:    orderUUID,
		UserUUID:     userUUID,
		BuildTimeSec: int64(time.Since(start).Seconds()),
	}

	if err := s.assemblyProducer.ProduceShipAssembled(ctx, event); err != nil {
		return fmt.Errorf("AssemblyService.AssembleShip produce ShipAssembledEvent error: %w", err)
	}

	return nil
}
