package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/pinai4/spaceship-factory/iam/internal/model"
)

type AuthService interface {
	Login(ctx context.Context, login, password string) (uuid.UUID, error)
	Whoami(ctx context.Context, sessionID uuid.UUID) (model.Session, error)
}
