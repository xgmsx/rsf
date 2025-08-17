package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/xgmsx/rsf/inventory/internal/model"
	inventoryV1 "github.com/xgmsx/rsf/shared/pkg/proto/inventory/v1"
)

func PartFilterFromProto(f *inventoryV1.PartsFilter) *model.PartsFilter {
	if f == nil {
		return nil
	}
	categories := make([]model.Category, 0, len(f.Categories))
	for _, c := range f.Categories {
		categories = append(categories, model.Category(c))
	}
	return &model.PartsFilter{
		UUIDs:                 f.Uuids,
		Names:                 f.Names,
		Categories:            categories,
		ManufacturerCountries: f.ManufacturerCountries,
		Tags:                  f.Tags,
	}
}

func PartToProto(p model.Part) *inventoryV1.Part {
	return &inventoryV1.Part{
		Uuid:          p.UUID,
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		StockQuantity: p.StockQuantity,
		Category:      inventoryV1.Category(p.Category),
		Dimensions:    DimensionsToProto(p.Dimensions),
		Manufacturer:  ManufacturerToProto(p.Manufacturer),
		Tags:          p.Tags,
		Metadata:      MetadataToProto(p.Metadata),
		CreatedAt:     timestamppb.New(p.CreatedAt),
		UpdatedAt:     timestamppb.New(p.UpdatedAt),
	}
}

func DimensionsToProto(d *model.Dimensions) *inventoryV1.Dimensions {
	if d == nil {
		return nil
	}
	return &inventoryV1.Dimensions{
		Length: d.Length,
		Width:  d.Width,
		Height: d.Height,
		Weight: d.Weight,
	}
}

func ManufacturerToProto(m *model.Manufacturer) *inventoryV1.Manufacturer {
	if m == nil {
		return nil
	}
	return &inventoryV1.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}

func MetadataToProto(meta map[string]*model.Value) map[string]*inventoryV1.Value {
	result := make(map[string]*inventoryV1.Value, len(meta))
	for k, v := range meta {
		result[k] = ValueToProto(v)
	}
	return result
}

func ValueToProto(v *model.Value) *inventoryV1.Value {
	if v == nil {
		return nil
	}
	switch {
	case v.DoubleValue != nil:
		return &inventoryV1.Value{Kind: &inventoryV1.Value_DoubleValue{DoubleValue: *v.DoubleValue}}
	case v.Int64Value != nil:
		return &inventoryV1.Value{Kind: &inventoryV1.Value_Int64Value{Int64Value: *v.Int64Value}}
	case v.BoolValue != nil:
		return &inventoryV1.Value{Kind: &inventoryV1.Value_BoolValue{BoolValue: *v.BoolValue}}
	case v.StringValue != nil:
		return &inventoryV1.Value{Kind: &inventoryV1.Value_StringValue{StringValue: *v.StringValue}}
	default:
		return &inventoryV1.Value{}
	}
}
