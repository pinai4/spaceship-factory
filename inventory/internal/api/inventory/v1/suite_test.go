package v1_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	v1 "github.com/pinai4/spaceship-factory/inventory/internal/api/inventory/v1"
	"github.com/pinai4/spaceship-factory/inventory/internal/model"
	"github.com/pinai4/spaceship-factory/inventory/internal/service/mocks"
	inventoryV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/inventory/v1"
)

type APISuite struct {
	suite.Suite

	ctx context.Context

	partService *mocks.PartService

	api inventoryV1.InventoryServiceServer
}

func (s *APISuite) SetupTest() {
	s.ctx = context.Background()

	s.partService = mocks.NewPartService(s.T())

	s.api = v1.NewAPI(
		s.partService,
	)
}

func (s *APISuite) TearDownTest() {
}

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(APISuite))
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
