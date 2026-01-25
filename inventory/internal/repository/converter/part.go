package converter

import (
	"fmt"

	"github.com/ploskirev/go-rocket/inventory/internal/model"
	repoModel "github.com/ploskirev/go-rocket/inventory/internal/repository/model"
	inventory_v1 "github.com/ploskirev/go-rocket/shared/pkg/proto/inventory/v1"
)

func PartToModel(part *repoModel.Part) *model.Part {
	return &model.Part{
		UUID:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      model.Category(part.Category),
		Dimensions:    *PartDimensionsToModel(&part.Dimensions),
		Manufacturer:  *PartManufacturerToModel(&part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      part.Metadata,
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
	}
}

func PartDimensionsToModel(d *repoModel.Dimensions) *model.Dimensions {
	return &model.Dimensions{
		Length: d.Length,
		Height: d.Height,
		Width:  d.Width,
		Weight: d.Weight,
	}
}

func PartManufacturerToModel(m *repoModel.Manifacturer) *model.Manufacturer {
	return &model.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}

func MetadataToModel(m map[string]inventory_v1.Value) map[string]any {
	r := make(map[string]any, len(m))

	for k, v := range m {
		switch t := v.ValueType.(type) {
		case *inventory_v1.Value_Int64Value:
			r[k] = t.Int64Value

		case *inventory_v1.Value_DoubleValue:
			r[k] = t.DoubleValue

		case *inventory_v1.Value_StringValue:
			r[k] = t.StringValue

		case *inventory_v1.Value_BoolValue:
			r[k] = t.BoolValue

		default:
			fmt.Println("Unsupported type for metadata")
		}
	}

	return r
}

func FiltersToRepoFilters(f *model.Filters) *repoModel.Filters {
	categories := make([]repoModel.Category, 0, len(f.Categories))
	for _, c := range f.Categories {
		categories = append(categories, repoModel.Category(c))
	}

	return &repoModel.Filters{
		UUIDS:                 f.UUIDS,
		Names:                 f.Names,
		Categories:            categories,
		ManufacturerCountries: f.ManufacturerCountries,
		Tags:                  f.Tags,
	}
}
