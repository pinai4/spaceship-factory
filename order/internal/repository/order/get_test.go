package order_test

import (
	"github.com/google/uuid"

	"github.com/pinai4/spaceship-factory/order/internal/model"
)

func (s *RepositorySuite) TestGetSuccess() {
	// Arrange
	order := buildTestOrder()
	orderUUID := order.UUID
	_ = s.repository.Create(s.ctx, order)
	_ = s.repository.Create(s.ctx, buildTestOrder())

	// Act
	res, err := s.repository.Get(s.ctx, orderUUID)

	// Assert
	s.Require().NoError(err)
	s.Require().Equal(order, res)
}

func (s *RepositorySuite) TestGetNotFound() {
	// Arrange
	orderUUID := uuid.New()

	// Act
	res, err := s.repository.Get(s.ctx, orderUUID)

	// Assert
	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrOrderNotFound)
	s.Require().Empty(res)
}
