package converter

import (
	"fmt"

	"github.com/ploskirev/go-rocket/inventory/internal/model"
	inventory_v1 "github.com/ploskirev/go-rocket/shared/pkg/proto/inventory/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func PartToProto(p model.Part) *inventory_v1.Part {
	updatedAt := timestamppb.Timestamp{}
	if p.UpdatedAt != nil {
		updatedAt = *timestamppb.New(*p.UpdatedAt)
	}

	return &inventory_v1.Part{
		Uuid:          p.UUID,
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		StockQuantity: p.StockQuantity,
		Tags:          p.Tags,
		Category:      CategoryToProto(&p.Category),
		Dimensions:    DimensionsToProto(&p.Dimensions),
		Manufacturer:  ManufacurerToProto(&p.Manufacturer),
		Metadata:      MetadataToProto(p.Metadata),
		CreatedAt:     timestamppb.New(p.CreatedAt),
		UpdatedAt:     &updatedAt,
	}
}

func MetadataToProto(m map[string]any) map[string]*inventory_v1.Value {
	r := make(map[string]*inventory_v1.Value, len(m))

	for k, v := range m {
		switch t := v.(type) {
		case int64:
			r[k] = &inventory_v1.Value{ValueType: &inventory_v1.Value_Int64Value{Int64Value: t}}

		case float64:
			r[k] = &inventory_v1.Value{ValueType: &inventory_v1.Value_DoubleValue{DoubleValue: t}}

		case string:
			r[k] = &inventory_v1.Value{ValueType: &inventory_v1.Value_StringValue{StringValue: t}}

		case bool:
			r[k] = &inventory_v1.Value{ValueType: &inventory_v1.Value_BoolValue{BoolValue: t}}

		default:
			fmt.Println("Unsupported type for metadata")
		}
	}

	return r
}

func CategoryToProto(c *model.Category) inventory_v1.Category {
	return inventory_v1.Category(*c)
}

func DimensionsToProto(d *model.Dimensions) *inventory_v1.Dimensions {
	return &inventory_v1.Dimensions{
		Length: d.Length,
		Height: d.Height,
		Width:  d.Width,
		Weight: d.Weight,
	}
}

func ManufacurerToProto(m *model.Manufacturer) *inventory_v1.Manufacturer {
	return &inventory_v1.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}

func FiltersToModel(f *inventory_v1.PartsFilter) *model.Filters {
	if f == nil {
		return &model.Filters{}
	}

	categories := make([]model.Category, 0, len(f.Categories))
	for _, c := range f.Categories {
		categories = append(categories, model.Category(c))
	}

	return &model.Filters{
		UUIDS:                 f.Uuids,
		Names:                 f.Names,
		Categories:            categories,
		ManufacturerCountries: f.ManufacturerCountries,
		Tags:                  f.Tags,
	}
}

func CategoryToModel(c inventory_v1.Category) model.Category {
	return model.Category(c)
}
