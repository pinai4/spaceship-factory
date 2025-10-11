package converter

import (
	"reflect"
	"strings"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/pinai4/spaceship-factory/inventory/internal/model"
	inventoryV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/inventory/v1"
)

func PartToProto(part model.Part) *inventoryV1.Part {
	parseCategory := func(mCategory model.Category) inventoryV1.Category {
		searchVal := "CATEGORY_" + string(mCategory)
		if val, ok := inventoryV1.Category_value[searchVal]; ok {
			return inventoryV1.Category(val)
		}
		return inventoryV1.Category_CATEGORY_UNSPECIFIED
	}

	parseMetadata := func(m map[string]any) map[string]*inventoryV1.Value {
		if m == nil {
			return nil
		}

		res := make(map[string]*inventoryV1.Value, len(m))
		for k, v := range m {
			switch val := v.(type) {
			case string:
				res[k] = &inventoryV1.Value{Kind: &inventoryV1.Value_StringValue{StringValue: val}}
			case int, int64:
				res[k] = &inventoryV1.Value{Kind: &inventoryV1.Value_Int64Value{Int64Value: reflect.ValueOf(val).Int()}}
			case float32, float64:
				res[k] = &inventoryV1.Value{Kind: &inventoryV1.Value_DoubleValue{DoubleValue: reflect.ValueOf(val).Float()}}
			case bool:
				res[k] = &inventoryV1.Value{Kind: &inventoryV1.Value_BoolValue{BoolValue: val}}
			}
		}
		return res
	}

	var updatedAt *timestamppb.Timestamp
	if part.UpdatedAt != nil {
		updatedAt = timestamppb.New(*part.UpdatedAt)
	}

	return &inventoryV1.Part{
		Uuid:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      parseCategory(part.Category),
		Dimensions: &inventoryV1.Dimensions{
			Length: part.Dimensions.Length,
			Width:  part.Dimensions.Width,
			Height: part.Dimensions.Height,
			Weight: part.Dimensions.Weight,
		},
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    part.Manufacturer.Name,
			Country: part.Manufacturer.Country,
			Website: part.Manufacturer.Website,
		},
		Tags:      part.Tags,
		Metadata:  parseMetadata(part.Metadata),
		CreatedAt: timestamppb.New(part.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

func PartsFilterToModel(filter *inventoryV1.PartsFilter) model.PartsFilter {
	if filter == nil {
		return model.PartsFilter{}
	}

	parseCategoryList := func(list []inventoryV1.Category) []model.Category {
		categories := make([]model.Category, len(list))
		for i, pCategory := range list {
			category := model.Category(strings.TrimPrefix(pCategory.String(), "CATEGORY_"))
			categories[i] = category
		}
		return categories
	}

	return model.PartsFilter{
		UUIDS:                 filter.GetUuids(),
		Names:                 filter.GetNames(),
		Categories:            parseCategoryList(filter.GetCategories()),
		ManufacturerCountries: filter.GetManufacturerCountries(),
		Tags:                  filter.GetTags(),
	}
}

func PartsToProto(list []model.Part) []*inventoryV1.Part {
	parts := make([]*inventoryV1.Part, len(list))
	for i, p := range list {
		parts[i] = PartToProto(p)
	}

	return parts
}
