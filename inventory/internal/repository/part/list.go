package part

import (
	"context"
	"slices"

	"github.com/ploskirev/go-rocket/inventory/internal/model"
	"github.com/ploskirev/go-rocket/inventory/internal/repository/converter"
	repoModel "github.com/ploskirev/go-rocket/inventory/internal/repository/model"
)

func (pr *partRepository) ListParts(_ context.Context, filters *model.Filters) []*model.Part {
	pr.mu.RLock()
	defer pr.mu.RUnlock()

	parts := make([]*repoModel.Part, 0, len(pr.storage))

	for _, p := range pr.storage {
		parts = append(parts, p)
	}

	if filters != nil {
		f := converter.FiltersToRepoFilters(filters)
		parts = filterParts(parts, f)
	}

	listParts := make([]*model.Part, 0, len(parts))
	for _, p := range parts {
		listParts = append(listParts, converter.PartToModel(p))
	}

	return listParts
}

type FiltersMaps struct {
	byUUIDs                 map[string]struct{}
	byNames                 map[string]struct{}
	byCategories            map[repoModel.Category]struct{}
	byManufactoredCountries map[string]struct{}
	byTagsSlice             []string
}

func convertFiltersToMaps(filters *repoModel.Filters) *FiltersMaps {
	filterUUIDsMap := make(map[string]struct{}, len(filters.UUIDS))
	for _, f := range filters.UUIDS {
		filterUUIDsMap[f] = struct{}{}
	}
	filterNamesMap := make(map[string]struct{}, len(filters.Names))
	for _, f := range filters.Names {
		filterNamesMap[f] = struct{}{}
	}
	filterCategoriesMap := make(map[repoModel.Category]struct{}, len(filters.Categories))
	for _, f := range filters.Categories {
		filterCategoriesMap[f] = struct{}{}
	}
	filterManufactoredCountriesMap := make(map[string]struct{}, len(filters.ManufacturerCountries))
	for _, f := range filters.ManufacturerCountries {
		filterManufactoredCountriesMap[f] = struct{}{}
	}

	return &FiltersMaps{
		byUUIDs:                 filterUUIDsMap,
		byNames:                 filterNamesMap,
		byCategories:            filterCategoriesMap,
		byManufactoredCountries: filterManufactoredCountriesMap,
		byTagsSlice:             filters.Tags,
	}
}

func filterParts(parts []*repoModel.Part, filters *repoModel.Filters) []*repoModel.Part {
	filtersMaps := convertFiltersToMaps(filters)
	filteredParts := make([]*repoModel.Part, 0, len(parts))

	for _, part := range parts {
		if len(filtersMaps.byUUIDs) > 0 {
			if _, ok := filtersMaps.byUUIDs[part.UUID]; !ok {
				continue
			}
		}
		if len(filtersMaps.byNames) > 0 {
			if _, ok := filtersMaps.byNames[part.Name]; !ok {
				continue
			}
		}
		if len(filtersMaps.byCategories) > 0 {
			if _, ok := filtersMaps.byCategories[part.Category]; !ok {
				continue
			}
		}
		if len(filtersMaps.byManufactoredCountries) > 0 {
			if _, ok := filtersMaps.byManufactoredCountries[part.Manufacturer.Country]; !ok {
				continue
			}
		}
		if len(filtersMaps.byTagsSlice) > 0 {
			skip := true
			for _, t := range part.Tags {
				if slices.Contains(filtersMaps.byTagsSlice, t) {
					skip = false
					break
				}
			}

			if skip {
				continue
			}
		}

		filteredParts = append(filteredParts, part)
	}

	return filteredParts
}
