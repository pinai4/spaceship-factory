package v1_test

import (
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pinai4/spaceship-factory/inventory/internal/converter"
	"github.com/pinai4/spaceship-factory/inventory/internal/model"
	inventoryV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/inventory/v1"
)

func (s *APISuite) TestGetPartSuccess() {
	var (
		reqPartUUID = uuid.NewString()

		req = &inventoryV1.GetPartRequest{
			Uuid: reqPartUUID,
		}

		part = buildTestPart()

		expectedProtoResponse = &inventoryV1.GetPartResponse{
			Part: converter.PartToProto(part),
		}
	)

	s.partService.On("Get", s.ctx, reqPartUUID).Return(part, nil)

	res, err := s.api.GetPart(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(expectedProtoResponse, res)
}

func (s *APISuite) TestGetPartNotFound() {
	var (
		reqPartUUID = uuid.NewString()

		req = &inventoryV1.GetPartRequest{
			Uuid: reqPartUUID,
		}
	)

	s.partService.On("Get", s.ctx, reqPartUUID).Return(model.Part{}, model.ErrPartNotFound)

	res, err := s.api.GetPart(s.ctx, req)
	s.Require().Error(err)
	s.Require().Nil(res)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.NotFound, st.Code())
}

func (s *APISuite) TestGetPartServiceError() {
	var (
		serviceErr  = errors.New("test error")
		reqPartUUID = uuid.NewString()

		req = &inventoryV1.GetPartRequest{
			Uuid: reqPartUUID,
		}
	)

	s.partService.On("Get", s.ctx, reqPartUUID).Return(model.Part{}, serviceErr)

	res, err := s.api.GetPart(s.ctx, req)
	s.Require().Error(err)
	s.Require().ErrorIs(err, serviceErr)
	s.Require().Nil(res)
}
