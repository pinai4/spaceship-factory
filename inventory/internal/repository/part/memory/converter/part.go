package converter

import (
	"github.com/pinai4/spaceship-factory/inventory/internal/model"
	repoModel "github.com/pinai4/spaceship-factory/inventory/internal/repository/part/memory/model"
)

func PartToModel(part repoModel.Part) model.Part {
	return model.Part{
		UUID:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      model.Category(part.Category),
		Dimensions:    model.Dimensions(part.Dimensions),
		Manufacturer:  model.Manufacturer(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      part.Metadata,
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
	}
}

func PartToRepoModel(part model.Part) repoModel.Part {
	return repoModel.Part{
		UUID:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      repoModel.Category(part.Category),
		Dimensions:    repoModel.Dimensions(part.Dimensions),
		Manufacturer:  repoModel.Manufacturer(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      part.Metadata,
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
	}
}
