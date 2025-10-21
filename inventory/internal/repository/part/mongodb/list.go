package mongodb

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/pinai4/spaceship-factory/inventory/internal/model"
	repoConverter "github.com/pinai4/spaceship-factory/inventory/internal/repository/part/mongodb/converter"
	repoModel "github.com/pinai4/spaceship-factory/inventory/internal/repository/part/mongodb/model"
)

func (r *repository) List(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	bsonFilter := bson.M{}

	// UUID filter
	if len(filter.UUIDS) > 0 {
		bsonFilter["_id"] = bson.M{"$in": filter.UUIDS}
	}

	// Name filter
	if len(filter.Names) > 0 {
		bsonFilter["name"] = bson.M{"$in": filter.Names}
	}

	// Category filter
	if len(filter.Categories) > 0 {
		bsonFilter["category"] = bson.M{"$in": filter.Categories}
	}

	// Manufacturer country filter (nested field)
	if len(filter.ManufacturerCountries) > 0 {
		bsonFilter["manufacturer.country"] = bson.M{"$in": filter.ManufacturerCountries}
	}

	// Tags filter — matches if *any* tag is present in the document’s tags array
	if len(filter.Tags) > 0 {
		bsonFilter["tags"] = bson.M{"$in": filter.Tags}
	}

	cursor, err := r.collection.Find(ctx, bsonFilter)
	if err != nil {
		return nil, fmt.Errorf("PartRepository.List collection.Find error: %w", err)
	}
	defer func() {
		if cerr := cursor.Close(ctx); cerr != nil {
			log.Printf("failed to close cursor: %v\n", cerr)
		}
	}()

	var repoParts []repoModel.Part
	if err := cursor.All(ctx, &repoParts); err != nil {
		return nil, fmt.Errorf("PartRepository.List cursor.All error: %w", err)
	}

	return repoConverter.PartsToModels(repoParts), nil
}
