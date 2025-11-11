package user_test

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/pinai4/spaceship-factory/iam/internal/model"
)

func (s *ServiceSuite) TestGetSuccess() {
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

	s.userRepository.On("Get", s.ctx, user.ID).Return(user, nil).Once()

	res, err := s.service.Get(s.ctx, user.ID)
	s.Require().NoError(err)
	s.Require().Equal(user, res)
}

func (s *ServiceSuite) TestGetRepoError() {
	var (
		repoErr = errors.New("test repo error")
		id      = uuid.New()
	)

	s.userRepository.On("Get", s.ctx, id).Return(model.User{}, repoErr).Once()

	res, err := s.service.Get(s.ctx, id)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Empty(res)
}
