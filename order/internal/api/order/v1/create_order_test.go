//go:build unit || !integration

package v1_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/pinai4/spaceship-factory/order/internal/model"
	grpcAuth "github.com/pinai4/spaceship-factory/platform/pkg/middleware/grpc"
	orderV1 "github.com/pinai4/spaceship-factory/shared/pkg/openapi/order/v1"
	commonV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/common/v1"
)

func (s *APISuite) TestCreateOrder_Success() {
	sessUserUUID := uuid.New()

	createOrder := model.CreateOrder{
		UserUUID:  sessUserUUID,
		PartUUIDs: []uuid.UUID{uuid.New(), uuid.New()},
	}

	requestOpenAPI := &orderV1.CreateOrderRequest{
		UserUUID:  createOrder.UserUUID,
		PartUuids: createOrder.PartUUIDs,
	}

	totalPrice := 0.02

	authCtx := grpcAuth.AddUserToContext(s.ctx, &commonV1.User{Uuid: sessUserUUID.String()})

	s.orderService.On("Create", authCtx, mock.Anything, createOrder).Return(totalPrice, nil)

	res, err := s.api.CreateOrder(authCtx, requestOpenAPI, orderV1.CreateOrderParams{})
	s.Require().NoError(err)
	s.Require().NotNil(res)

	resImpl, ok := res.(*orderV1.CreateOrderResponse)
	s.Require().True(ok)
	s.Require().Equal(totalPrice, resImpl.GetTotalPrice())
}

func (s *APISuite) TestCreateOrder_BadRequestSessionUserConflictsError() {
	createOrder := model.CreateOrder{
		UserUUID:  uuid.New(),
		PartUUIDs: []uuid.UUID{},
	}

	requestOpenAPI := &orderV1.CreateOrderRequest{
		UserUUID:  createOrder.UserUUID,
		PartUuids: createOrder.PartUUIDs,
	}

	s.orderService.AssertNotCalled(s.T(), "Create", s.ctx, mock.Anything, mock.Anything)

	res, err := s.api.CreateOrder(s.ctx, requestOpenAPI, orderV1.CreateOrderParams{})
	s.Require().NoError(err)
	s.Require().NotNil(res)

	resErr, ok := res.(*orderV1.BadRequestError)
	s.Require().True(ok)
	s.Require().Equal(400, resErr.GetCode())
}

func (s *APISuite) TestCreateOrder_BadRequestPartsListEmptyError() {
	sessUserUUID := uuid.New()

	createOrder := model.CreateOrder{
		UserUUID:  sessUserUUID,
		PartUUIDs: []uuid.UUID{},
	}

	requestOpenAPI := &orderV1.CreateOrderRequest{
		UserUUID:  createOrder.UserUUID,
		PartUuids: createOrder.PartUUIDs,
	}

	authCtx := grpcAuth.AddUserToContext(s.ctx, &commonV1.User{Uuid: sessUserUUID.String()})

	s.orderService.AssertNotCalled(s.T(), "Create", s.ctx, mock.Anything, mock.Anything)

	res, err := s.api.CreateOrder(authCtx, requestOpenAPI, orderV1.CreateOrderParams{})
	s.Require().NoError(err)
	s.Require().NotNil(res)

	resErr, ok := res.(*orderV1.BadRequestError)
	s.Require().True(ok)
	s.Require().Equal(400, resErr.GetCode())
}

func (s *APISuite) TestCreateOrder_BadRequestPartsNotAvailableError() {
	sessUserUUID := uuid.New()

	createOrder := model.CreateOrder{
		UserUUID:  sessUserUUID,
		PartUUIDs: []uuid.UUID{uuid.New(), uuid.New()},
	}

	requestOpenAPI := &orderV1.CreateOrderRequest{
		UserUUID:  createOrder.UserUUID,
		PartUuids: createOrder.PartUUIDs,
	}

	authCtx := grpcAuth.AddUserToContext(s.ctx, &commonV1.User{Uuid: sessUserUUID.String()})

	s.orderService.On("Create", authCtx, mock.Anything, createOrder).Return(float64(0), model.ErrOrderedPartsNotAvailable)

	res, err := s.api.CreateOrder(authCtx, requestOpenAPI, orderV1.CreateOrderParams{})
	s.Require().NoError(err)
	s.Require().NotNil(res)

	resErr, ok := res.(*orderV1.BadRequestError)
	s.Require().True(ok)
	s.Require().Equal(400, resErr.GetCode())
}
