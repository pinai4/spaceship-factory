package user

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/pinai4/spaceship-factory/iam/internal/model"
)

func (s *service) Register(ctx context.Context, id uuid.UUID, userRegistrationInfo model.UserRegistrationInfo) error {
	passwordHash, err := s.passwordHasher.Hash(userRegistrationInfo.Password)
	if err != nil {
		return fmt.Errorf("UserService.Register failed to hash password: %w", err)
	}

	user := model.User{
		ID:           id,
		Info:         userRegistrationInfo.Info,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
	}

	if err := s.userRepository.Create(ctx, user); err != nil {
		return fmt.Errorf("UserService.Register create user error: %w", err)
	}

	return nil
}
