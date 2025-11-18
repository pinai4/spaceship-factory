//go:build integration

package e2e

import (
	"context"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/pinai4/spaceship-factory/platform/pkg/testcontainers/app"
	"github.com/pinai4/spaceship-factory/platform/pkg/testcontainers/mongo"
	"github.com/pinai4/spaceship-factory/platform/pkg/testcontainers/network"
	inventoryV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/inventory/v1"
)

// TestEnvironment — structure for storing test environment resources
type TestEnvironment struct {
	Network *network.Network
	Mongo   *mongo.Container
	Auth    *app.Container
	App     *app.Container
}

func (e *TestEnvironment) InsertTestPart(ctx context.Context, part *inventoryV1.Part) (string, error) {
	parseMetadata := func(metadataProto map[string]*inventoryV1.Value) map[string]any {
		result := make(map[string]any, len(metadataProto))

		for k, v := range metadataProto {
			if v == nil {
				result[k] = nil
				continue
			}

			switch kind := v.Kind.(type) {
			case *inventoryV1.Value_StringValue:
				result[k] = kind.StringValue
			case *inventoryV1.Value_Int64Value:
				result[k] = kind.Int64Value
			case *inventoryV1.Value_DoubleValue:
				result[k] = kind.DoubleValue
			case *inventoryV1.Value_BoolValue:
				result[k] = kind.BoolValue
			}
		}

		return result
	}

	partDoc := bson.M{
		"_id":            part.Uuid,
		"name":           part.Name,
		"description":    part.Description,
		"price":          part.Price,
		"stock_quantity": part.StockQuantity,
		"category":       strings.TrimPrefix(part.Category.String(), "CATEGORY_"),
		"dimensions": bson.M{
			"length": part.Dimensions.Length,
			"width":  part.Dimensions.Width,
			"height": part.Dimensions.Height,
			"weight": part.Dimensions.Weight,
		},
		"manufacturer": bson.M{
			"name":    part.Manufacturer.Name,
			"country": part.Manufacturer.Country,
			"website": part.Manufacturer.Website,
		},
		"tags":       part.Tags,
		"metadata":   parseMetadata(part.Metadata),
		"created_at": primitive.NewDateTimeFromTime(part.CreatedAt.AsTime()),
	}

	_, err := e.Mongo.Client().Database(e.Mongo.Config().Database).Collection(mongoPartsCollectionName).InsertOne(ctx, partDoc)
	if err != nil {
		return "", err
	}

	return part.Uuid, nil
}

func (e *TestEnvironment) GetTestParts() []*inventoryV1.Part {
	now := timestamppb.New(time.Now())

	part1 := &inventoryV1.Part{
		Uuid:          "00000000-0000-0000-0000-000000000001",
		Name:          "Turbo Engine X200",
		Description:   "High-performance turbo engine suitable for small aircraft.",
		Price:         125000.50,
		StockQuantity: 8,
		Category:      inventoryV1.Category_CATEGORY_ENGINE,
		Dimensions: &inventoryV1.Dimensions{
			Length: 120.5,
			Width:  80.2,
			Height: 95.3,
			Weight: 450.0,
		},
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    "AeroTech Industries",
			Country: "USA",
			Website: "https://aerotech.example.com",
		},
		Tags: []string{"engine", "turbo", "aircraft"},
		Metadata: map[string]*inventoryV1.Value{
			"power_kw":      {Kind: &inventoryV1.Value_DoubleValue{DoubleValue: 980.5}},
			"certified":     {Kind: &inventoryV1.Value_BoolValue{BoolValue: true}},
			"serial_number": {Kind: &inventoryV1.Value_StringValue{StringValue: "SN-ENGX200-001"}},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	part2 := &inventoryV1.Part{
		Uuid:          "00000000-0000-0000-0000-000000000002",
		Name:          "Titanium Wing Panel",
		Description:   "Lightweight titanium alloy wing panel with anti-corrosion coating.",
		Price:         32000.0,
		StockQuantity: 25,
		Category:      inventoryV1.Category_CATEGORY_WING,
		Dimensions: &inventoryV1.Dimensions{
			Length: 250.0,
			Width:  60.0,
			Height: 5.0,
			Weight: 120.0,
		},
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    "SkyMetal Works",
			Country: "Germany",
			Website: "https://skymetal.example.com",
		},
		Tags: []string{"wing", "titanium", "aircraft"},
		Metadata: map[string]*inventoryV1.Value{
			"material":     {Kind: &inventoryV1.Value_StringValue{StringValue: "Titanium Alloy"}},
			"batch_number": {Kind: &inventoryV1.Value_Int64Value{Int64Value: 20241001}},
			"is_tested":    {Kind: &inventoryV1.Value_BoolValue{BoolValue: true}},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	part3 := &inventoryV1.Part{
		Uuid:          "00000000-0000-0000-0000-000000000003",
		Name:          "Fusion Rotor F500",
		Description:   "Advanced fusion rotor providing high efficiency for medium-sized aircraft.",
		Price:         89000.75,
		StockQuantity: 12,
		Category:      inventoryV1.Category_CATEGORY_ENGINE,
		Dimensions: &inventoryV1.Dimensions{
			Length: 140.0,
			Width:  85.0,
			Height: 100.0,
			Weight: 380.0,
		},
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    "AeroFusion Corp",
			Country: "Japan",
			Website: "https://aerofusion.example.com",
		},
		Tags: []string{"engine", "rotor", "aircraft", "fusion"},
		Metadata: map[string]*inventoryV1.Value{
			"power_kw":       {Kind: &inventoryV1.Value_DoubleValue{DoubleValue: 760.2}},
			"certified":      {Kind: &inventoryV1.Value_BoolValue{BoolValue: true}},
			"serial_number":  {Kind: &inventoryV1.Value_StringValue{StringValue: "SN-FR-F500-007"}},
			"warranty_years": {Kind: &inventoryV1.Value_Int64Value{Int64Value: 5}},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	return []*inventoryV1.Part{part1, part2, part3}
}

func (e *TestEnvironment) ClearPartsCollection(ctx context.Context) error {
	_, err := e.Mongo.Client().Database(e.Mongo.Config().Database).Collection(mongoPartsCollectionName).DeleteMany(ctx, bson.M{})
	if err != nil {
		return err
	}

	return nil
}
