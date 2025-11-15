package seeds

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/pinai4/spaceship-factory/iam/internal/model"
	"github.com/pinai4/spaceship-factory/iam/internal/repository"
)

var (
	userID    = uuid.MustParse("00000000-0000-0000-0000-111111111111")
	userLogin = "tester"

	// hash for password "secret"
	userHashPassword = "$2a$04$BEqr.hPxbQXgBBefNy53quzSlscu8aYRvcks.5OHiyAfIDOaRNtgu" //nolint:gosec

	notificationMethodProvider = "telegram"

	sessionID  = uuid.MustParse("00000000-0000-0000-0000-222222222222")
	sessionTTL = 365 * 24 * time.Hour
)

type seeder struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
}

func New(userRepo repository.UserRepository, sessionRepo repository.SessionRepository) *seeder {
	return &seeder{userRepo: userRepo, sessionRepo: sessionRepo}
}

func (s *seeder) getUser() model.User {
	return model.User{
		ID: userID,
		Info: model.UserInfo{
			Login: userLogin,
			Email: "tester@tester.com",
			NotificationMethods: []model.NotificationMethod{
				{
					Provider: notificationMethodProvider,
					Target:   "111111111",
				},
			},
		},
		PasswordHash: userHashPassword,
		CreatedAt:    time.Now(),
	}
}

func (s *seeder) seedUser(ctx context.Context, user model.User) error {
	_, err := s.userRepo.Get(ctx, userID)
	if err == nil {
		return nil
	}
	if !errors.Is(err, model.ErrUserNotFound) {
		return err
	}

	return s.userRepo.Create(ctx, user)
}

func (s *seeder) seedSession(ctx context.Context, user model.User) error {
	_, err := s.sessionRepo.Get(ctx, sessionID)
	if err == nil {
		return nil
	}
	if !errors.Is(err, model.ErrSessionNotFound) {
		return err
	}

	now := time.Now()
	session := model.Session{
		ID:        sessionID,
		User:      user,
		CreatedAt: now,
		ExpiresAt: now.Add(sessionTTL),
	}

	return s.sessionRepo.Create(ctx, session, sessionTTL)
}

func (s *seeder) Seed(ctx context.Context) error {
	user := s.getUser()
	if err := s.seedUser(ctx, user); err != nil {
		return err
	}
	if err := s.seedSession(ctx, user); err != nil {
		return err
	}

	return nil
}
