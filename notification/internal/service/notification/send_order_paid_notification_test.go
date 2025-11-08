//go:build unit || !integration

package notification_test

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

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

	s.telegramClient.On("SendMessage", s.ctx, chatID, mock.MatchedBy(func(m string) bool {
		return strings.Contains(m, event.OrderUUID) &&
			strings.Contains(m, event.UserUUID) &&
			strings.Contains(m, event.PaymentMethod) &&
			strings.Contains(m, event.TransactionUUID)
	})).Return(nil).Once()

	err := s.service.SendOrderPaidNotification(s.ctx, event)
	s.Require().NoError(err)
}

func (s *ServiceSuite) TestSendOrderPaidNotificationTelegramClientError() {
	telegramClientErr := errors.New("test client error")

	event := model.OrderPaidEvent{
		EventUUID:       uuid.New().String(),
		OrderUUID:       uuid.New().String(),
		UserUUID:        uuid.New().String(),
		PaymentMethod:   "CARD",
		TransactionUUID: uuid.New().String(),
	}

	s.telegramClient.On("SendMessage", s.ctx, chatID, mock.MatchedBy(func(m string) bool {
		return strings.Contains(m, event.OrderUUID) &&
			strings.Contains(m, event.UserUUID) &&
			strings.Contains(m, event.PaymentMethod) &&
			strings.Contains(m, event.TransactionUUID)
	})).Return(telegramClientErr).Once()

	err := s.service.SendOrderPaidNotification(s.ctx, event)
	s.Require().Error(err)
	s.Require().ErrorIs(err, telegramClientErr)
}
