package generate

import (
	"time"

	"github.com/ianhecker/pokemon-tcg-services/internal/justtcg/v1/cards"
)

func ftop(f float64) *cards.Price {
	price := cards.Price(f)
	return &price
}

func tptr(t time.Time) *time.Time {
	return &t
}
