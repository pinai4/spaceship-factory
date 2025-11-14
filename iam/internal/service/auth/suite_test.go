//go:build unit || !integration

package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	repoMocks "github.com/pinai4/spaceship-factory/iam/internal/repository/mocks"
	"github.com/pinai4/spaceship-factory/iam/internal/service"
	"github.com/pinai4/spaceship-factory/iam/internal/service/auth"
	serviceMocks "github.com/pinai4/spaceship-factory/iam/internal/service/mocks"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	userRepository    *repoMocks.UserRepository
	sessionRepository *repoMocks.SessionRepository
	passwordHasher    *serviceMocks.PasswordHasher

	sessionTTL time.Duration

	service service.AuthService
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.userRepository = repoMocks.NewUserRepository(s.T())
	s.sessionRepository = repoMocks.NewSessionRepository(s.T())
	s.passwordHasher = serviceMocks.NewPasswordHasher(s.T())

	s.sessionTTL = time.Second

	s.service = auth.NewService(s.userRepository, s.sessionRepository, s.passwordHasher, s.sessionTTL)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestAuthService(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
