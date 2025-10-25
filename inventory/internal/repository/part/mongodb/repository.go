package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/pinai4/spaceship-factory/inventory/internal/model"
	def "github.com/pinai4/spaceship-factory/inventory/internal/repository"
)

const collectionName = "parts"

var _ def.PartRepository = (*repository)(nil)

type repository struct {
	collection *mongo.Collection
}

func NewRepository(ctx context.Context, db *mongo.Database) *repository {
	collection := db.Collection(collectionName)

	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "name", Value: 1}},
			Options: options.Index().SetUnique(false),
		},
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := collection.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		panic(err)
	}

	repo := &repository{
		collection: collection,
	}

	if err := repo.seed(ctx); err != nil {
		panic(fmt.Errorf("seeding failed: %w", err))
	}

	return repo
}

func (r *repository) seed(ctx context.Context) error {
	list, err := r.List(ctx, model.PartsFilter{})
	if err != nil {
		return err
	}
	if len(list) > 0 {
		return nil
	}

	now := time.Now()

	part1 := model.Part{
		UUID:          uuid.Nil.String(),
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
		CreatedAt: now,
		UpdatedAt: &now,
	}
	if err := r.Add(ctx, part1); err != nil {
		return err
	}

	part2 := model.Part{
		UUID:          uuid.NewString(),
		Name:          "Titanium Wing Panel",
		Description:   "Lightweight titanium alloy wing panel with anti-corrosion coating.",
		Price:         32000.0,
		StockQuantity: 25,
		Category:      model.CategoryWing,
		Dimensions: model.Dimensions{
			Length: 250.0,
			Width:  60.0,
			Height: 5.0,
			Weight: 120.0,
		},
		Manufacturer: model.Manufacturer{
			Name:    "SkyMetal Works",
			Country: "Germany",
			Website: "https://skymetal.example.com",
		},
		Tags: []string{"wing", "titanium", "aircraft"},
		Metadata: map[string]any{
			"material":     "Titanium Alloy",
			"batch_number": 20241001,
			"is_tested":    true,
		},
		CreatedAt: now,
	}
	if err := r.Add(ctx, part2); err != nil {
		return err
	}

	return nil
}
