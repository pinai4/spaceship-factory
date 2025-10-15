package part_test

import (
	"errors"

	"github.com/google/uuid"

	"github.com/pinai4/spaceship-factory/inventory/internal/model"
)

func (s *ServiceSuite) TestGetSuccess() {
	part := buildTestPart()
	partUUID := part.UUID

	s.partRepository.On("Get", s.ctx, partUUID).Return(part, nil)

	res, err := s.service.Get(s.ctx, partUUID)
	s.Require().NoError(err)
	s.Require().Equal(part, res)
}

func (s *ServiceSuite) TestGetRepoError() {
	var (
		repoErr  = errors.New("test repo error")
		partUUID = uuid.NewString()
	)

	s.partRepository.On("Get", s.ctx, partUUID).Return(model.Part{}, repoErr)

	res, err := s.service.Get(s.ctx, partUUID)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Empty(res)
}
