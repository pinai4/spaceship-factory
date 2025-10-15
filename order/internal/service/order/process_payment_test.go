package order_test

import (
	"errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/pinai4/spaceship-factory/order/internal/model"
)

func (s *ServiceSuite) TestProcessPaymentSuccess() {
	order := buildTestOrder()
	orderUUID := order.UUID
	userUUID := order.UserUUID

	paymentMethod := model.OrderPaymentMethodInvestorMoney
	tranUUID := uuid.New()

	updatedOrder := order
	updatedOrder.TransactionUUID = &tranUUID
	updatedOrder.PaymentMethod = &paymentMethod
	updatedOrder.Status = model.OrderStatusPaid

	s.orderRepository.On("Get", s.ctx, orderUUID).Return(order, nil).Once()
	s.paymentClient.On("PayOrder", s.ctx, orderUUID.String(), userUUID.String(), string(paymentMethod)).Return(tranUUID.String(), nil).Once()
	s.orderRepository.On("Update", s.ctx, orderUUID, updatedOrder).Return(nil).Once()

	res, err := s.service.ProcessPayment(s.ctx, orderUUID, paymentMethod)
	s.Require().NoError(err)
	s.Require().Equal(tranUUID, res)
}

func (s *ServiceSuite) TestProcessPaymentClientError() {
	clientErr := errors.New("test client error")

	order := buildTestOrder()
	orderUUID := order.UUID
	userUUID := order.UserUUID

	paymentMethod := model.OrderPaymentMethodInvestorMoney

	s.orderRepository.On("Get", s.ctx, orderUUID).Return(order, nil).Once()
	s.paymentClient.On("PayOrder", s.ctx, orderUUID.String(), userUUID.String(), string(paymentMethod)).Return("", clientErr).Once()
	s.orderRepository.AssertNotCalled(s.T(), "Update", s.ctx, mock.Anything, mock.Anything)

	res, err := s.service.ProcessPayment(s.ctx, orderUUID, paymentMethod)
	s.Require().Error(err)
	s.Require().ErrorIs(err, clientErr)
	s.Require().Empty(res)
}

func (s *ServiceSuite) TestProcessPaymentRepoNotFoundError() {
	repoErr := errors.New("test repo not found error")

	orderUUID := uuid.New()
	paymentMethod := model.OrderPaymentMethodInvestorMoney

	s.orderRepository.On("Get", s.ctx, orderUUID).Return(model.Order{}, repoErr).Once()
	s.paymentClient.AssertNotCalled(s.T(), "PayOrder", s.ctx, mock.Anything, mock.Anything, mock.Anything)
	s.orderRepository.AssertNotCalled(s.T(), "Update", s.ctx, mock.Anything, mock.Anything)

	res, err := s.service.ProcessPayment(s.ctx, orderUUID, paymentMethod)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Empty(res)
}

func (s *ServiceSuite) TestProcessPaymentRepoUpdateError() {
	repoErr := errors.New("test repo update error")

	order := buildTestOrder()
	orderUUID := order.UUID
	userUUID := order.UserUUID

	paymentMethod := model.OrderPaymentMethodInvestorMoney
	tranUUID := uuid.New()

	updatedOrder := order
	updatedOrder.TransactionUUID = &tranUUID
	updatedOrder.PaymentMethod = &paymentMethod
	updatedOrder.Status = model.OrderStatusPaid

	s.orderRepository.On("Get", s.ctx, orderUUID).Return(order, nil).Once()
	s.paymentClient.On("PayOrder", s.ctx, orderUUID.String(), userUUID.String(), string(paymentMethod)).Return(tranUUID.String(), nil).Once()
	s.orderRepository.On("Update", s.ctx, orderUUID, updatedOrder).Return(repoErr).Once()

	res, err := s.service.ProcessPayment(s.ctx, orderUUID, paymentMethod)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Empty(res)
}
