package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/pinai4/spaceship-factory/iam/internal/model"
)

type UserService interface {
	Register(ctx context.Context, id uuid.UUID, userRegistrationInfo model.UserRegistrationInfo) error
	Get(ctx context.Context, id uuid.UUID) (model.User, error)
}
