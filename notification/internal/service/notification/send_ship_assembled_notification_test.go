//go:build unit || !integration

package notification_test

import (
	"github.com/google/uuid"

	"github.com/pinai4/spaceship-factory/notification/internal/model"
)

func (s *ServiceSuite) TestSendShipAssembledNotificationSuccess() {
	event := model.ShipAssembledEvent{
		EventUUID:    uuid.New().String(),
		OrderUUID:    uuid.New().String(),
		UserUUID:     uuid.New().String(),
		BuildTimeSec: 10,
	}

	err := s.service.SendShipAssembledNotification(s.ctx, event)
	s.Require().NoError(err)
}
