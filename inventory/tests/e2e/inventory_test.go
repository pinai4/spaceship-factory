package e2e

import (
	"context"

	inventoryV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/inventory/v1"
)

func (s *E2ESuite) SetupTest() {
	s.ctx, s.ctxCancel = context.WithCancel(s.suiteCtx)
}

func (s *E2ESuite) TearDownTest() {
	err := s.env.ClearPartsCollection(s.ctx)
	s.Require().NoError(err)

	s.ctxCancel()
}

func (s *E2ESuite) TestGetPart() {
	part := s.env.GetTestParts()[0]
	partUUID, err := s.env.InsertTestPart(s.ctx, part)
	s.Require().NoError(err)

	resp, err := s.inventoryV1Client.GetPart(s.ctx, &inventoryV1.GetPartRequest{Uuid: partUUID})
	s.Require().NoError(err)
	s.Require().NotNil(resp.GetPart())
	s.Require().Equal(partUUID, resp.GetPart().Uuid)
	s.Require().NotNil(resp.GetPart().GetMetadata())
	s.Require().Equal(
		part.GetMetadata()["power_kw"].GetDoubleValue(),
		resp.GetPart().GetMetadata()["power_kw"].GetDoubleValue(),
	)
	s.Require().Equal(
		part.GetMetadata()["certified"].GetBoolValue(),
		resp.GetPart().GetMetadata()["certified"].GetBoolValue(),
	)
	s.Require().Equal(
		part.GetMetadata()["serial_number"].GetStringValue(),
		resp.GetPart().GetMetadata()["serial_number"].GetStringValue(),
	)
}

func (s *E2ESuite) TestListParts() {
	part1 := s.env.GetTestParts()[0]
	part1UUID, err := s.env.InsertTestPart(s.ctx, part1)
	s.Require().NoError(err)

	part2 := s.env.GetTestParts()[1]
	part2UUID, err := s.env.InsertTestPart(s.ctx, part2)
	s.Require().NoError(err)

	part3 := s.env.GetTestParts()[2]
	part3UUID, err := s.env.InsertTestPart(s.ctx, part3)
	s.Require().NoError(err)

	s.Run("Filter Empty", func() {
		resp, err := s.inventoryV1Client.ListParts(s.ctx, &inventoryV1.ListPartsRequest{})
		s.Require().NoError(err)
		s.Require().Equal(3, len(resp.GetParts()))
		s.Require().ElementsMatch(
			[]string{part1UUID, part2UUID, part3UUID},
			[]string{resp.GetParts()[0].GetUuid(), resp.GetParts()[1].GetUuid(), resp.GetParts()[2].GetUuid()},
		)
	})

	s.Run("Filter by UUIDs", func() {
		resp, err := s.inventoryV1Client.ListParts(s.ctx, &inventoryV1.ListPartsRequest{
			Filter: &inventoryV1.PartsFilter{
				Uuids: []string{part1UUID, part2UUID},
			},
		})
		s.Require().NoError(err)
		s.Require().Equal(2, len(resp.GetParts()))
		s.Require().ElementsMatch(
			[]string{part1UUID, part2UUID},
			[]string{resp.GetParts()[0].GetUuid(), resp.GetParts()[1].GetUuid()},
		)
	})
}
