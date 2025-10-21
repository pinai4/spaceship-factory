package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/pinai4/spaceship-factory/inventory/internal/model"
	repoConverter "github.com/pinai4/spaceship-factory/inventory/internal/repository/part/mongodb/converter"
)

func (r *repository) Add(ctx context.Context, part model.Part) error {
	repoPart := repoConverter.PartToRepoModel(part)

	_, err := r.collection.InsertOne(ctx, repoPart)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return model.ErrPartAlreadyExists
		}
		return fmt.Errorf("PartRepository.Add add part error: %w", err)
	}

	return nil
}
