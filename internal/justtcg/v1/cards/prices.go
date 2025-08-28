package cards

import (
	"time"
)

type Prices struct {
	Market      float64   `json:"market"`
	MinPrice30d float64   `json:"minPrice30d"`
	MaxPrice30d float64   `json:"maxPrice30d"`
	LastUpdated time.Time `json:"lastUpdated"`
}

func MakePrices(
	market float64,
	minPrice30d float64,
	maxPrice30d float64,
	lastUpdated int64,
) Prices {
	return Prices{
		Market:      market,
		MinPrice30d: minPrice30d,
		MaxPrice30d: maxPrice30d,
		LastUpdated: time.Unix(lastUpdated, 0).UTC(),
	}
}

func MakePricesFromVariant(variant Variant) Prices {
	return MakePrices(
		variant.Price,
		variant.MinPrice30d,
		variant.MaxPrice30d,
		variant.LastUpdated,
	)
}
