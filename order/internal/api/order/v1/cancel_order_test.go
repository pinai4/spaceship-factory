package v1_test

import (
	"github.com/google/uuid"

	"github.com/pinai4/spaceship-factory/order/internal/model"
	orderV1 "github.com/pinai4/spaceship-factory/shared/pkg/openapi/order/v1"
)

func (s *APISuite) TestCancelOrderSuccess() {
	var (
		orderUUID = uuid.New()

		paramsOpenAPI = orderV1.CancelOrderParams{
			OrderUUID: orderUUID,
		}

		expectedResponseOpenAPI = &orderV1.CancelOrderNoContent{}
	)

	s.orderService.On("Cancel", s.ctx, orderUUID).Return(nil)

	res, err := s.api.CancelOrder(s.ctx, paramsOpenAPI)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(expectedResponseOpenAPI, res)
}

func (s *APISuite) TestCancelOrderNotFound() {
	var (
		orderUUID = uuid.New()

		paramsOpenAPI = orderV1.CancelOrderParams{
			OrderUUID: orderUUID,
		}
	)

	s.orderService.On("Cancel", s.ctx, orderUUID).Return(model.ErrOrderNotFound)

	res, err := s.api.CancelOrder(s.ctx, paramsOpenAPI)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	resErr, ok := res.(*orderV1.NotFoundError)
	s.Require().True(ok)
	s.Require().Equal(404, resErr.GetCode())
}

func (s *APISuite) TestCancelOrderConflictError() {
	var (
		orderUUID = uuid.New()

		paramsOpenAPI = orderV1.CancelOrderParams{
			OrderUUID: orderUUID,
		}
	)

	s.orderService.On("Cancel", s.ctx, orderUUID).Return(model.ErrOrderCancelNotAllowed)

	res, err := s.api.CancelOrder(s.ctx, paramsOpenAPI)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	resErr, ok := res.(*orderV1.ConflictError)
	s.Require().True(ok)
	s.Require().Equal(409, resErr.GetCode())
}
