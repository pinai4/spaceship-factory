package model

import (
	"time"
)

type Category string

const (
	CategoryUnspecified Category = "UNSPECIFIED"
	CategoryEngine      Category = "ENGINE"
	CategoryFuel        Category = "FUEL"
	CategoryPorthole    Category = "PORTHOLE"
	CategoryWing        Category = "WING"
)

func (c Category) IsValid() bool {
	switch c {
	case CategoryUnspecified, CategoryEngine, CategoryFuel, CategoryPorthole, CategoryWing:
		return true
	default:
		return false
	}
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

type Part struct {
	UUID          string
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Category      Category
	Dimensions    Dimensions
	Manufacturer  Manufacturer
	Tags          []string
	Metadata      map[string]any
	CreatedAt     time.Time
	UpdatedAt     *time.Time
}
