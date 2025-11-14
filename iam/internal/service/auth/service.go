package auth

import (
	"time"

	"github.com/pinai4/spaceship-factory/iam/internal/repository"
	def "github.com/pinai4/spaceship-factory/iam/internal/service"
)

var _ def.AuthService = (*service)(nil)

type service struct {
	userRepository    repository.UserRepository
	sessionRepository repository.SessionRepository
	passwordHasher    def.PasswordHasher
	sessionTTL        time.Duration
}

func NewService(
	userRepository repository.UserRepository,
	sessionRepository repository.SessionRepository,
	passwordHasher def.PasswordHasher,
	sessionTTL time.Duration,
) *service {
	return &service{
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
		passwordHasher:    passwordHasher,
		sessionTTL:        sessionTTL,
	}
}
