package cards

import (
	"fmt"

	"github.com/ianhecker/pokemon-tcg-services/internal/pokemontcg"
)

type Variant struct {
	Condition   pokemontcg.Condition `json:"condition"`
	Price       float64              `json:"price"`
	MinPrice30d float64              `json:"minPrice30d"`
	MaxPrice30d float64              `json:"maxPrice30d"`
	LastUpdated int64                `json:"lastUpdated"`
}

type Variants []Variant

func (variants Variants) ToMap() (map[pokemontcg.Condition]Prices, error) {

	m := make(map[pokemontcg.Condition]Prices)
	for _, v := range variants {
		condition := v.Condition

		_, exists := m[condition]
		if exists {
			return nil, fmt.Errorf("duplicated condition in prices: '%s'", condition)
		}
		m[condition] = MakePricesFromVariant(v)
	}
	return m, nil
}
