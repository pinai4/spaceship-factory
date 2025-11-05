//go:build unit || !integration

package assembly_test

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/pinai4/spaceship-factory/assembly/internal/model"
)

func (s *ServiceSuite) TestAssembleShipSuccess() {
	orderUUID := uuid.NewString()
	userUUID := uuid.NewString()

	s.assemblyProducer.On("ProduceShipAssembled", s.ctx, mock.MatchedBy(func(e model.ShipAssembledEvent) bool {
		return e.OrderUUID == orderUUID &&
			e.UserUUID == userUUID
		// EventUUID is intentionally ignored
	})).Return(nil).Once()

	err := s.service.AssembleShip(s.ctx, orderUUID, userUUID, 100*time.Millisecond)
	s.Require().NoError(err)
}

func (s *ServiceSuite) TestAssembleShipProducerError() {
	producerErr := errors.New("test producer error")

	orderUUID := uuid.NewString()
	userUUID := uuid.NewString()

	s.assemblyProducer.On("ProduceShipAssembled", s.ctx, mock.MatchedBy(func(e model.ShipAssembledEvent) bool {
		return e.OrderUUID == orderUUID &&
			e.UserUUID == userUUID
		// EventUUID is intentionally ignored
	})).Return(producerErr).Once()

	err := s.service.AssembleShip(s.ctx, orderUUID, userUUID, 100*time.Millisecond)
	s.Require().Error(err)
	s.Require().ErrorIs(err, producerErr)
}
