package order_test

import (
	"errors"

	"github.com/google/uuid"

	"github.com/pinai4/spaceship-factory/order/internal/model"
)

func (s *ServiceSuite) TestGetSuccess() {
	order := buildTestOrder()
	orderUUID := order.UUID

	s.orderRepository.On("Get", s.ctx, orderUUID).Return(order, nil).Once()

	res, err := s.service.Get(s.ctx, orderUUID)
	s.Require().NoError(err)
	s.Require().Equal(order, res)
}

func (s *ServiceSuite) TestGetRepoError() {
	var (
		repoErr   = errors.New("test repo error")
		orderUUID = uuid.New()
	)

	s.orderRepository.On("Get", s.ctx, orderUUID).Return(model.Order{}, repoErr).Once()

	res, err := s.service.Get(s.ctx, orderUUID)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Empty(res)
}
