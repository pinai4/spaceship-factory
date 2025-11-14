package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/pinai4/spaceship-factory/iam/internal/model"
)

type SessionRepository interface {
	Create(ctx context.Context, session model.Session, ttl time.Duration) error
	Get(ctx context.Context, id uuid.UUID) (model.Session, error)
}
