package cards

import (
	"github.com/ianhecker/pokemon-tcg-services/internal/pokemontcg"
)

type Variant struct {
	Printing    Printing             `json:"printing"`
	Condition   pokemontcg.Condition `json:"condition"`
	Price       float64              `json:"price"`
	MinPrice30d float64              `json:"minPrice30d"`
	MaxPrice30d float64              `json:"maxPrice30d"`
	LastUpdated int64                `json:"lastUpdated"`
}
