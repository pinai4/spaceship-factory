package memory_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/pinai4/spaceship-factory/inventory/internal/model"
	"github.com/pinai4/spaceship-factory/inventory/internal/repository"
	"github.com/pinai4/spaceship-factory/inventory/internal/repository/part/memory"
)

type RepositorySuite struct {
	suite.Suite

	ctx context.Context

	repository repository.PartRepository
}

func (s *RepositorySuite) SetupTest() {
	s.ctx = context.Background()

	s.repository = memory.NewRepository()
}

func (s *RepositorySuite) TearDownTest() {
}

func TestRepositoryIntegration(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}

func (s *RepositorySuite) countParts() int {
	s.T().Helper()

	parts, err := s.repository.List(s.ctx, model.PartsFilter{})
	s.Require().NoError(err)

	return len(parts)
}

func buildTestPart() model.Part {
	return model.Part{
		UUID:          uuid.NewString(),
		Name:          "Turbo Engine X200",
		Description:   "High-performance turbo engine suitable for small aircraft.",
		Price:         125000.50,
		StockQuantity: 8,
		Category:      model.CategoryEngine,
		Dimensions: model.Dimensions{
			Length: 120.5,
			Width:  80.2,
			Height: 95.3,
			Weight: 450.0,
		},
		Manufacturer: model.Manufacturer{
			Name:    "AeroTech Industries",
			Country: "USA",
			Website: "https://aerotech.example.com",
		},
		Tags: []string{"engine", "turbo", "aircraft"},
		Metadata: map[string]any{
			"power_kw":      980.5,
			"certified":     true,
			"serial_number": "SN-ENGX200-001",
		},
		CreatedAt: time.Now(),
	}
}
