package memory_test

import (
	"github.com/google/uuid"

	"github.com/pinai4/spaceship-factory/order/internal/model"
)

func (s *RepositorySuite) TestUpdateSuccess() {
	// Arrange
	order := buildTestOrder()
	orderUUID := order.UUID
	_ = s.repository.Create(s.ctx, order)

	tranUUID := uuid.New()
	pm := model.OrderPaymentMethodCreditCard

	updatedOrder := order
	updatedOrder.UserUUID = uuid.New()
	updatedOrder.PartUUIDs = []uuid.UUID{uuid.New()}
	updatedOrder.TotalPrice = 100
	updatedOrder.TransactionUUID = &tranUUID
	updatedOrder.PaymentMethod = &pm
	updatedOrder.Status = model.OrderStatusPaid

	// Act
	err := s.repository.Update(s.ctx, orderUUID, updatedOrder)

	// Assert
	s.Require().NoError(err)

	res, err := s.repository.Get(s.ctx, orderUUID)
	s.Require().NoError(err)
	s.Require().Equal(updatedOrder, res)
}

func (s *RepositorySuite) TestUpdateNotFound() {
	// Arrange
	orderUUID := uuid.New()

	// Act
	err := s.repository.Update(s.ctx, orderUUID, model.Order{})

	// Assert
	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrOrderNotFound)
}
