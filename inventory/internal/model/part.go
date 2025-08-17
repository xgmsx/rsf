package model

import "time"

type PartCategory string

type PartsFilter struct {
	UUIDs                 []string
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
	Tags                  []string
}

type Dimensions struct {
	Length float64
	Width  float64
	Height float64
	Weight float64
}

type Manufacturer struct {
	Name    string
	Country string
	Website string
}

type Value struct {
	DoubleValue *float64
	Int64Value  *int64
	BoolValue   *bool
	StringValue *string
}

type Part struct {
	UUID          string
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Category      Category
	Dimensions    *Dimensions
	Manufacturer  *Manufacturer
	Tags          []string
	Metadata      map[string]*Value
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Category int32

const (
	Category_CATEGORY_UNSPECIFIED Category = 0
	Category_CATEGORY_ENGINE      Category = 1
	Category_CATEGORY_FUEL        Category = 2
	Category_CATEGORY_PORTHOLE    Category = 3
	Category_CATEGORY_WING        Category = 4
	Category_CATEGORY_SHIELD      Category = 5
)
