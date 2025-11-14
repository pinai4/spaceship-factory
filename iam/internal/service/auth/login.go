package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/pinai4/spaceship-factory/iam/internal/model"
)

func (s *service) Login(ctx context.Context, login, password string) (uuid.UUID, error) {
	user, err := s.userRepository.GetByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return uuid.Nil, model.ErrInvalidCredentials
		}
		return uuid.Nil, err
	}

	if err := s.passwordHasher.Compare(user.PasswordHash, password); err != nil {
		return uuid.Nil, model.ErrInvalidCredentials
	}

	sessionID := uuid.New()
	session := model.Session{
		ID:        sessionID,
		User:      user,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(s.sessionTTL),
	}

	if err := s.sessionRepository.Create(ctx, session, s.sessionTTL); err != nil {
		return uuid.Nil, fmt.Errorf("AuthService.Login failed to create session: %w", err)
	}

	return sessionID, nil
}
