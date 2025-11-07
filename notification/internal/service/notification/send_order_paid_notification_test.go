//go:build unit || !integration

package notification_test

import (
	"github.com/google/uuid"

	"github.com/pinai4/spaceship-factory/notification/internal/model"
)

func (s *ServiceSuite) TestSendOrderPaidNotificationSuccess() {
	event := model.OrderPaidEvent{
		EventUUID:       uuid.New().String(),
		OrderUUID:       uuid.New().String(),
		UserUUID:        uuid.New().String(),
		PaymentMethod:   "CARD",
		TransactionUUID: uuid.New().String(),
	}

	err := s.service.SendOrderPaidNotification(s.ctx, event)
	s.Require().NoError(err)
}
