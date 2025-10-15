package v1_test

import (
	"github.com/google/uuid"

	paymentV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/payment/v1"
)

func (s *APISuite) TestPayOrderSuccess() {
	var (
		expectedTranUUID = uuid.NewString()

		req = &paymentV1.PayOrderRequest{
			OrderUuid:     uuid.NewString(),
			UserUuid:      uuid.NewString(),
			PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
		}
	)

	s.paymentService.On("PayOrder", s.ctx).Return(expectedTranUUID, nil)

	res, err := s.api.PayOrder(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(expectedTranUUID, res.TransactionUuid)
}
