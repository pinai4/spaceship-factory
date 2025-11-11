//go:build unit || !integration

package user_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	repoMocks "github.com/pinai4/spaceship-factory/iam/internal/repository/mocks"
	"github.com/pinai4/spaceship-factory/iam/internal/service"
	serviceMocks "github.com/pinai4/spaceship-factory/iam/internal/service/mocks"
	"github.com/pinai4/spaceship-factory/iam/internal/service/user"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	userRepository *repoMocks.UserRepository
	passwordHasher *serviceMocks.PasswordHasher

	service service.UserService
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.userRepository = repoMocks.NewUserRepository(s.T())
	s.passwordHasher = serviceMocks.NewPasswordHasher(s.T())

	s.service = user.NewService(s.userRepository, s.passwordHasher)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestService(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
