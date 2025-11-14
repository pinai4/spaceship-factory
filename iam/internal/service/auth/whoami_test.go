//go:build unit || !integration

package auth_test

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/pinai4/spaceship-factory/iam/internal/model"
)

func (s *ServiceSuite) TestWhoami_Success() {
	now := time.Now()

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

	session := model.Session{
		ID:        uuid.New(),
		User:      user,
		CreatedAt: now,
		UpdatedAt: &now,
		ExpiresAt: now,
	}

	s.sessionRepository.On("Get", s.ctx, session.ID).Return(session, nil).Once()

	res, err := s.service.Whoami(s.ctx, session.ID)
	s.Require().NoError(err)
	s.Require().Equal(session, res)
}

func (s *ServiceSuite) TestWhoami_RepoError() {
	var (
		repoErr = errors.New("test repo error")
		id      = uuid.New()
	)

	s.sessionRepository.On("Get", s.ctx, id).Return(model.Session{}, repoErr).Once()

	res, err := s.service.Whoami(s.ctx, id)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Empty(res)
}
