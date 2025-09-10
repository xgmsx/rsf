package part

import (
	"context"
	"slices"
	"sync"
	"time"

	"github.com/xgmsx/rsf/inventory/internal/model"
	def "github.com/xgmsx/rsf/inventory/internal/repository"
	"github.com/xgmsx/rsf/inventory/internal/utils"
)

var _ def.PartRepository = (*partsRepository)(nil)

type partsRepository struct {
	mu   sync.RWMutex
	data map[string]*model.Part
}

func NewPartRepository() *partsRepository {
	now := time.Now()

	part1 := &model.Part{
		UUID:          "111e4567-e89b-12d3-a456-426614174001",
		Name:          "Hyperdrive Engine",
		Description:   "A class-9 hyperdrive engine capable of faster-than-light travel.",
		Price:         450000.00,
		StockQuantity: 3,
		Category:      model.Category_CATEGORY_ENGINE,
		Dimensions: &model.Dimensions{
			Length: 120.0,
			Width:  80.0,
			Height: 100.0,
			Weight: 500.0,
		},
		Manufacturer: &model.Manufacturer{
			Name:    "Hyperdrive Corp",
			Country: "USA",
			Website: "https://hyperdrive.example.com",
		},
		Tags: []string{"engine", "hyperdrive", "space"},
		Metadata: map[string]*model.Value{
			"power_output":    {DoubleValue: utils.ToPtr(9.5)},
			"is_experimental": {BoolValue: utils.ToPtr(true)},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	part2 := &model.Part{
		UUID:          "222e4567-e89b-12d3-a456-426614174002",
		Name:          "Quantum Shield Generator",
		Description:   "Advanced shield generator providing protection against cosmic radiation.",
		Price:         175000.00,
		StockQuantity: 5,
		Category:      model.Category_CATEGORY_SHIELD,
		Dimensions: &model.Dimensions{
			Length: 60.0,
			Width:  40.0,
			Height: 50.0,
			Weight: 150.0,
		},
		Manufacturer: &model.Manufacturer{
			Name:    "Quantum Tech",
			Country: "Germany",
			Website: "https://quantumtech.example.com",
		},
		Tags: []string{"shield", "quantum", "defense"},
		Metadata: map[string]*model.Value{
			"energy_consumption": {DoubleValue: utils.ToPtr(3.2)},
			"warranty_years":     {Int64Value: utils.ToPtr(int64(5))},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	parts := make(map[string]*model.Part, 2)

	parts[part1.UUID] = part1
	parts[part2.UUID] = part2

	return &partsRepository{
		data: parts,
	}
}

func (r *partsRepository) GetPart(_ context.Context, uuid string) (model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	part, ok := r.data[uuid]
	if !ok {
		return model.Part{}, model.ErrPartDoesNotExist
	}
	return *part, nil
}

func (r *partsRepository) ListParts(ctx context.Context, filter *model.PartsFilter) ([]model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []model.Part

	for _, part := range r.data {
		if filter != nil {
			if len(filter.UUIDs) > 0 && !slices.Contains(filter.UUIDs, part.UUID) {
				continue
			}
			if len(filter.Names) > 0 && !slices.Contains(filter.Names, part.Name) {
				continue
			}
			if len(filter.Categories) > 0 && !slices.Contains(filter.Categories, model.Category(part.Category)) {
				continue
			}
			if len(filter.ManufacturerCountries) > 0 && (part.Manufacturer == nil || !slices.Contains(filter.ManufacturerCountries, part.Manufacturer.Country)) {
				continue
			}
			if len(filter.Tags) > 0 {
				skip := false
				for _, tag := range filter.Tags {
					if !slices.Contains(part.Tags, tag) {
						skip = true
						break
					}
				}
				if skip {
					continue
				}
			}
		}
		result = append(result, *part)
	}
	return result, nil
}
