package order_test

import (
	"errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/pinai4/spaceship-factory/order/internal/model"
)

func (s *ServiceSuite) TestCreateSuccess() {
	orderUUID := uuid.New()
	createOrder := model.CreateOrder{
		UserUUID:  uuid.New(),
		PartUUIDs: []uuid.UUID{uuid.New(), uuid.New()},
	}

	parts := []model.Part{
		{
			UUID:          createOrder.PartUUIDs[0].String(),
			Price:         0.01,
			StockQuantity: 1,
		},
		{
			UUID:          createOrder.PartUUIDs[1].String(),
			Price:         0.01,
			StockQuantity: 1,
		},
	}

	order := model.Order{
		UUID:            orderUUID,
		UserUUID:        createOrder.UserUUID,
		PartUUIDs:       createOrder.PartUUIDs,
		TotalPrice:      0.02,
		TransactionUUID: nil,
		PaymentMethod:   nil,
		Status:          model.OrderStatusPendingPayment,
	}

	s.inventoryClient.On("ListParts", s.ctx, []string{createOrder.PartUUIDs[0].String(), createOrder.PartUUIDs[1].String()}).Return(parts, nil).Once()
	s.orderRepository.On("Create", s.ctx, order).Return(nil).Once()

	res, err := s.service.Create(s.ctx, orderUUID, createOrder)
	s.Require().NoError(err)
	s.Require().Equal(0.02, res)
}

func (s *ServiceSuite) TestCreateClientError() {
	clientErr := errors.New("test client error")

	orderUUID := uuid.New()
	createOrder := model.CreateOrder{
		UserUUID:  uuid.New(),
		PartUUIDs: []uuid.UUID{uuid.New()},
	}

	s.inventoryClient.On("ListParts", s.ctx, []string{createOrder.PartUUIDs[0].String()}).Return(nil, clientErr).Once()
	s.orderRepository.AssertNotCalled(s.T(), "Create", s.ctx, mock.Anything)

	res, err := s.service.Create(s.ctx, orderUUID, createOrder)
	s.Require().Error(err)
	s.Require().ErrorIs(err, clientErr)
	s.Require().Empty(res)
}

func (s *ServiceSuite) TestCreateOrderedPartsNotAvailableError() {
	orderUUID := uuid.New()
	createOrder := model.CreateOrder{
		UserUUID:  uuid.New(),
		PartUUIDs: []uuid.UUID{uuid.New(), uuid.New()},
	}

	// second part is absent
	parts := []model.Part{
		{
			UUID:          createOrder.PartUUIDs[0].String(),
			Price:         0.01,
			StockQuantity: 1,
		},
	}

	s.inventoryClient.On("ListParts", s.ctx, []string{createOrder.PartUUIDs[0].String(), createOrder.PartUUIDs[1].String()}).Return(parts, nil).Once()
	s.orderRepository.AssertNotCalled(s.T(), "Create", s.ctx, mock.Anything)

	res, err := s.service.Create(s.ctx, orderUUID, createOrder)
	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrOrderedPartsNotAvailable)
	s.Require().Empty(res)
}

func (s *ServiceSuite) TestCreateRepoError() {
	repoErr := errors.New("test repo error")

	orderUUID := uuid.New()
	createOrder := model.CreateOrder{
		UserUUID:  uuid.New(),
		PartUUIDs: []uuid.UUID{uuid.New(), uuid.New()},
	}

	parts := []model.Part{
		{
			UUID:          createOrder.PartUUIDs[0].String(),
			Price:         0.01,
			StockQuantity: 1,
		},
		{
			UUID:          createOrder.PartUUIDs[1].String(),
			Price:         0.01,
			StockQuantity: 1,
		},
	}

	order := model.Order{
		UUID:            orderUUID,
		UserUUID:        createOrder.UserUUID,
		PartUUIDs:       createOrder.PartUUIDs,
		TotalPrice:      0.02,
		TransactionUUID: nil,
		PaymentMethod:   nil,
		Status:          model.OrderStatusPendingPayment,
	}

	s.inventoryClient.On("ListParts", s.ctx, []string{createOrder.PartUUIDs[0].String(), createOrder.PartUUIDs[1].String()}).Return(parts, nil).Once()
	s.orderRepository.On("Create", s.ctx, order).Return(repoErr).Once()

	res, err := s.service.Create(s.ctx, orderUUID, createOrder)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Empty(res)
}
