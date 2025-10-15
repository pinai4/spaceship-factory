package payment_test

func (s *ServiceSuite) TestPayOrderSuccess() {
	_, err := s.service.PayOrder(s.ctx)
	s.Require().NoError(err)
}
