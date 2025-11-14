package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/pinai4/spaceship-factory/iam/internal/model"
)

func (s *service) Whoami(ctx context.Context, sessionID uuid.UUID) (model.Session, error) {
	session, err := s.sessionRepository.Get(ctx, sessionID)
	if err != nil {
		return model.Session{}, fmt.Errorf("AuthService.Whoami failed to get session: %w", err)
	}

	return session, nil
}
