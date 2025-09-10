package testutil

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/xgmsx/rsf/inventory/internal/model"
	"github.com/xgmsx/rsf/inventory/internal/utils"
)

func GetNewPart() model.Part {
	return model.Part{
		UUID:          gofakeit.UUID(),
		Name:          gofakeit.ProductName(),
		Description:   gofakeit.Product().Description,
		Price:         gofakeit.Product().Price,
		StockQuantity: int64(gofakeit.Int8()),
		Category:      model.Category_CATEGORY_ENGINE,
		Dimensions: &model.Dimensions{
			Length: float64(gofakeit.Int8()),
			Width:  float64(gofakeit.Int8()),
			Height: float64(gofakeit.Int8()),
			Weight: float64(gofakeit.Int8()),
		},
		Manufacturer: &model.Manufacturer{
			Name:    gofakeit.Company(),
			Country: gofakeit.Country(),
			Website: gofakeit.URL(),
		},
		Tags: []string{"engine"},
		Metadata: map[string]*model.Value{
			"power_output":    {DoubleValue: utils.ToPtr(9.5)},
			"is_experimental": {BoolValue: utils.ToPtr(true)},
		},
		CreatedAt: gofakeit.Date(),
		UpdatedAt: gofakeit.Date(),
	}
}

func GetNewParts(count int) []model.Part {
	parts := make([]model.Part, count)
	for i := range parts {
		parts[i] = GetNewPart()
	}
	return parts
}
