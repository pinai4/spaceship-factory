package payment

import def "github.com/pinai4/spaceship-factory/payment/internal/service"

var _ def.PaymentService = (*service)(nil)

type service struct{}

func NewService() *service {
	return &service{}
}
