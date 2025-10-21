package mongodb

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/pinai4/spaceship-factory/inventory/internal/model"
	repoConverter "github.com/pinai4/spaceship-factory/inventory/internal/repository/part/mongodb/converter"
	repoModel "github.com/pinai4/spaceship-factory/inventory/internal/repository/part/mongodb/model"
)

func (r *repository) Get(ctx context.Context, uuid string) (model.Part, error) {
	var repoPart repoModel.Part
	err := r.collection.FindOne(ctx, bson.M{"_id": uuid}).Decode(&repoPart)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Part{}, model.ErrPartNotFound
		}
		return model.Part{}, fmt.Errorf("PartRepository.Get get part error: %w", err)
	}

	return repoConverter.PartToModel(repoPart), nil
}
