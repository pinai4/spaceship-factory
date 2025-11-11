package passwordhasher_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/pinai4/spaceship-factory/iam/internal/service"
	"github.com/pinai4/spaceship-factory/iam/internal/service/passwordhasher"
)

type ServiceSuite struct {
	suite.Suite

	passwordHasher service.PasswordHasher
}

func (s *ServiceSuite) SetupTest() {
	s.passwordHasher = passwordhasher.NewPasswordHasher(4)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestPasswordHasher(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}

func (s *ServiceSuite) TestHashSuccess() {
	plainPassword := "test_password"

	hash, err := s.passwordHasher.Hash(plainPassword)
	s.Require().NoError(err)
	s.Require().NoError(s.passwordHasher.Compare(hash, plainPassword))
}

func (s *ServiceSuite) TestCompareMatch() {
	// hash of "test_password"
	passwordHash := "$2a$04$HdM0Be8ThigLd0Za/kaNq.sbr02sdbBOeyqj/v4ru0o8X2o2s08.y"

	err := s.passwordHasher.Compare(passwordHash, "test_password")
	s.Require().NoError(err)
}

func (s *ServiceSuite) TestCompareNotMatch() {
	// hash of "test_password"
	passwordHash := "$2a$04$HdM0Be8ThigLd0Za/kaNq.sbr02sdbBOeyqj/v4ru0o8X2o2s08.y"

	err := s.passwordHasher.Compare(passwordHash, "other_plain_password")
	s.Require().Error(err)
}

func (s *ServiceSuite) TestCompareHashFormatError() {
	passwordHash := "invalid_hash"

	err := s.passwordHasher.Compare(passwordHash, "test_password")
	s.Require().Error(err)
}
