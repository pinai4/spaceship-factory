//go:build unit || !integration

package order_test

import (
	"errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/pinai4/spaceship-factory/order/internal/model"
)

func (s *ServiceSuite) TestAssembleSuccess() {
	order := buildTestOrder()
	orderUUID := order.UUID

	updatedOrder := order
	updatedOrder.Status = model.OrderStatusAssembled

	s.orderRepository.On("Get", s.ctx, orderUUID).Return(order, nil).Once()
	s.orderRepository.On("Update", s.ctx, orderUUID, updatedOrder).Return(nil).Once()

	err := s.service.Assemble(s.ctx, orderUUID)
	s.Require().NoError(err)
}

func (s *ServiceSuite) TestAssembleRepoNotFoundError() {
	repoErr := errors.New("test repo not found error")

	orderUUID := uuid.New()

	s.orderRepository.On("Get", s.ctx, orderUUID).Return(model.Order{}, repoErr).Once()
	s.orderRepository.AssertNotCalled(s.T(), "Update", s.ctx, mock.Anything, mock.Anything)

	err := s.service.Assemble(s.ctx, orderUUID)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
}

func (s *ServiceSuite) TestAssembleRepoUpdateError() {
	repoErr := errors.New("test repo update error")

	order := buildTestOrder()
	orderUUID := order.UUID

	updatedOrder := order
	updatedOrder.Status = model.OrderStatusAssembled

	s.orderRepository.On("Get", s.ctx, orderUUID).Return(order, nil).Once()
	s.orderRepository.On("Update", s.ctx, orderUUID, updatedOrder).Return(repoErr).Once()

	err := s.service.Assemble(s.ctx, orderUUID)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
}
