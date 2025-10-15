package part_test

import (
	"github.com/go-faster/errors"
	"github.com/google/uuid"

	"github.com/pinai4/spaceship-factory/inventory/internal/model"
)

func (s *ServiceSuite) TestListSuccess() {
	var (
		filter = model.PartsFilter{
			UUIDS:                 []string{uuid.NewString(), uuid.NewString()},
			Names:                 []string{"Test Name1", "Test Name2"},
			Categories:            []model.Category{model.CategoryWing},
			ManufacturerCountries: make([]string, 0),
			Tags:                  nil,
		}

		parts = []model.Part{buildTestPart(), buildTestPart()}
	)

	s.partRepository.On("List", s.ctx, filter).Return(parts, nil)

	res, err := s.service.List(s.ctx, filter)
	s.Require().NoError(err)
	s.Require().Equal(parts, res)
}

func (s *ServiceSuite) TestListRepoError() {
	var (
		repoErr = errors.New("test error")

		filter = model.PartsFilter{
			Tags: []string{"Tag1", "Tag2"},
		}
	)

	s.partRepository.On("List", s.ctx, filter).Return([]model.Part{}, repoErr)

	res, err := s.service.List(s.ctx, filter)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Empty(res)
}
