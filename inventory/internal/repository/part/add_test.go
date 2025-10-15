package part_test

import "github.com/pinai4/spaceship-factory/inventory/internal/model"

func (s *RepositorySuite) TestAddSuccess() {
	// Arrange
	part := buildTestPart()
	partUUID := part.UUID

	// Act
	err := s.repository.Add(s.ctx, part)

	// Assert
	s.Require().NoError(err)

	count := s.countParts()
	s.Require().Equal(count, 1)

	res, err := s.repository.Get(s.ctx, partUUID)
	s.Require().NoError(err)
	s.Require().Equal(part, res)
}

func (s *RepositorySuite) TestAddPartAlreadyExistsError() {
	// Arrange
	part := buildTestPart()
	_ = s.repository.Add(s.ctx, part)

	// Act
	err := s.repository.Add(s.ctx, part)

	// Assert
	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrPartAlreadyExists)
}
