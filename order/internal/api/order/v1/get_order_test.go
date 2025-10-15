package v1_test

import (
	"errors"

	"github.com/google/uuid"

	"github.com/pinai4/spaceship-factory/order/internal/converter"
	"github.com/pinai4/spaceship-factory/order/internal/model"
	orderV1 "github.com/pinai4/spaceship-factory/shared/pkg/openapi/order/v1"
)

func (s *APISuite) TestGetOrderSuccess() {
	var (
		order     = buildTestOrder()
		orderUUID = order.UUID

		paramsOpenAPI = orderV1.GetOrderParams{
			OrderUUID: orderUUID,
		}

		expectedResponseOpenAPI = converter.OrderToOpenAPI(order)
	)

	s.orderService.On("Get", s.ctx, orderUUID).Return(order, nil)

	res, err := s.api.GetOrder(s.ctx, paramsOpenAPI)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(expectedResponseOpenAPI, res)
}

func (s *APISuite) TestGetOrderNotFound() {
	var (
		orderUUID = uuid.New()

		paramsOpenAPI = orderV1.GetOrderParams{
			OrderUUID: orderUUID,
		}
	)

	s.orderService.On("Get", s.ctx, orderUUID).Return(model.Order{}, model.ErrOrderNotFound)

	res, err := s.api.GetOrder(s.ctx, paramsOpenAPI)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	resErr, ok := res.(*orderV1.NotFoundError)
	s.Require().True(ok)
	s.Require().Equal(404, resErr.GetCode())
}

func (s *APISuite) TestGetOrderServiceError() {
	var (
		serviceErr = errors.New("test error")

		orderUUID = uuid.New()

		paramsOpenAPI = orderV1.GetOrderParams{
			OrderUUID: orderUUID,
		}
	)

	s.orderService.On("Get", s.ctx, orderUUID).Return(model.Order{}, serviceErr)

	res, err := s.api.GetOrder(s.ctx, paramsOpenAPI)
	s.Require().Error(err)
	s.Require().ErrorIs(err, serviceErr)
	s.Require().Nil(res)
}
