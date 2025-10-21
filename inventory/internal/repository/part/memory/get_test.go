package memory_test

import (
	"github.com/google/uuid"

	"github.com/pinai4/spaceship-factory/inventory/internal/model"
)

func (s *RepositorySuite) TestGetSuccess() {
	// Arrange
	part := buildTestPart()
	partUUID := part.UUID
	_ = s.repository.Add(s.ctx, part)
	_ = s.repository.Add(s.ctx, buildTestPart())

	// Act
	res, err := s.repository.Get(s.ctx, partUUID)

	// Assert
	s.Require().NoError(err)
	s.Require().Equal(part, res)
}

func (s *RepositorySuite) TestGetNotFound() {
	// Arrange
	partUUID := uuid.NewString()

	// Act
	res, err := s.repository.Get(s.ctx, partUUID)

	// Assert
	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrPartNotFound)
	s.Require().Empty(res)
}
