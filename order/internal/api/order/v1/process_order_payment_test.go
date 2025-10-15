package v1_test

import (
	"github.com/google/uuid"

	"github.com/pinai4/spaceship-factory/order/internal/model"
	orderV1 "github.com/pinai4/spaceship-factory/shared/pkg/openapi/order/v1"
)

func (s *APISuite) TestProcessOrderPaymentSuccess() {
	var (
		orderUUID = uuid.New()

		paramsOpenAPI = orderV1.ProcessOrderPaymentParams{
			OrderUUID: orderUUID,
		}

		paymentMethod = model.OrderPaymentMethodInvestorMoney

		requestOpenAPI = &orderV1.ProcessOrderPaymentRequest{
			PaymentMethod: orderV1.ProcessOrderPaymentRequestPaymentMethod(paymentMethod),
		}

		tranUUID                = uuid.New()
		expectedResponseOpenAPI = &orderV1.ProcessOrderPaymentResponse{
			TransactionUUID: tranUUID,
		}
	)

	s.orderService.On("ProcessPayment", s.ctx, orderUUID, paymentMethod).Return(tranUUID, nil)

	res, err := s.api.ProcessOrderPayment(s.ctx, requestOpenAPI, paramsOpenAPI)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(expectedResponseOpenAPI, res)
}

func (s *APISuite) TestProcessOrderPaymentNotFound() {
	var (
		orderUUID = uuid.New()

		paramsOpenAPI = orderV1.ProcessOrderPaymentParams{
			OrderUUID: orderUUID,
		}

		paymentMethod = model.OrderPaymentMethodInvestorMoney

		requestOpenAPI = &orderV1.ProcessOrderPaymentRequest{
			PaymentMethod: orderV1.ProcessOrderPaymentRequestPaymentMethod(paymentMethod),
		}
	)

	s.orderService.On("ProcessPayment", s.ctx, orderUUID, paymentMethod).Return(uuid.UUID{}, model.ErrOrderNotFound)

	res, err := s.api.ProcessOrderPayment(s.ctx, requestOpenAPI, paramsOpenAPI)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	resErr, ok := res.(*orderV1.NotFoundError)
	s.Require().True(ok)
	s.Require().Equal(404, resErr.GetCode())
}
