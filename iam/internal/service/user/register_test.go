//go:build unit || !integration

package user_test

import (
	"errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/pinai4/spaceship-factory/iam/internal/model"
)

func (s *ServiceSuite) TestRegisterSuccess() {
	id := uuid.New()

	registrationInfo := model.UserRegistrationInfo{
		Info: model.UserInfo{
			Login: "test_login",
			Email: "test_email",
			NotificationMethods: []model.NotificationMethod{
				{
					Provider: "test_provider",
					Target:   "test_target",
				},
			},
		},
		Password: "test_password",
	}

	passwordHash := "test_password_hash"

	s.passwordHasher.On("Hash", registrationInfo.Password).Return(passwordHash, nil).Once()
	s.userRepository.On("Create", s.ctx, mock.MatchedBy(func(u model.User) bool {
		return u.ID == id &&
			u.Info.Login == registrationInfo.Info.Login &&
			u.Info.Email == registrationInfo.Info.Email &&
			u.Info.NotificationMethods[0] == registrationInfo.Info.NotificationMethods[0] &&
			u.PasswordHash == passwordHash &&
			u.UpdatedAt == nil
	})).Return(nil).Once()

	err := s.service.Register(s.ctx, id, registrationInfo)
	s.Require().NoError(err)
}

func (s *ServiceSuite) TestRegisterHasherError() {
	hasherErr := errors.New("test passwordHasher error")

	id := uuid.New()

	registrationInfo := model.UserRegistrationInfo{
		Info: model.UserInfo{
			Login: "test_login",
			Email: "test_email",
			NotificationMethods: []model.NotificationMethod{
				{
					Provider: "test_provider",
					Target:   "test_target",
				},
			},
		},
		Password: "test_password",
	}

	s.passwordHasher.On("Hash", registrationInfo.Password).Return("", hasherErr).Once()

	err := s.service.Register(s.ctx, id, registrationInfo)
	s.Require().Error(err)
	s.Require().ErrorIs(err, hasherErr)
}

func (s *ServiceSuite) TestRegisterRepoError() {
	repoErr := errors.New("test repo error")

	id := uuid.New()

	registrationInfo := model.UserRegistrationInfo{
		Info: model.UserInfo{
			Login: "test_login",
			Email: "test_email",
			NotificationMethods: []model.NotificationMethod{
				{
					Provider: "test_provider",
					Target:   "test_target",
				},
			},
		},
		Password: "test_password",
	}

	passwordHash := "test_password_hash"

	s.passwordHasher.On("Hash", registrationInfo.Password).Return(passwordHash, nil).Once()
	s.userRepository.On("Create", s.ctx, mock.MatchedBy(func(u model.User) bool {
		return u.ID == id &&
			u.Info.Login == registrationInfo.Info.Login &&
			u.Info.Email == registrationInfo.Info.Email &&
			u.Info.NotificationMethods[0] == registrationInfo.Info.NotificationMethods[0] &&
			u.PasswordHash == passwordHash &&
			u.UpdatedAt == nil
	})).Return(repoErr).Once()

	err := s.service.Register(s.ctx, id, registrationInfo)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
}
