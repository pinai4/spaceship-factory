//go:build unit || !integration

package auth_test

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/pinai4/spaceship-factory/iam/internal/model"
)

func (s *ServiceSuite) TestLogin_Success() {
	now := time.Now()

	password := "test_password"

	user := model.User{
		ID: uuid.New(),
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
		PasswordHash: "test_password_hash",
		CreatedAt:    now,
		UpdatedAt:    &now,
	}

	s.userRepository.
		On("GetByLogin", s.ctx, user.Info.Login).
		Return(user, nil).
		Once()

	s.passwordHasher.
		On("Compare", user.PasswordHash, password).
		Return(nil).
		Once()

	var capturedSession model.Session
	s.sessionRepository.
		On("Create",
			s.ctx,
			mock.MatchedBy(func(s model.Session) bool {
				expected := model.Session{
					User: user,
				}
				return assert.ObjectsAreEqualValues(expected.User, s.User)
			}),
			s.sessionTTL,
		).
		Run(func(args mock.Arguments) {
			capturedSession = args.Get(1).(model.Session)
		}).
		Return(nil).
		Once()

	res, err := s.service.Login(s.ctx, user.Info.Login, password)
	s.Require().NoError(err)
	s.Require().Equal(capturedSession.ID, res)
}

func (s *ServiceSuite) TestLogin_UserRepoError() {
	repoErr := errors.New("test repo error")

	login := "test_login"
	password := "test_password"

	s.userRepository.
		On("GetByLogin", s.ctx, login).
		Return(model.User{}, repoErr).
		Once()

	res, err := s.service.Login(s.ctx, login, password)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Empty(res)
}

func (s *ServiceSuite) TestLogin_UserNotFound() {
	login := "test_login"
	password := "test_password"

	s.userRepository.
		On("GetByLogin", s.ctx, login).
		Return(model.User{}, model.ErrUserNotFound).
		Once()

	res, err := s.service.Login(s.ctx, login, password)
	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrInvalidCredentials)
	s.Require().Empty(res)
}

func (s *ServiceSuite) TestLogin_PasswordIncorrect() {
	passHasherErr := errors.New("test password invalid")

	now := time.Now()

	password := "test_password"

	user := model.User{
		ID: uuid.New(),
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
		PasswordHash: "test_password_hash",
		CreatedAt:    now,
		UpdatedAt:    &now,
	}

	s.userRepository.
		On("GetByLogin", s.ctx, user.Info.Login).
		Return(user, nil).
		Once()

	s.passwordHasher.
		On("Compare", user.PasswordHash, password).
		Return(passHasherErr).
		Once()

	res, err := s.service.Login(s.ctx, user.Info.Login, password)
	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrInvalidCredentials)
	s.Require().Empty(res)
}

func (s *ServiceSuite) TestLogin_SessionRepoError() {
	repoErr := errors.New("test repo error")

	now := time.Now()

	password := "test_password"

	user := model.User{
		ID: uuid.New(),
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
		PasswordHash: "test_password_hash",
		CreatedAt:    now,
		UpdatedAt:    &now,
	}

	s.userRepository.
		On("GetByLogin", s.ctx, user.Info.Login).
		Return(user, nil).
		Once()

	s.passwordHasher.
		On("Compare", user.PasswordHash, password).
		Return(nil).
		Once()

	s.sessionRepository.
		On("Create",
			s.ctx,
			mock.MatchedBy(func(s model.Session) bool {
				expected := model.Session{
					User: user,
				}
				return assert.ObjectsAreEqualValues(expected.User, s.User)
			}),
			s.sessionTTL,
		).
		Return(repoErr).
		Once()

	res, err := s.service.Login(s.ctx, user.Info.Login, password)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Empty(res)
}
