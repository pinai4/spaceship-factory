package v1_test

import (
	"errors"

	"github.com/google/uuid"

	"github.com/pinai4/spaceship-factory/inventory/internal/converter"
	"github.com/pinai4/spaceship-factory/inventory/internal/model"
	inventoryV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/inventory/v1"
)

func (s *APISuite) TestListPartsSuccess() {
	var (
		protoFilter = &inventoryV1.PartsFilter{
			Uuids:                 []string{uuid.NewString(), uuid.NewString()},
			Names:                 []string{"Name1", "Name2", "Name3"},
			Categories:            []inventoryV1.Category{inventoryV1.Category_CATEGORY_FUEL},
			ManufacturerCountries: nil,
			Tags:                  []string{},
		}

		req = &inventoryV1.ListPartsRequest{
			Filter: protoFilter,
		}

		modelFilter = converter.PartsFilterToModel(protoFilter)

		parts = []model.Part{buildTestPart(), buildTestPart()}

		expectedProtoResponse = &inventoryV1.ListPartsResponse{
			Parts: converter.PartsToProto(parts),
		}
	)

	s.partService.On("List", s.ctx, modelFilter).Return(parts, nil)

	res, err := s.api.ListParts(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(expectedProtoResponse, res)
}

func (s *APISuite) TestListPartsServiceError() {
	var (
		serviceErr  = errors.New("test error")
		protoFilter = &inventoryV1.PartsFilter{}

		req = &inventoryV1.ListPartsRequest{
			Filter: protoFilter,
		}

		modelFilter = converter.PartsFilterToModel(protoFilter)
	)

	s.partService.On("List", s.ctx, modelFilter).Return([]model.Part{}, serviceErr)

	res, err := s.api.ListParts(s.ctx, req)
	s.Require().Error(err)
	s.Require().ErrorIs(err, serviceErr)
	s.Require().Nil(res)
}
