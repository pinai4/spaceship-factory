//go:build unit || !integration

package notification_test

import (
	"errors"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/pinai4/spaceship-factory/notification/internal/model"
)

func (s *ServiceSuite) TestSendShipAssembledNotificationSuccess() {
	event := model.ShipAssembledEvent{
		EventUUID:    uuid.New().String(),
		OrderUUID:    uuid.New().String(),
		UserUUID:     uuid.New().String(),
		BuildTimeSec: 1000001,
	}

	s.telegramClient.On("SendMessage", s.ctx, chatID, mock.MatchedBy(func(m string) bool {
		return strings.Contains(m, event.OrderUUID) &&
			strings.Contains(m, event.UserUUID) &&
			strings.Contains(m, strconv.FormatInt(event.BuildTimeSec, 10))
	})).Return(nil).Once()

	err := s.service.SendShipAssembledNotification(s.ctx, event)
	s.Require().NoError(err)
}

func (s *ServiceSuite) TestSendShipAssembledNotificationTelegramClientError() {
	telegramClientErr := errors.New("test client error")

	event := model.ShipAssembledEvent{
		EventUUID:    uuid.New().String(),
		OrderUUID:    uuid.New().String(),
		UserUUID:     uuid.New().String(),
		BuildTimeSec: 1000001,
	}

	s.telegramClient.On("SendMessage", s.ctx, chatID, mock.MatchedBy(func(m string) bool {
		return strings.Contains(m, event.OrderUUID) &&
			strings.Contains(m, event.UserUUID) &&
			strings.Contains(m, strconv.FormatInt(event.BuildTimeSec, 10))
	})).Return(telegramClientErr).Once()

	err := s.service.SendShipAssembledNotification(s.ctx, event)
	s.Require().Error(err)
	s.Require().ErrorIs(err, telegramClientErr)
}
