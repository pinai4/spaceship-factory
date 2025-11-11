package user

import (
	"context"

	"github.com/google/uuid"

	"github.com/pinai4/spaceship-factory/iam/internal/model"
)

func (s *service) Get(ctx context.Context, id uuid.UUID) (model.User, error) {
	return s.userRepository.Get(ctx, id)
}
