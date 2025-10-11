package part

import (
	"context"
	"slices"

	"github.com/pinai4/spaceship-factory/inventory/internal/model"
	repoConverter "github.com/pinai4/spaceship-factory/inventory/internal/repository/converter"
)

func (r *repository) List(_ context.Context, filter model.PartsFilter) ([]model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	hasTagListsIntersectionFunc := func(a, b []string) bool {
		set := make(map[string]struct{}, len(a))
		for _, v := range a {
			set[v] = struct{}{}
		}

		for _, v := range b {
			if _, ok := set[v]; ok {
				return true
			}
		}
		return false
	}

	parts := make([]model.Part, 0, len(r.data))
	for _, repoPart := range r.data {
		part := repoConverter.PartToModel(repoPart)
		if len(filter.UUIDS) > 0 && !slices.Contains(filter.UUIDS, part.UUID) {
			continue
		}
		if len(filter.Names) > 0 && !slices.Contains(filter.Names, part.Name) {
			continue
		}
		// TODO check this case carefully
		if len(filter.Categories) > 0 && !slices.Contains(filter.Categories, part.Category) {
			continue
		}
		if len(filter.ManufacturerCountries) > 0 && !slices.Contains(filter.ManufacturerCountries, part.Manufacturer.Country) {
			continue
		}
		if len(filter.Tags) > 0 && (len(part.Tags) == 0 || !hasTagListsIntersectionFunc(filter.Tags, part.Tags)) {
			continue
		}
		parts = append(parts, part)
	}

	return parts, nil
}
