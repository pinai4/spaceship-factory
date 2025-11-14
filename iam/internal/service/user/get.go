package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/pinai4/spaceship-factory/iam/internal/model"
)

func (s *service) Get(ctx context.Context, id uuid.UUID) (model.User, error) {
	user, err := s.userRepository.Get(ctx, id)
	if err != nil {
		return model.User{}, fmt.Errorf("UserService.Get failed to get user: %w", err)
	}

	return user, nil
}
