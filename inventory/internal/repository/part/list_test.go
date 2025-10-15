package part_test

import (
	"github.com/pinai4/spaceship-factory/inventory/internal/model"
)

func (s *RepositorySuite) TestListSuccess() {
	part1 := buildTestPart()
	part1UUID := part1.UUID
	part1.Category = model.CategoryPorthole

	part2 := buildTestPart()
	part2UUID := part2.UUID

	part3 := buildTestPart()
	part3.Tags = []string{"tag1", "tag2"}

	part4 := buildTestPart()
	part4.Tags = []string{"tag2", "tag3"}

	part5 := buildTestPart()
	_ = s.repository.Add(s.ctx, part1)
	_ = s.repository.Add(s.ctx, part2)
	_ = s.repository.Add(s.ctx, part3)
	_ = s.repository.Add(s.ctx, part4)
	_ = s.repository.Add(s.ctx, part5)

	tests := []struct {
		name     string
		filter   model.PartsFilter
		expected []model.Part
	}{
		{
			name:     "Filter Empty",
			filter:   model.PartsFilter{},
			expected: []model.Part{part1, part2, part4, part5, part3},
		},
		{
			name: "Filter by UUIDs",
			filter: model.PartsFilter{
				UUIDS: []string{part1UUID, part2UUID},
			},
			expected: []model.Part{part1, part2},
		},
		{
			name: "Filter by UUIDs and Categories",
			filter: model.PartsFilter{
				UUIDS:      []string{part1UUID, part2UUID},
				Categories: []model.Category{model.CategoryPorthole},
			},
			expected: []model.Part{part1},
		},
		{
			name: "Filter by Tags",
			filter: model.PartsFilter{
				Tags: []string{"tag2"},
			},
			expected: []model.Part{part3, part4},
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			res, err := s.repository.List(s.ctx, test.filter)
			s.Require().NoError(err)
			s.Require().ElementsMatch(test.expected, res)
		})
	}
}
