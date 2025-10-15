package order_test

import "github.com/pinai4/spaceship-factory/order/internal/model"

func (s *RepositorySuite) TestCreateSuccess() {
	// Arrange
	order := buildTestOrder()
	orderUUID := order.UUID

	// Act
	err := s.repository.Create(s.ctx, order)

	// Assert
	s.Require().NoError(err)

	res, err := s.repository.Get(s.ctx, orderUUID)
	s.Require().NoError(err)
	s.Require().Equal(order, res)
}

func (s *RepositorySuite) TestCreateOrderAlreadyExistsError() {
	// Arrange
	order := buildTestOrder()
	_ = s.repository.Create(s.ctx, order)

	// Act
	err := s.repository.Create(s.ctx, order)

	// Assert
	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrOrderAlreadyExists)
}
