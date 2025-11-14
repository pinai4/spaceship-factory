package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/pinai4/spaceship-factory/iam/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user model.User) error
	Get(ctx context.Context, id uuid.UUID) (model.User, error)
	GetByLogin(ctx context.Context, login string) (model.User, error)
}
